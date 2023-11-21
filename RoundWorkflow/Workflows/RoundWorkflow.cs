using Dapr.Workflow;
using DurableTask.Core.Exceptions;
using System;

using RoundWorkflow.Activities;
using RoundWorkflow.Models;

namespace RoundWorkflow.Workflows
{
    public class VolleyballRoundWorkflow : Workflow<RoundInput, RoundResult>
    {
        public override async Task<RoundResult> RunAsync(WorkflowContext context, RoundInput roundInputs)
        {
            string roundId = context.InstanceId;

            // Notify 
            await context.CallActivityAsync(
                nameof(NotifyActivity),
                new Notification($"Starting new Tournament round: {roundId}"));

            // Pair teams  

            context.SetCustomStatus("Pairing teams for round");

            var roundPairings = await context.CallActivityAsync<List<Pair>>(
                nameof(PairTeamsActivity),
                roundInputs.Teams);

            context.SetCustomStatus($"Starting Games for round {roundInputs.Round}");

            // Kicking off Games and retrieving results
            var games = new List<Task<GameResult>>();
            var roundResults = Array.Empty<GameResult>();

            try
            {
                foreach (var pairing in roundPairings)
                {
                    games.Add(context.CallActivityAsync<GameResult>(
                        nameof(StartGameActivity),
                        pairing));
                }

                roundResults = await Task.WhenAll(games);
            }
            catch (Exception ex)
            {

                if (ex.InnerException is TaskFailedException)
                {
                    await context.CallActivityAsync(
                        nameof(NotifyActivity),
                        new Notification($"Round {roundInputs.Round} failed due to {ex.Message}"));
                    context.SetCustomStatus("Round failed");
                    return new RoundResult(Succeeded: false, null);
                }
            }

            // TODO: Activity to determine winners and formulate state payload 

            context.SetCustomStatus($"Saving state for round {roundInputs.Round}");

            try
            {
                // Save the payload into state
                await context.CallActivityAsync(
                    nameof(SaveRoundActivity),
                    roundResults);
            }
            catch (TaskFailedException)
            {
                // Notify the user saving state failed
                await context.CallActivityAsync(
                    nameof(NotifyActivity),
                    new Notification($"Saving state for round {roundInputs.Round} failed."));

                context.SetCustomStatus("Saving round state failed.");

                return new RoundResult(Succeeded: false, null);
            }

            return new RoundResult(true, null);
        }
    }
}
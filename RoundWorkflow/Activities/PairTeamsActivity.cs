using Dapr.Client;
using Dapr.Workflow;
using RoundWorkflow.Models;

namespace RoundWorkflow.Activities
{
    class PairTeamsActivity : WorkflowActivity<List<string>, List<Pair>>
    {
        readonly ILogger _logger;

        public PairTeamsActivity(ILoggerFactory loggerFactory, DaprClient client)
        {
            _logger = loggerFactory.CreateLogger<PairTeamsActivity>();
        }

        public override Task<List<Pair>> RunAsync(WorkflowActivityContext context, List<string> teams)
        {
            // Randomly order the array 
            Random rand = new();
            var reorderedTeams = teams.OrderBy(x => rand.Next()).ToArray();

            var pairings = new List<Pair>();

            for (int x = 0; x <= reorderedTeams.Length; x++)
            {
                var teamA = reorderedTeams[x];
                var teamB = reorderedTeams[x += 1];

                pairings.Add(new Pair(teamA, teamB));

                _logger.LogInformation("Matchup Paired: {teamA},{teamB}", teamA, teamB);

            }

            _logger.LogInformation("All teams successfully paired");

            return Task.FromResult<List<Pair>>(pairings);

        }
    }
}
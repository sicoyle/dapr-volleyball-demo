using Dapr.Client;
using Dapr.Workflow;
using RoundWorkflow.Models;

namespace RoundWorkflow.Activities
{
    class StartGameActivity : WorkflowActivity<Pair, GameResult>
    {
        readonly ILogger _logger;
        private readonly DaprClient _client;

        public StartGameActivity(ILoggerFactory loggerFactory, DaprClient client)
        {
            _logger = loggerFactory.CreateLogger<StartGameActivity>();
            _client = client;
        }

        public override Task<GameResult> RunAsync(WorkflowActivityContext context, Pair matchup)
        {
            //TODO: CASSIE
            // call the game sim for a particular pairing via service invocation
            // Return a GameResult object to the workflow from the activity
        }
    }
}
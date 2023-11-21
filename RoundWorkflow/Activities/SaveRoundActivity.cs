using System.Resources;
using Dapr.Client;
using Dapr.Workflow;
using RoundWorkflow.Models;

namespace RoundWorkflow.Activities
{
    class SaveRoundActivity : WorkflowActivity<GameResult[], object?>
    {
        readonly ILogger _logger;
        private readonly DaprClient _client;

        public SaveRoundActivity(ILoggerFactory loggerFactory, DaprClient client)
        {
            _logger = loggerFactory.CreateLogger<StartGameActivity>();
            _client = client;
        }

        public override Task<object?> RunAsync(WorkflowActivityContext context, GameResult[] gameResults)
        {
            //change the shape of the gameresults object and store in state store while using outbox (i.e. state transaction and publish a message to roundresults topic or something similar)
            return null;
        }
    }
}
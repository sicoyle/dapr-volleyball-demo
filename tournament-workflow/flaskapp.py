from flask import Flask, request, jsonify
import threading
from time import sleep
from dapr.ext.workflow import WorkflowRuntime, DaprWorkflowContext, WorkflowActivityContext
from dapr.conf import Settings
from dapr.clients import DaprClient
from dapr.clients.exceptions import DaprInternalError
import requests

app = Flask(__name__)

settings = Settings()

counter = 0
instanceId = "exampleInstanceID26"
workflowComponent = "dapr"
workflowName = "hello_world_wf"
inputData = "Hi Counter!"
workflowOptions = dict()
workflowOptions["task_queue"] =  "testQueue"
eventName = "event1"
eventData = "eventData"
nonExistentIDError = "no such instance exists"


def hello_world_wf(ctx: DaprWorkflowContext, input):
    print(f'{input}')
    yield ctx.call_activity(hello_act, input=1)
    yield ctx.call_activity(call_child_workflow, input=1)
    yield ctx.call_activity(generate_report, input=1)

def call_child_workflow(ctx: WorkflowActivityContext, input):
    print(f'Sam TODO to call other workflow!', flush=True)

    # API endpoint for starting a workflow (replace with your actual endpoint)
    workflow_api_endpoint = "http://localhost:3500/v1.0-beta1/workflows/dapr/RoundVolleyballWorkflow/start?instanceID=volleyball"

    # Input data for the child workflow (replace with your actual data)
    child_workflow_input = {"input_data": "TODO Child Workflow Input Game Data"}

    dapr_app_id = "workflow"

      # Headers to include in the request
    headers = {
        "Content-Type": "application/json",
        "dapr-app-id": dapr_app_id,
    }

    # Make a POST request to start the child workflow
    response = requests.post(workflow_api_endpoint, json=child_workflow_input, headers=headers)


    if response.status_code == 202:
        print("Child workflow started successfully.")
    else:
        print(f"Failed to start child workflow. Status code: {response.status_code}, Response: {response.text}")

    
def generate_report(ctx: WorkflowActivityContext, input):
    print(f'Sam success work on report!', flush=True)

def hello_act(ctx: WorkflowActivityContext, input):
    global counter
    counter += input
    print(f'Hi Sam! New counter value is: {counter}!', flush=True)

def run_workflow():
    with DaprClient() as d:
        host = settings.DAPR_RUNTIME_HOST
        port = settings.DAPR_GRPC_PORT
        workflowRuntime = WorkflowRuntime(host, port)
        workflowRuntime = WorkflowRuntime()
        workflowRuntime.register_workflow(hello_world_wf)
        workflowRuntime.register_activity(hello_act)
        workflowRuntime.register_activity(call_child_workflow)
        workflowRuntime.register_activity(generate_report)

        workflowRuntime.start()

        print("==========Start Counter Increase as per Input:==========")
        start_resp = d.start_workflow(
            instance_id=instanceId,
            workflow_component=workflowComponent,
            workflow_name=workflowName,
            input=inputData,
            workflow_options=workflowOptions
        )
        print(f"start_resp {start_resp.instance_id}")

        # while True:
        #     getResponse = d.get_workflow(instance_id=instanceId, workflow_component=workflowComponent)
        #     print(f"sam status of workflow: {getResponse.runtime_status}")

        #     if getResponse.runtime_status == "COMPLETED":
        #         print("Workflow completed.")
        #         break

        sleep(15)


        # d.purge_workflow(instance_id=instanceId, workflow_component=workflowComponent)
        # try:
        #     getResponse = d.get_workflow(instance_id=instanceId, workflow_component=workflowComponent)
        # except DaprInternalError as err:
        #     if nonExistentIDError in err._message:
        #         print("Instance Successfully Purged")

        workflowRuntime.shutdown()

# Define an endpoint for starting the workflow
@app.route('/start-workflow', methods=['POST'])
def start_workflow_endpoint():
    # Retrieve input data from the request
    input_data = request.json.get('input_data', None)

    # Run the workflow in a separate thread to avoid blocking the HTTP server
    workflow_thread = threading.Thread(target=run_workflow)
    workflow_thread.start()

    return jsonify({'status': 'Workflow started successfully'}), 202

if __name__ == '__main__':
    # Run the Flask app
    app.run(port=5001)

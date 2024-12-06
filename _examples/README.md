# Examples

This directory contains a collection of examples that illustrate how to use copilot-go.

| Example | Description |
|---------|-------------|
| [auth](auth) | An agent that handles auth. |
| [azure_openai](azure-openai) | An agent that uses Azure OpenAI Service. |
| [copilot_api](azure-copilot_api) | An agent that uses the Copilot API. |
| [events](events) | An agent that demonstrates responding with errors, confirmations, and references. |
| [github](github) | An agent that demonstrates using the context and the GitHub API to get details about GitHub resources. |


## Running the Examples

To run the example, you need to copy the `.env.sample` file to `.env` and fill in the required values.

Once you have the `.env` file, you can run the example with the following command:

```sh
go run main.go
```

Or you can use the vscode launch configuration to run the example.

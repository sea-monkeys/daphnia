# Create a store from a Dataset

The dataset comes from https://huggingface.co/datasets/MattCoddity/dockerNLcommands

## Create a store from a Dataset

```bash
go test -v create_vectordb_test.go
```

Wait for a moment (2 or 3 minutes on a MacBook Pro M2) and you will get a message like this:

```bash
Total commands processed: 2415
PASS
ok      04-ollama-docker        163.504s
```

## Run the samples

```bash
go run main.go
```

```bash
go test -v q1_test.go
```

When the LLM is ready, the answer response time is around 100 millisecond.

module 03-ollama-demo

go 1.23.1

require (
	github.com/ollama/ollama v0.3.13
	github.com/sea-monkeys/artemia v0.0.0 // indirect
	github.com/sea-monkeys/daphnia v0.0.2

)

require github.com/google/uuid v1.6.0 // indirect

replace github.com/sea-monkeys/daphnia => ..

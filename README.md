# dexetera

A home for simulation archetypes and AI agent training environments.

## Build the example sim

In order to build the example sim with WebAssembly, you run:

```shell
GOOS=js GOARCH=wasm go build -o ./app/example_sim.wasm ./cmd/example_sim/main.go 
```

## Run the example sim

```shell
# view the app running at http://localhost:8000
cd app/ && python3 -m http.server 8000
```

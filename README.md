# dexetera

A home for interactive simulation environments.

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

You can then interact with the sim using [DexAct](https://pypi.org/project/dexact/).

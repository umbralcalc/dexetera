# dexetera...

A toolkit to build purely-frontend web applications for interactive [stochadex](https://github.com/umbralcalc/stochadex) simulations.

## Build the example simulation

In order to build an example simulation with WebAssembly, you run:

```shell
GOOS=js GOARCH=wasm go build -o ./app/src/example_sim_1/main.wasm ./cmd/example_sim_1/main.go 
```

## Run the example simulation

```shell
# view the app running at http://localhost:8000/example_sim_1.html
cd app/ && python3 -m http.server 8000
```

You can then interact with the simulation using [dexAct](https://pypi.org/project/dexact/).

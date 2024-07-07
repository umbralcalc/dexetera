# dexetera...

Purely-frontend web applications for interactive [stochadex](https://github.com/umbralcalc/stochadex) simulations.

For more context, here is the associated WIP article: [https://umbralcalc.github.io/posts/dexetera.html](https://umbralcalc.github.io/posts/dexetera.html).

## Build a simulation

To build one of the simulations with WebAssembly, you run:

```shell
GOOS=js GOARCH=wasm go build -o ./app/src/flounceball/main.wasm ./cmd/flounceball/main.go 
```

## Run a simulation

```shell
# view the app running at http://localhost:8000/flounceball.html
cd app/ && python3 -m http.server 8000
```

You can then interact with the simulation using [dexAct](https://pypi.org/project/dexact/).

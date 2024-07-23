# dexetera...

The dexetera framework is for developing purely-frontend web applications for interactive [stochadex](https://github.com/umbralcalc/stochadex) simulations. The compiled simulation can be 'stepped' in time through interactions with a websocket server on the user's machine. Custom visualisations for these interactive simulations can also be developed in the JavaScript code which runs the simulation steps (see, e.g., `app/flounceball.html`). For more context, here is the associated technical article about this project: [https://umbralcalc.github.io/posts/dexetera.html](https://umbralcalc.github.io/posts/dexetera.html).

## Build a simulation

To build one of the simulations with WebAssembly, you run:

```shell
GOOS=js GOARCH=wasm go build -o ./app/src/flounceball/main.wasm ./cmd/flounceball/main.go 
```

The resulting `main.wasm` binary can then be executed in JavaScript code. When executed, this registers `stepSimulation` function as a callable.

## Run a simulation

```shell
# view the app running at http://localhost:8000/flounceball.html
cd app/ && python3 -m http.server 8000
```

You can then interact with the simulation using any desired websocket server application. To make things straightforward for the python developer, we've made a PyPI package called [dexAct](https://pypi.org/project/dexact/) that simplifies this experience into:

```shell
# install dexact into your local python environment
pip install dexact

# run the python action server while the app is running at http://localhost:8000/flounceball.html
python cmd/flounceball/action_server.py
```

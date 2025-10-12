# dexetera...

Decision-making games for the python programmer built using the [stochadex](https://github.com/umbralcalc/stochadex) simulation framework. Each of these games is a purely-frontend web application built by compiling an interactive stochadex simulation into WebAssembly. The games are then 'stepped' in time through updates sent via a python websocket server on the user's machine and visualisations for them are all written in JavaScript (see, e.g., `app/flounceball.html`).

## Build a game

To build one of the games with WebAssembly, you run:

```shell
GOOS=js GOARCH=wasm go build -o ./app/src/flounceball/main.wasm ./cmd/flounceball/main.go 
```

The resulting `main.wasm` binary can then be executed in JavaScript code. When executed, this registers `stepSimulation` function as a callable.

## Run a game

```shell
# view the app running at http://localhost:8000/flounceball.html
cd app/ && python -m http.server 8000
```

You can then control the game using any desired websocket server application. To make things straightforward for the python programmer, there's a PyPI package called [dexAct](https://pypi.org/project/dexact/) that simplifies the experience into:

```shell
# install dexact into your local python environment
pip install dexact

# run the python action server while the app is 
# running at http://localhost:8000/flounceball.html
python cmd/flounceball/action_server.py
```

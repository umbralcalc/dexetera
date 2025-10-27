## Context on what this repo is for

This repo is to facilitate Go programmers to create decision-making games by exposing simulations built using the [stochadex](https://github.com/umbralcalc/stochadex) framework as WebAssembly-compiled frontends in the users' browser. 

The intended way that a user interacts with this application is:
- views the github pages site for this repo in index.html
- clicks into playing one of the games via, e.g., the `/app/minimal_example.html` endpoint
- runs their python websocket server code locally where the user outputs an action state (the ActionState protocol in action_state.proto) to the client and this updates what they visualise in their browser (so this is where animated visualisations must be customisable) and provides the main driving mechanism for the simulation to progress
- in response to this action state message, the websocket client in the browser written by us then steps the simulation via a JS callback function in the WebAssembly-compiled simulation stepper and the simulation then returns PartitionState messages which we send back to the user's server code to tell it how the state is changed... note how it can be configured to filter which PartitionState messages even get sent back to the user in, e.g., the app/minimal_example.html line with `serverPartitionNames: ["counter_state"],`

Each of these games is meant to be a purely-frontend web application built by compiling an interactive stochadex simulation into WebAssembly. The games are then 'stepped' in time through updates sent via a python websocket server on the user's machine and visualisations for them are all written in JavaScript (see, e.g., `app/minimal_example.html`).

The `pkg/simio` package has been updated to the latest version of the stochadex framework and this supports a JS callback function being passed into it in the `app/` code.

The games writtten using the stochadex framework can be found in the `pkg/games` package.

## What I'd like you to help me with

I'm in the middle of refactoring the code in this repo to make it more maintainable and easily extensible to new games in the future. I'd like you to help me with this by suggesting refactoring changes towards this longer-term goal but get me to approve these step-by-step before you implement them.

I want the Go programmer to be able to specify a lot of how this JS frontend looks for these decision-making games by generating the visualisations with Go code if possible, while doing this in a way which is maintainable and works well with the abstractions already provided by the stochadex package (at the very least the ones exposed in the State protocol buffer layer, but ideally even more at the `simulator.Iteration` and `simulator.OutputFunction` level too).

I don't want the stochadex framework to have to change, so all of the proposed changes here should be extensions written in this repo ideally, but I'm willing to consider feature extensions to the stochadex if you think it's really important.

A good first step should be to update the pkg/games/game.go and pkg/games/minimal_example.go to use the simulator.ConfigGenerator to build the simulation configs step-by-step and then you can generate the implementation and settings configs needed in cmd/minimal_example/main.go using .GenerateConfigs()

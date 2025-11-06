## Context to help building new games in this repo

Please read the README.md to understand the general flow for building new games in this repo. There is only one very simple example of the end-to-end flow (team_sport) and you can gather the general structure of how this fits together from looking at its files:

- pkg/team_sport - backend game logic and configuration
- team_sport/ - generated game
- cmd/team_sport/generate_game - script used to generate the team_sport/ folder
- cmd/team_sport/register_step - compiles to `team_sport/src/main.wasm` using team_sport/build.sh script
- cmd/team_sport/action_server.py - how the user interacts with the simulation site with their python action state sending server
- cmd/team_sport/cheatsheet.md - a handy description of what the keys and values in the state dictionary received by the user in their Python `ActionTaker.take_next_action` code can be interpreted to mean

## Context to help write game simulations into the `simulator.ConfigGenerator`

Please start by reading the docs here to understand the API: https://umbralcalc.github.io/stochadex/pkg/simulator.html

The simulation framework followed a shared-state actor pattern where each `simulator.PartitionConfig` you register into the `simulator.ConfigGenerator` via `.SetPartition` defines its own width of partition of the global state data via `InitStateValues` and mutates this data via its own logic in its own concrete implementation of the `simulator.Iteration` interface, configured in the `simulator.PartitionConfig.Iteration`. Each partition mutates its own state but can read-only view any other partition's state history window in the `stateHistories` input into the `simulator.Iteration.Iterate` method call. Note that the `parititonIndex` is used to find which state history belongs to the partition being iterated and the order in which partitions are registered into `simulator.ConfigGenerator` determines their `partitionIndex`.

Let's look at what each field in the `simulator.PartitionConfig` does:

```
type PartitionConfig struct {
    Name               string  // This is the name for the partition (it's index is determined by the registration order into simulator.ConfigGenerator)
    Iteration          Iteration  // This is the implementation of this partition's state mutation iteration
    Params             Params // These are arbitrary parameters which can be configured at the start of the simulation
    ParamsAsPartitions map[string][]string // These are optional parameters which join into the parameters above by mapping partition names onto their corresponding partition indices - used only to conveniently set more human-readable partition name parameters but not used internally by the simulation, which always uses partitionIndex ordering 
    ParamsFromUpstream map[string]NamedUpstreamConfig // These are parameters which, for each simulation step, are determined to be the state values []float64 produced by any computationally upstream partition iterations and can hence be used to construct arbitrary fixed computational graph structures to coordinate partitions for simulation steps
    InitStateValues    []float64 // These are the initial state values which this partition subsequently mutates and also determines the width of this partition's state 
    StateHistoryDepth  int // This is the length/depth of the state history rolling window data kept to share as read-only between all partitions for all time - for example, index 0 of the `simulator.StateHistory.Values.RawRowView(0)` is the most recent state values for that partition from the last whole simulation step
    Seed               int // This is the seed input for any random number generation that may be needed 
}
```

The simulation run overall is configured using the `simulator.SimulationConfig` that is registered into the `simulator.ConfigGenerator` via `.SetSimulation`. Let's now look at what each field in the `simulator.SimulationConfig` does:

```
type SimulationConfig struct {
    OutputCondition      OutputCondition // This defines when the OutputFunction is called for each partition during a simulation run - for dexetera games we always configure this to be the OnlyNamesCondition
    OutputFunction       OutputFunction // This defines how simulation information is output during a run - for dexetera games we always configure this to be the JsCallbackOutputFunction so that the simulation can update the world view state of the game (see `RegisterStep` in pkg/simio/step.go to see how all of this works)
    TerminationCondition TerminationCondition // This defines when the simulation should end
    TimestepFunction     TimestepFunction // This defines how the global simulation time variable should be updated for each step
    InitTimeValue        float64 // The defines the initial simulation time
}
```

Please now read these docs which contain all of the useful pre-written partition iterations you might find useful to build simulations with:

- https://umbralcalc.github.io/stochadex/pkg/continuous.html
- https://umbralcalc.github.io/stochadex/pkg/discrete.html
- https://umbralcalc.github.io/stochadex/pkg/general.html
- https://umbralcalc.github.io/stochadex/pkg/kernels.html

You can also use these docs in the list above to help guide you on best practices for creating partition iterations of your own, and how to test them.


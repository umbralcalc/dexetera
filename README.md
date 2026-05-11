# dexetera

A Go module for turning [stochadex](https://github.com/umbralcalc/stochadex) simulations into self-contained, embeddable, in-browser interactive dashboards.

If you have a stochadex simulation and want to drop a live, interactive widget of it into a markdown blog post — sliders, buttons, real-time visualisation, the simulation actually running in the reader's browser via WebAssembly — you `go get` dexetera, write a tiny `dashboard.Config` describing the controls, and run a one-line codegen step. The output is a copy-pasteable HTML snippet.

The [growth](growth/) folder in this repo is a working end-to-end example you can preview locally and read as a template.

## Workflow for a downstream stochadex project

Suppose you have a project at `github.com/you/myblog-sims` that hosts your stochadex simulations and feeds widgets into your blog at `github.com/you/yourblog`.

### 1. Add dexetera as a dependency

```bash
cd ~/code/myblog-sims
go get github.com/umbralcalc/dexetera
```

### 2. Express your simulation as a `dashboard.Config`

Write a constructor in your project that returns a `*dashboard.Config`. Declare which partition states get streamed out, which take action input, the canvas visualization, the controls (sliders, readouts, optional reset button), the action driver, and the stochadex simulation builder. See [pkg/growth/growth.go](pkg/growth/growth.go) for the full pattern; the relevant builder calls look like:

```go
import (
    "github.com/umbralcalc/dexetera/pkg/dashboard"
    "github.com/umbralcalc/stochadex/pkg/simulator"
)

func NewConfig() *dashboard.Config {
    visConfig := dashboard.NewVisualizationBuilder().
        WithCanvas(320, 160).
        WithBackground("#ffffff").
        AddLineChart("population", 18, 18, 284, 124, &dashboard.ChartOptions{
            Color: "#3c78d8", LineWidth: 2,
        }).
        Build()

    return dashboard.NewConfigBuilder("growth").
        WithDescription("Logistic growth: drag the sliders to set r and K live.").
        WithServerPartition("population").
        WithActionStatePartition("population").
        WithVisualization(visConfig).
        WithSimulation(BuildGrowthSimulation).
        WithSlider(dashboard.Slider{
            Name: "r", Label: "r (growth rate)",
            Partition: "population", ValueIndex: 0,
            Min: 0, Max: 0.2, Step: 0.005, Default: 0.05, Decimals: 3,
        }).
        WithSlider(dashboard.Slider{
            Name: "K", Label: "K (carrying capacity)",
            Partition: "population", ValueIndex: 1,
            Min: 0, Max: 1000, Step: 10, Default: 500, Decimals: 0,
        }).
        WithReadout(dashboard.Readout{
            Partition: "population",
            Template:  "t = {t} · N = {v}",
        }).
        WithResetButton().
        WithInlineDriver(50). // tick interval (ms)
        Build()
}
```

### 3. Add a wasm entry point under `cmd/<name>/register_step/`

Five lines: hand the Config to `simio.RegisterStep`. See [cmd/growth/register_step/register_step.go](cmd/growth/register_step/register_step.go).

```go
//go:build js && wasm
package main

import (
    "github.com/umbralcalc/dexetera/pkg/simio"
    "github.com/you/myblog-sims/pkg/foo"
)

func main() {
    simio.RegisterStep(foo.NewConfig())
}
```

### 4. Add a codegen entry point under `cmd/<name>/generate/`

Five lines: call `dashboard.MustGenerateWidget`. See [cmd/growth/generate/generate.go](cmd/growth/generate/generate.go).

```go
package main

import (
    "github.com/umbralcalc/dexetera/pkg/dashboard"
    "github.com/you/myblog-sims/pkg/foo"
)

func main() {
    dashboard.MustGenerateWidget(foo.NewConfig(), dashboard.WidgetOptions{
        // The URLs the embedded snippet will use. Defaults work for
        // local preview via test.html; override for blog embedding.
        RuntimeBaseURL: "/assets/dexetera/runtime/",
        WasmURL:        "/assets/dexetera/widgets/foo/main.wasm",
    })
}
```

### 5. Generate the widget and build the wasm

```bash
go run ./cmd/foo/generate     # writes foo/widget.html, foo/test.html, foo/build.sh
./foo/build.sh                # compiles foo/src/main.wasm
```

### 6. Preview locally

```bash
python3 -m http.server 8000
# open http://localhost:8000/foo/test.html
```

The `test.html` wrapper uses local relative paths regardless of `WidgetOptions`, so it always works from a static-file server pointing at the repo root.

### 7. Embed in your blog

Two things need to happen on the blog repo's side:

**Sync the runtime once.** Copy this repo's `runtime/` folder to your blog's static assets, e.g.:

```bash
cp -r /path/to/dexetera/runtime/ /path/to/yourblog/assets/dexetera/runtime/
```

Then re-run that whenever you update dexetera. (One-line in a Makefile.)

**Sync each widget's wasm.** For each widget your blog uses, copy its built wasm:

```bash
mkdir -p /path/to/yourblog/assets/dexetera/widgets/foo
cp foo/src/main.wasm /path/to/yourblog/assets/dexetera/widgets/foo/main.wasm
```

**Paste the snippet into your post.** Open `foo/widget.html` and paste its contents into your markdown post (in a raw-HTML block — for Pandoc/Jekyll-with-Pandoc that's a fenced ```` ```{=html} ```` block). The snippet is a single `<div>` plus a scoped `<style>` plus an IIFE `<script>` — drop in anywhere prose flows, no further setup.

The snippet expects the runtime + wasm to be reachable at the URLs you set in step 4. Adjust those URLs to match your blog's asset layout.

## What the snippet looks like

The generated `widget.html` is one self-contained block:

- **`<div id="dexetera-foo" class="dexetera-widget">`** — the widget root. The id is unique per widget so multiple widgets can coexist on a page.
- **`<style>`** — all CSS scoped to `#dexetera-foo`. Won't bleed into the host page; multiple dexetera widgets on the same page won't fight.
- **The dashboard layout** — a panel grid (canvas + readouts on one panel, sliders + reset button on another).
- **`<script>`** — IIFE that loads `runtime/renderer.js` (deduplicated across widgets via a shared promise on `self.__dexeteraLoading`), spawns a Web Worker pointing at `runtime/worker.js`, and wires up sliders → `setActions` → wasm → renderer → DOM readouts.

## Action drivers

Per-step action input flows through one of two drivers, picked by the Config:

- **`inline`** (`WithInlineDriver(intervalMs)`) — actions come from in-page UI (slider events, button clicks). Tick rate is configurable. **This is what blog widgets typically use.**
- **`websocket`** (`WithWebsocketDriver(url)`) — actions come from an external WebSocket, e.g. a [dexact](https://pypi.org/project/dexact/) Python server running an `ActionTaker.take_next_action(time, states) -> list[float]`. Useful for offline experiments or when the action logic doesn't belong in the browser.

Adding a third driver (e.g. one that replays actions from a recorded log) is a self-contained ~50-line file under [runtime/drivers/](runtime/drivers/). Each driver defines `self.createDriver(env, options)` and the worker dynamically loads whichever one the Config asks for.

## Action delivery: per-partition named vs. broadcast

Per-step action input flows through an `ActionState` protobuf message ([proto/action_state.proto](proto/action_state.proto)) with two delivery paths:

- **Per-partition named** (`partitions` map): each entry writes to the partition whose name matches the map key. The path the inline driver uses, and the path most dashboards want.
- **Broadcast** (`values` slice): the same slice is delivered to every partition listed in `ActionStatePartitionNames`. Retained as a wire-compatibility shim for existing dexact Python clients.

The named path takes precedence when both are present. See [pkg/simio/dispatch.go](pkg/simio/dispatch.go) for the full semantics.

## Repo layout

```
pkg/dashboard/        Config, ConfigBuilder, VisualizationBuilder,
                      WidgetOptions, GenerateWidget — the public Go API
pkg/simio/            Wasm-side runtime: RegisterStep + ApplyActionState
pkg/growth/           The end-to-end smoke-test simulation
cmd/growth/
    register_step/    Wasm main for growth (template for your projects)
    generate/         Codegen main for growth (template for your projects)
runtime/              JS runtime — sync this folder into your blog's
                      static assets, once. Contains renderer.js,
                      worker.js, the proto stubs, drivers/.
proto/                action_state.proto + regen script
growth/               growth's generated widget + local-preview wrapper
                      (safe to delete; regenerate via `go run ./cmd/growth/generate`)
```

## Regenerating proto stubs

```bash
./proto/generate_proto.sh
```

Writes Go output to `pkg/simio/action_state.pb.go` and JS output to `runtime/action_state_pb.js`. Requires `protoc` and `protoc-gen-go` on `$PATH`.

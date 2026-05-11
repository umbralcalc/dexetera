package dashboard

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

// WidgetOptions controls where the generated widget snippet expects its
// runtime files and wasm binary to be hosted, plus the output paths.
//
// All fields are optional; sensible defaults make the generated output
// runnable from a static-file server pointing at the parent directory
// (the "test.html" wrapper uses local paths regardless, so local preview
// works out of the box).
type WidgetOptions struct {
	// OutputDir is where the codegen writes widget.html, test.html, and
	// build.sh. Defaults to "<config.Name>/" if empty.
	OutputDir string

	// WidgetID is the unique HTML id of the widget root element, used to
	// scope all CSS rules so multiple widgets on a page don't collide.
	// Defaults to "dexetera-<config.Name>" if empty. Must be a valid CSS id.
	WidgetID string

	// RuntimeBaseURL is the URL prefix where the embedding page can find
	// the dexetera runtime files (renderer.js, worker.js, wasm_exec.js,
	// proto stubs, drivers/). Trailing slash recommended.
	//
	// For embedding in a Jekyll-style blog, set this to wherever you've
	// copied the dexetera runtime/ folder, typically something like
	// "/assets/dexetera/runtime/". Defaults to "./runtime/".
	RuntimeBaseURL string

	// WasmURL is the URL of the wasm binary on the embedding page's
	// domain. For blog embedding, set to wherever you've copied the
	// built wasm, e.g. "/assets/dexetera/widgets/<name>/main.wasm".
	// Defaults to "./src/main.wasm".
	WasmURL string
}

func (o *WidgetOptions) applyDefaults(name string) {
	if o.OutputDir == "" {
		o.OutputDir = name + "/"
	}
	if o.WidgetID == "" {
		o.WidgetID = "dexetera-" + name
	}
	if o.RuntimeBaseURL == "" {
		o.RuntimeBaseURL = "./runtime/"
	}
	if o.WasmURL == "" {
		o.WasmURL = "./src/main.wasm"
	}
}

// GenerateWidget writes the embeddable widget snippet (widget.html), a
// local-preview wrapper (test.html), and a wasm build script (build.sh)
// for the given Config + options. The output files form a self-contained
// preview locally and a copy-pasteable snippet for blog embedding.
//
// Outputs under opts.OutputDir:
//
//	widget.html  Drop-in <div>+<style>+<script> snippet for embedding.
//	             References RuntimeBaseURL and WasmURL absolutely so the
//	             embedding page only needs to host the runtime + wasm.
//	test.html    Standalone HTML page that wraps the widget for local
//	             preview. Uses local relative paths (../runtime/, ./src/...)
//	             regardless of opts so that `python3 -m http.server` from
//	             the repo root can preview without further configuration.
//	build.sh     Script that compiles cmd/<name>/register_step to
//	             src/main.wasm.
//
// The output directory is created if it doesn't exist; existing files
// in it are overwritten.
func GenerateWidget(config *Config, opts WidgetOptions) error {
	opts.applyDefaults(config.Name)

	if err := os.MkdirAll(opts.OutputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	body, err := renderWidgetBody(config, opts.WidgetID, opts.RuntimeBaseURL, opts.WasmURL)
	if err != nil {
		return fmt.Errorf("failed to render widget body: %w", err)
	}
	if err := writeFile(opts.OutputDir, "widget.html", body); err != nil {
		return err
	}

	// test.html: same widget body, but with local relative URLs so the
	// page works via `python3 -m http.server` from the repo root.
	testBody, err := renderWidgetBody(config, opts.WidgetID, "../runtime/", "./src/main.wasm")
	if err != nil {
		return fmt.Errorf("failed to render test widget body: %w", err)
	}
	testPage := wrapTestHTML(config.Name, testBody)
	if err := writeFile(opts.OutputDir, "test.html", testPage); err != nil {
		return err
	}

	if err := generateBuildScript(opts.OutputDir, config.Name); err != nil {
		return err
	}
	return nil
}

func writeFile(dir, name, contents string) error {
	path := filepath.Join(dir, name)
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create %s: %w", name, err)
	}
	defer f.Close()
	if _, err := f.WriteString(contents); err != nil {
		return fmt.Errorf("failed to write %s: %w", name, err)
	}
	return nil
}

// renderWidgetBody produces the embeddable widget snippet — a single
// <div> wrapping a scoped <style> block, the panel layout, and an IIFE
// <script> that loads the runtime, instantiates a worker, and wires up
// the canvas + sliders + readouts + reset button.
//
// All CSS selectors are prefixed with "#<widgetID>" so the styles stay
// confined to this widget — multiple dexetera widgets can coexist on the
// same page without fighting over .panel, .slider, etc.
func renderWidgetBody(cfg *Config, widgetID, runtimeBase, wasmURL string) (string, error) {
	visConfig := cfg.VisualizationConfig
	hasControls := len(cfg.Sliders) > 0 || cfg.ShowReset

	// Marshal the renderer / sliders / readouts / driver as JSON so the
	// widget script reads them as a plain object literal — same pattern
	// the previous codegen used, just inlined now.
	cfgJSON, err := marshalGameConfig(cfg)
	if err != nil {
		return "", err
	}

	data := struct {
		WidgetID       string
		RuntimeBase    string
		WasmURL        string
		Description    string
		CanvasWidth    int
		CanvasHeight   int
		Sliders        []Slider
		Readouts       []Readout
		HasControls    bool
		ShowReset      bool
		GameConfigJSON string
	}{
		WidgetID:       widgetID,
		RuntimeBase:    runtimeBase,
		WasmURL:        wasmURL,
		Description:    cfg.Description,
		CanvasWidth:    visConfig.CanvasWidth,
		CanvasHeight:   visConfig.CanvasHeight,
		Sliders:        cfg.Sliders,
		Readouts:       cfg.Readouts,
		HasControls:    hasControls,
		ShowReset:      cfg.ShowReset,
		GameConfigJSON: cfgJSON,
	}

	tmpl, err := template.New("widget").Parse(widgetTemplate)
	if err != nil {
		return "", err
	}
	var buf stringBuilder
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// stringBuilder is a tiny os.File-shaped wrapper around a strings.Builder
// so text/template can write into it.
type stringBuilder struct{ b []byte }

func (s *stringBuilder) Write(p []byte) (int, error) { s.b = append(s.b, p...); return len(p), nil }
func (s *stringBuilder) String() string              { return string(s.b) }

// widgetTemplate is the embeddable widget snippet. It's structured as
// one wrapping <div id="{{.WidgetID}}">, an inline scoped <style>, the
// panel layout, and an IIFE <script>. The script loads renderer.js from
// RuntimeBase (deduplicating across widgets), then spins up its own
// Worker pointing at RuntimeBase/worker.js with the gameConfig.driver.
const widgetTemplate = `<div id="{{.WidgetID}}" class="dexetera-widget">
<style>
#{{.WidgetID}} { font-family: system-ui, -apple-system, sans-serif; color: #2c3e50; line-height: 1.5; }
#{{.WidgetID}} .description { margin: 0 0 1em; color: #2c3e50; opacity: 0.85; font-size: 1rem; }
#{{.WidgetID}} code { font-family: ui-monospace, SFMono-Regular, Menlo, monospace; font-size: 0.95em; background: rgba(60,120,216,0.08); padding: 0.05em 0.3em; border-radius: 3px; }
#{{.WidgetID}} .dashboard { display: grid; grid-template-columns: repeat(auto-fit, minmax(320px, 1fr)); gap: 0.9em; }
#{{.WidgetID}} .panel { border: 1px solid #2c3e50; border-radius: 6px; padding: 0.8em 0.9em; background: #ffffff; display: flex; flex-direction: column; gap: 0.6em; box-sizing: border-box; }
#{{.WidgetID}} .panel-title { font-weight: 600; color: #2c3e50; font-size: 1rem; }
#{{.WidgetID}} canvas { display: block; width: 100%; max-width: {{.CanvasWidth}}px; height: auto; aspect-ratio: {{.CanvasWidth}} / {{.CanvasHeight}}; margin: 0 auto; background: #ffffff; }
#{{.WidgetID}} .panel-readout { margin: 0; font-size: 1rem; color: #2c3e50; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
#{{.WidgetID}} .slider { display: grid; grid-template-columns: 1fr 80px; grid-template-areas: "name readout" "input input"; align-items: center; gap: 0.3em 1em; font-size: 1rem; }
#{{.WidgetID}} .slider-name { grid-area: name; color: #2c3e50; }
#{{.WidgetID}} .slider input[type="range"] { grid-area: input; width: 100%; accent-color: #3c78d8; }
#{{.WidgetID}} .slider-readout { grid-area: readout; text-align: right; color: #3c78d8; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
#{{.WidgetID}} .panel-actions { display: flex; flex-wrap: wrap; gap: 0.6em; margin-top: 0.4em; }
#{{.WidgetID}} button.button-secondary { cursor: pointer; border: 1px solid #2c3e50; background: #ffffff; color: #2c3e50; padding: 0.4em 0.85em; border-radius: 6px; font-size: 1rem; font-family: inherit; }
#{{.WidgetID}} button.button-secondary:hover { background: #f4f6f9; }
#{{.WidgetID}} .status { margin: 1em 0 0; text-align: right; font-size: 0.9em; color: #2c3e50; opacity: 0.6; font-family: ui-monospace, SFMono-Regular, Menlo, monospace; }
</style>
{{if .Description}}<p class="description">{{.Description}}</p>{{end}}
<div class="dashboard">
    <section class="panel">
        <div class="panel-title">Simulation</div>
        <canvas width="{{.CanvasWidth}}" height="{{.CanvasHeight}}"></canvas>
        {{range .Readouts}}
        <p class="panel-readout" data-readout="{{.Partition}}">&nbsp;</p>
        {{end}}
    </section>
    {{if .HasControls}}
    <section class="panel">
        <div class="panel-title">Live controls</div>
        {{range .Sliders}}
        <label class="slider">
            <span class="slider-name">{{.Label}}</span>
            <input type="range" data-slider="{{.Name}}"
                   min="{{.Min}}" max="{{.Max}}" step="{{.Step}}" value="{{.Default}}">
            <span class="slider-readout" data-slider-readout="{{.Name}}">&nbsp;</span>
        </label>
        {{end}}
        {{if .ShowReset}}
        <div class="panel-actions">
            <button type="button" class="button-secondary" data-reset>Reset simulation</button>
        </div>
        {{end}}
    </section>
    {{end}}
</div>
<p class="status" data-status>Loading…</p>
</div>
<script>
(function () {
    var widget = document.getElementById('{{.WidgetID}}');
    var RUNTIME_BASE = '{{.RuntimeBase}}';
    var WASM_URL = '{{.WasmURL}}';
    var gameConfig = {{.GameConfigJSON}};

    // Load the renderer script lazily, sharing one promise across every
    // dexetera widget on the same page.
    function ensureRenderer() {
        if (self.dexetera && self.dexetera.createRenderer) return Promise.resolve();
        if (self.__dexeteraLoading) return self.__dexeteraLoading;
        self.__dexeteraLoading = new Promise(function (resolve, reject) {
            var s = document.createElement('script');
            s.src = RUNTIME_BASE + 'renderer.js';
            s.onload = function () { resolve(); };
            s.onerror = function () { reject(new Error('failed to load ' + s.src)); };
            document.head.appendChild(s);
        });
        return self.__dexeteraLoading;
    }

    function $(sel) { return widget.querySelector(sel); }
    function $$(sel) { return widget.querySelectorAll(sel); }

    function setStatus(msg) {
        var el = $('[data-status]');
        if (el) el.textContent = msg;
    }

    var slidersByPartition = (function () {
        var grouped = {};
        for (var i = 0; i < gameConfig.sliders.length; i++) {
            var s = gameConfig.sliders[i];
            if (!grouped[s.partition]) grouped[s.partition] = [];
            grouped[s.partition].push(s);
        }
        return grouped;
    })();

    var worker = null;

    function publishActions() {
        for (var i = 0; i < gameConfig.sliders.length; i++) {
            var s = gameConfig.sliders[i];
            var input = $('[data-slider="' + s.name + '"]');
            var ro = $('[data-slider-readout="' + s.name + '"]');
            if (input && ro) ro.textContent = parseFloat(input.value).toFixed(s.decimals);
        }
        if (!worker) return;
        var partitions = {};
        for (var partition in slidersByPartition) {
            if (!Object.prototype.hasOwnProperty.call(slidersByPartition, partition)) continue;
            var group = slidersByPartition[partition];
            var maxIdx = -1;
            for (var j = 0; j < group.length; j++) maxIdx = Math.max(maxIdx, group[j].valueIndex);
            var values = new Array(maxIdx + 1);
            for (var k = 0; k <= maxIdx; k++) values[k] = 0;
            for (var l = 0; l < group.length; l++) {
                var sl = group[l];
                var inp = $('[data-slider="' + sl.name + '"]');
                values[sl.valueIndex] = parseFloat(inp ? inp.value : sl.default);
            }
            partitions[partition] = values;
        }
        worker.postMessage({ action: 'setActions', partitions: partitions });
    }

    function applyReadout(template, decimals, partitionState) {
        var s = template;
        s = s.replace(/\{t\}/g, Math.floor(partitionState.timesteps));
        s = s.replace(/\{v(\d*)\}/g, function (_, idx) {
            var i = idx === '' ? 0 : parseInt(idx, 10);
            var v = partitionState.state.values[i];
            return (v === undefined) ? '' : v.toFixed(decimals);
        });
        return s;
    }

    function startWorker(renderer) {
        if (worker) worker.terminate();
        worker = new Worker(RUNTIME_BASE + 'worker.js');
        worker.onmessage = function (e) {
            var msg = e.data;
            if (msg.type === 'partitionState') {
                renderer.update(msg.data);
                renderer.render();
                for (var i = 0; i < gameConfig.readouts.length; i++) {
                    var r = gameConfig.readouts[i];
                    if (r.partition !== msg.data.partitionName) continue;
                    var el = $('[data-readout="' + r.partition + '"]');
                    if (el) el.textContent = applyReadout(r.template, r.decimals, msg.data);
                }
            } else if (msg.type === 'status') {
                setStatus(msg.data);
            } else if (msg.type === 'error') {
                console.error('dexetera worker error:', msg.data);
                setStatus('Error: ' + msg.data);
            }
        };
        worker.onerror = function (err) {
            console.error('dexetera worker error:', err);
            setStatus('Worker error: ' + err.message);
        };
        worker.postMessage({
            action: 'start',
            wasmBinary: new URL(WASM_URL, document.baseURI).href,
            driver: gameConfig.driver,
        });
        publishActions();
    }

    ensureRenderer().then(function () {
        var canvas = $('canvas');
        var renderer = self.dexetera.createRenderer(canvas, gameConfig.visualization);

        for (var i = 0; i < gameConfig.sliders.length; i++) {
            var s = gameConfig.sliders[i];
            var el = $('[data-slider="' + s.name + '"]');
            if (el) el.addEventListener('input', publishActions);
        }
        if (gameConfig.showReset) {
            var btn = $('[data-reset]');
            if (btn) btn.addEventListener('click', function () { startWorker(renderer); });
        }
        publishActions();
        startWorker(renderer);
    }).catch(function (err) {
        console.error(err);
        setStatus('Failed to load dexetera runtime: ' + err.message);
    });
})();
</script>
`

// wrapTestHTML wraps a widget body in a minimal standalone HTML page so
// that opening test.html in a browser (served via a static-file server)
// previews the widget locally without any blog setup.
func wrapTestHTML(name, widgetBody string) string {
	return `<!DOCTYPE html>
<html lang="en">
<head>
<meta charset="UTF-8">
<meta name="viewport" content="width=device-width, initial-scale=1.0">
<title>dexetera — ` + name + ` (local preview)</title>
<style>
body { margin: 0; padding: 2em 1.5em; max-width: 900px; margin: 0 auto;
       font-family: system-ui, -apple-system, sans-serif; color: #2c3e50; }
h1 { font-size: 2.4em; font-weight: 600; letter-spacing: -0.02em; margin: 0 0 0.6em; }
header p { margin: 0 0 1.4em; opacity: 0.65; font-size: 0.95em; }
</style>
</head>
<body>
<header>
<h1>` + name + `</h1>
<p>Local preview. The embeddable snippet is in widget.html.</p>
</header>
` + widgetBody + `
</body>
</html>
`
}

// jsConfig and friends are the JSON shape the inline widget script reads.
// They mirror Slider / Readout / DriverSpec but with lowercase JSON tags
// so the script can index them naturally as plain JS objects.

type jsSlider struct {
	Name       string  `json:"name"`
	Partition  string  `json:"partition"`
	ValueIndex int     `json:"valueIndex"`
	Default    float64 `json:"default"`
	Decimals   int     `json:"decimals"`
}

type jsReadout struct {
	Partition string `json:"partition"`
	Template  string `json:"template"`
	Decimals  int    `json:"decimals"`
}

type jsConfig struct {
	Visualization map[string]interface{} `json:"visualization"`
	Sliders       []jsSlider             `json:"sliders"`
	Readouts      []jsReadout            `json:"readouts"`
	ShowReset     bool                   `json:"showReset"`
	Driver        map[string]interface{} `json:"driver"`
}

func marshalGameConfig(cfg *Config) (string, error) {
	visConfig := cfg.VisualizationConfig
	renderers := make([]map[string]interface{}, 0, len(visConfig.Renderers))
	for _, r := range visConfig.Renderers {
		renderers = append(renderers, map[string]interface{}{
			"type":          r.Type,
			"partitionName": r.PartitionName,
			"properties":    r.Properties,
		})
	}

	sliders := make([]jsSlider, 0, len(cfg.Sliders))
	for _, s := range cfg.Sliders {
		decimals := s.Decimals
		if decimals == 0 {
			decimals = 3
		}
		sliders = append(sliders, jsSlider{
			Name:       s.Name,
			Partition:  s.Partition,
			ValueIndex: s.ValueIndex,
			Default:    s.Default,
			Decimals:   decimals,
		})
	}
	readouts := make([]jsReadout, 0, len(cfg.Readouts))
	for _, r := range cfg.Readouts {
		decimals := r.Decimals
		if decimals == 0 {
			decimals = 2
		}
		readouts = append(readouts, jsReadout{
			Partition: r.Partition,
			Template:  r.Template,
			Decimals:  decimals,
		})
	}

	driverOpts := cfg.Driver.Options
	if driverOpts == nil {
		driverOpts = map[string]interface{}{}
	}

	jc := jsConfig{
		Visualization: map[string]interface{}{
			"canvasWidth":      visConfig.CanvasWidth,
			"canvasHeight":     visConfig.CanvasHeight,
			"backgroundColor":  visConfig.BackgroundColor,
			"updateIntervalMs": visConfig.UpdateIntervalMs,
			"renderers":        renderers,
		},
		Sliders:   sliders,
		Readouts:  readouts,
		ShowReset: cfg.ShowReset,
		Driver: map[string]interface{}{
			"kind":    cfg.Driver.Kind,
			"options": driverOpts,
		},
	}
	out, err := json.Marshal(jc)
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// generateBuildScript writes a per-widget build.sh that compiles the
// example's wasm into <outputDir>/src/main.wasm. Same shape as before.
func generateBuildScript(outputDir, name string) error {
	const buildScript = `#!/usr/bin/env bash
# Generated build script for {{.Name}}.
# Compiles cmd/{{.Name}}/register_step to src/main.wasm.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"

mkdir -p "$SCRIPT_DIR/src"
cd "$PROJECT_ROOT"
GOOS=js GOARCH=wasm go build -o "$SCRIPT_DIR/src/main.wasm" ./cmd/{{.Name}}/register_step
echo "Built $SCRIPT_DIR/src/main.wasm"
`
	tmpl, err := template.New("build").Parse(buildScript)
	if err != nil {
		return err
	}
	data := struct{ Name string }{Name: name}
	path := filepath.Join(outputDir, "build.sh")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	if err := tmpl.Execute(file, data); err != nil {
		return err
	}
	return os.Chmod(path, 0755)
}

// MustGenerateWidget is the panic-on-error convenience wrapper used by
// the cmd/<name>/generate mains: most of those programs do nothing but
// call this and would otherwise just `if err != nil { panic(err) }` it.
func MustGenerateWidget(config *Config, opts WidgetOptions) {
	if err := GenerateWidget(config, opts); err != nil {
		log.Fatalf("dashboard: %v", err)
	}
}

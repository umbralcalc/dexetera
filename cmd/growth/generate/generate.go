// generate writes the growth example's embeddable widget snippet
// (widget.html), local-preview wrapper (test.html), and wasm build
// script (build.sh) into the growth/ folder. Re-run whenever the
// dashboard.Config in pkg/growth or the codegen templates change.
//
// To embed the snippet in a Jekyll-style blog post, configure
// RuntimeBaseURL and WasmURL to the absolute URLs the blog hosts the
// runtime and wasm at, e.g. "/assets/dexetera/runtime/" and
// "/assets/dexetera/widgets/growth/main.wasm". The defaults below leave
// them empty — the codegen falls back to "./runtime/" + "./src/main.wasm",
// which keeps the local-preview test.html happy.
package main

import (
	"github.com/umbralcalc/dexetera/pkg/dashboard"
	"github.com/umbralcalc/dexetera/pkg/growth"
)

func main() {
	dashboard.MustGenerateWidget(growth.NewConfig(), dashboard.WidgetOptions{})
}

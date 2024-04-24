package main

import (
	"net/http"

	"github.com/aim4ik11/architecture-lab-3/painter"
	"github.com/aim4ik11/architecture-lab-3/painter/lang"
	"github.com/aim4ik11/architecture-lab-3/ui"
)

func main() {
	var (
		pv ui.Visualizer

		opLoop painter.Loop
		parser lang.Parser
	)

	pv.Title = "Simple painter"

	pv.OnScreenReady = opLoop.Start
	opLoop.Receiver = &pv

	go func() {
		http.Handle("/", lang.HttpHandler(&opLoop, &parser))
		_ = http.ListenAndServe("localhost:17000", nil)
	}()

	pv.Main()
	opLoop.StopAndWait()
}

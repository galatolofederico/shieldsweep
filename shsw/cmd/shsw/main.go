package main

import (
	"flag"

	"github.com/galatolofederico/shieldsweep/shsw/internal/engine"
)

func main() {
	var home string

	flag.StringVar(&home, "home", "/etc/shsw", "ShieldSweep home directory (where shsw.json is located)")
	flag.Parse()

	engine := engine.NewEngine(home)
	engine.Run()
}

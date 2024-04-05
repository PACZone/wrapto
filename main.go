package main

import "github.com/PACZone/wrapto/core"

func main() {
	core, err := core.NewCore()
	if err != nil {
		panic(err.Error()) // TODO: Log
	}

	core.Start()
}

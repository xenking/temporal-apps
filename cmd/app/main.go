package main

import (
	"github.com/cristalhq/acmd"
)

var commands []acmd.Command

func main() {
	r := acmd.RunnerOf(commands, acmd.Config{
		AppName:        "temporal-apps",
		AppDescription: "Temporal Apps CLI",
		Version:        "0.1.0",
	})

	if err := r.Run(); err != nil {
		r.Exit(err)
	}
}

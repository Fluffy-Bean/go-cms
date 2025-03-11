package main

import (
	"github.com/Fluffy-Bean/cms/cmd"
)

func main() {
	conf := cmd.Config{
		Host: "0.0.0.0:7070",
	}

	cmd.Execute(conf)
}

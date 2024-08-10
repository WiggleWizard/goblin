package main

import (
	"goblin/cmd"
	"log"
	"os"
)

func main() {
	// Golang, why do I have to force timestamps off and force the log output to be stdout? Why isn't this the default?
	// Go can goblin. Goblin on deez nutz
	log.SetFlags(0)
	log.SetOutput(os.Stdout)

	cmd.Execute()
}

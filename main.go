package main

import (
	"log"
	"os"

	"github.com/starkandwayne/molten-core/commands"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	logger := log.New(os.Stdout, "", 0)

	app := kingpin.New("mc", "MoltenCore Cli")
	commands.Configure(logger, app)
	kingpin.MustParse(app.Parse(os.Args[1:]))
}

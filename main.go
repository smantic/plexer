package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/smantic/plexer/cmd"
	"github.com/smantic/plexer/cmd/listen"
)

func main() {

	flag.Usage = func() {
		fmt.Printf(cmd.HelpStr)
	}

	if len(os.Args) == 1 {
		fmt.Printf(cmd.HelpStr)
		return
	}

	switch os.Args[1] {
	case "help":
		fmt.Printf(cmd.HelpStr)
		return
	case "listen":
		listen.Run(os.Args[2:])
	}

	flag.Parse()
}

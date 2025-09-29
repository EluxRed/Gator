package main

import (
	"fmt"
	"log"
	"os"

	"github.com/EluxRed/Gator/internal/config"
)

type state struct {
	configPtr *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	current_state := &state{configPtr: &cfg}

	cmds := commands{registeredCommands: make(map[string]func(*state, command) error)}

	if len(os.Args) < 2 {
		log.Fatalln(fmt.Errorf("command name missing"))
	}
	cmd_login := command{Name: os.Args[1], Args: os.Args[2:]}

	cmds.register("login", handlerLogin)

	if err := cmds.run(current_state, cmd_login); err != nil {
		log.Fatalln(err)
	}
}

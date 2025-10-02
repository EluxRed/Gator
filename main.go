package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/EluxRed/Gator/internal/config"
	"github.com/EluxRed/Gator/internal/database"

	_ "github.com/lib/pq" //The underscore tells Go that you're importing it for its side effects, not because you need to use it.
)

type state struct {
	dbPtr     *database.Queries
	configPtr *config.Config
}

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := sql.Open("postgres", cfg.DB_URL)
	if err != nil {
		log.Fatalln(err)
	}

	dbQueries := database.New(db)

	current_state := &state{dbPtr: dbQueries, configPtr: &cfg}

	cmds := commands{registeredCommands: make(map[string]func(*state, command) error)}

	if len(os.Args) < 2 {
		log.Fatalln(fmt.Errorf("command name missing"))
	}
	cmd := command{Name: os.Args[1], Args: os.Args[2:]}

	switch cmd.Name {
	case "login":
		cmds.register(cmd.Name, handlerLogin)
	case "register":
		cmds.register(cmd.Name, handlerRegister)
	case "reset":
		cmds.register(cmd.Name, handlerReset)
	case "users":
		cmds.register(cmd.Name, handlerUsers)
	case "agg":
		cmds.register(cmd.Name, handlerAgg)
	case "addfeed":
		cmds.register(cmd.Name, middlewareLoggedIn(handlerAddFeed))
	case "feeds":
		cmds.register(cmd.Name, handlerFeeds)
	case "follow":
		cmds.register(cmd.Name, middlewareLoggedIn(handlerFollow))
	case "following":
		cmds.register(cmd.Name, middlewareLoggedIn(handlerFollowing))
	case "unfollow":
		cmds.register(cmd.Name, middlewareLoggedIn(handlerUnfollow))
	case "browse":
		cmds.register(cmd.Name, middlewareLoggedIn(handlerBrowse))
	default:
		log.Fatalln(fmt.Errorf("wrong command"))
	}

	if err := cmds.run(current_state, cmd); err != nil {
		log.Fatalln(err)
	}
}

func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		currentUsr, err := s.dbPtr.GetUser(context.Background(), s.configPtr.Current_User_Name)
		if err != nil {
			return err
		}
		return handler(s, cmd, currentUsr)
	}
}

func checkParams(args []string, num int) error {
	if len(args) != num {
		return fmt.Errorf("parameters passed: %v\nparameters expected: %v", len(args), num)
	}
	return nil
}

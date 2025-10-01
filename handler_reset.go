package main

import (
	"context"
)

func handlerReset(s *state, cmd command) error {
	if err := checkParams(cmd.Args, 0); err != nil {
		return err
	}
	if err := s.dbPtr.DeleteAllUsers(context.Background()); err != nil {
		return err
	}
	return nil
}

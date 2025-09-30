package main

import (
	"context"
	"fmt"
)

func handlerReset(s *state, cmd command) error {
	if len(cmd.Args) > 0 {
		return fmt.Errorf("parameters passed: %v\nparameters expected: %v", len(cmd.Args), 0)
	}
	if err := s.dbPtr.DeleteAllUsers(context.Background()); err != nil {
		return err
	}
	return nil
}

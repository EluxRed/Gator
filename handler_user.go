package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/EluxRed/Gator/internal/database"
	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if err := checkParams(cmd.Args, 1); err != nil {
		return err
	}
	name := cmd.Args[0]
	_, err := s.dbPtr.GetUser(context.Background(), name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user does not exist: %w", err)
		} else {
			return err
		}
	}
	if err := s.configPtr.SetUser(name); err != nil {
		return fmt.Errorf("couldn't set current user: %w", err)
	}
	fmt.Println("user set to:", name)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if err := checkParams(cmd.Args, 1); err != nil {
		return err
	}
	usr := database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(), Name: cmd.Args[0]}
	_, err := s.dbPtr.GetUser(context.Background(), usr.Name)
	if err != nil {
		if !errors.Is(err, sql.ErrNoRows) {
			return err
		}
	} else {
		return fmt.Errorf("username already exists")
	}
	createdUsr, err := s.dbPtr.CreateUser(context.Background(), usr)
	if err != nil {
		return err
	}
	s.configPtr.SetUser(createdUsr.Name)
	fmt.Println("the user was created")
	printUser(createdUsr)
	return nil
}

func handlerUsers(s *state, cmd command) error {
	if err := checkParams(cmd.Args, 0); err != nil {
		return err
	}
	users, err := s.dbPtr.GetUsers(context.Background())
	if err != nil {
		return err
	}
	currentUsr := s.configPtr.Current_User_Name
	for i := range users {
		if currentUsr == users[i] {
			users[i] = currentUsr + " (current)"
		}
		fmt.Println("* ", users[i])
	}

	return nil
}

func printUser(user database.User) {
	fmt.Printf(" * ID:      %v\n", user.ID)
	fmt.Printf(" * Name:    %v\n", user.Name)
}

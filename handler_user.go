package main

import (
	"fmt"
	"time"
	"github.com/mirmeep/gator-blog/internal/database"
	"context"
	"github.com/google/uuid"
	"database/sql"
)


func handlerLogin(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("command login requires 1 argument, username")
	}

	username := cmd.args[0]
	_, err := s.db.GetUser(context.Background(), sql.NullString{
			String: username,
			Valid: true,
		})
	if err != nil {
		return fmt.Errorf("failed to get user: %w", err)
	}

	err = s.conf.SetUser(username)
	if err != nil {
		return fmt.Errorf("cannot set username: %w", err)
	}

	fmt.Printf("user %s logged in\n", username)

	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.args) == 0 {
		return fmt.Errorf("command register requires 1 argument, name")
	}

	username := cmd.args[0]

	_, err := s.db.CreateUser(context.Background(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name: sql.NullString{
			String: username,
			Valid: true,
		},
	})

	if err != nil {
		return fmt.Errorf("user creation failed: %w", err)
	}

	s.conf.SetUser(username)
	fmt.Printf("User %s created!", username)

	return nil
}

func handlerReset(s *state, cmd command) error {
	if len(cmd.args) > 0 {
		return fmt.Errorf("command reset doesn't take params")
	}

	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not delete users: %w", err)
	}

	fmt.Println("Database users successfully deleted!")
	return nil
}
// Command seed populates a SQLite database with users, groups, and messages
// defined in cmd/seed/data/*.json, all written through the real
// UserService/GroupService — not raw DB inserts — so the seeded data
// exercises the same validation and side effects as the app itself.
//
// This command is deliberately theme-agnostic: it knows nothing about what
// content it's seeding, only the User/Group/Message shape. Swapping in a
// different theme (or richer content later pulled from an external source)
// is purely a matter of editing cmd/seed/data/*.json — no Go changes.
//
// Point DATABASE_PATH at a scratch file before running so this doesn't
// touch real data, and run from the repo root so the relative data paths
// resolve, e.g.:
//
//	DATABASE_PATH=./internal/db/seed.sample.sqlite3 go run ./cmd/seed
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/alextilot/golang-htmx-chatapp/config"
	"github.com/alextilot/golang-htmx-chatapp/internal/db"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
)

const (
	seedPassword = "Password123!" // every seeded user gets this password
	// File names match their model: internal/model/user.go / group.go.
	usersFile  = "cmd/seed/data/user.json"
	groupsFile = "cmd/seed/data/group.json"
)

type userSeed struct {
	Username string `json:"username"`
	Email    string `json:"email"`
}

type messageSeed struct {
	Sender  string `json:"sender"`
	Content string `json:"content"`
}

type groupSeed struct {
	Name     string        `json:"name"`
	Members  []string      `json:"members"` // usernames; first is the creator
	Messages []messageSeed `json:"messages"`
}

func loadJSON[T any](path string) (T, error) {
	var out T
	data, err := os.ReadFile(path)
	if err != nil {
		return out, err
	}
	if err := json.Unmarshal(data, &out); err != nil {
		return out, fmt.Errorf("parsing %s: %w", path, err)
	}
	return out, nil
}

func main() {
	ctx := context.Background()
	cfg := config.Load()

	userSeeds, err := loadJSON[[]userSeed](usersFile)
	if err != nil {
		log.Fatalf("load %s: %v", usersFile, err)
	}

	groupSeeds, err := loadJSON[[]groupSeed](groupsFile)
	if err != nil {
		log.Fatalf("load %s: %v", groupsFile, err)
	}

	database := db.NewDatabase(db.Config{DSN: cfg.DatabasePath})
	defer database.Close()
	database.AutoMigrate()

	repos := repository.NewRepositories(repository.Deps{DB: database.Conn})
	services := service.NewServices(service.Deps{Repos: repos})

	log.Printf("🌱 seeding %s / %s into %s", usersFile, groupsFile, cfg.DatabasePath)

	users := make(map[string]*model.User, len(userSeeds))
	for _, u := range userSeeds {
		created, err := services.UserService.Signup(ctx, service.SignupInput{
			Username:       u.Username,
			Email:          u.Email,
			Password:       seedPassword,
			RepeatPassword: seedPassword,
		})
		if err != nil {
			log.Fatalf("signup %s: %v", u.Username, err)
		}
		users[u.Username] = created
	}

	for _, g := range groupSeeds {
		creator := users[g.Members[0]]

		grp, err := services.GroupService.Create(ctx, creator.ID, service.CreateGroupInput{
			Name: g.Name,
			Type: model.GroupTypeGroup,
		})
		if err != nil {
			log.Fatalf("create group %s: %v", g.Name, err)
		}

		for _, member := range g.Members[1:] {
			u := users[member]
			if err := services.GroupService.AddMember(ctx, grp.ID, creator.ID, service.AddMemberInput{UserID: u.ID}); err != nil {
				log.Fatalf("add member %s to %s: %v", member, g.Name, err)
			}
		}

		for _, m := range g.Messages {
			sender := users[m.Sender]
			if _, err := services.GroupService.SendMessage(ctx, grp.ID, sender.ID, service.SendMessageInput{Content: m.Content}); err != nil {
				log.Fatalf("send message in %s: %v", g.Name, err)
			}
		}

		log.Printf("✅ seeded group %q with %d members, %d messages", g.Name, len(g.Members), len(g.Messages))
	}

	fmt.Println()
	fmt.Println("Seed complete. Log in as any seeded user with password:", seedPassword)
	for _, u := range userSeeds {
		fmt.Printf("  %s\n", u.Username)
	}
}

// Command seed populates a SQLite database with Pokémon-themed trainers,
// groups, and messages, all written through the real UserService/
// GroupService — not raw DB inserts — so the seeded data exercises the same
// validation and side effects as the app itself.
//
// Point DATABASE_PATH at a scratch file before running so this doesn't
// touch real data, e.g.:
//
//	DATABASE_PATH=./internal/db/seed.pokemon.sqlite3 go run ./cmd/seed
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/alextilot/golang-htmx-chatapp/config"
	"github.com/alextilot/golang-htmx-chatapp/internal/db"
	"github.com/alextilot/golang-htmx-chatapp/internal/model"
	"github.com/alextilot/golang-htmx-chatapp/internal/repository"
	"github.com/alextilot/golang-htmx-chatapp/internal/service"
)

// This is a fake password.
const seedPassword = "Pikachu1!"

type trainerSeed struct {
	username string
	email    string
}

type messageSeed struct {
	sender  string
	content string
}

type groupSeed struct {
	name     string
	members  []string // usernames; first is the creator
	messages []messageSeed
}

func main() {
	ctx := context.Background()
	cfg := config.Load()

	database := db.NewDatabase(db.Config{DSN: cfg.DatabasePath})
	defer database.Close()
	database.AutoMigrate()

	repos := repository.NewRepositories(repository.Deps{DB: database.Conn})
	services := service.NewServices(service.Deps{Repos: repos})

	log.Printf("🌱 seeding Pokémon data into %s", cfg.DatabasePath)

	trainers := []trainerSeed{
		{"ash_ketchum", "ash@pallet.town"},
		{"misty", "misty@cerulean.gym"},
		{"brock", "brock@pewter.gym"},
		{"gary_oak", "gary@oak.lab"},
		{"may", "may@hoenn.region"},
		{"dawn", "dawn@sinnoh.region"},
		{"serena", "serena@kalos.region"},
	}

	users := make(map[string]*model.User, len(trainers))
	for _, t := range trainers {
		u, err := services.UserService.Signup(ctx, service.SignupInput{
			Username:       t.username,
			Email:          t.email,
			Password:       seedPassword,
			RepeatPassword: seedPassword,
		})
		if err != nil {
			log.Fatalf("signup %s: %v", t.username, err)
		}
		users[t.username] = u
	}

	groups := []groupSeed{
		{
			name:    "Pallet Town Crew",
			members: []string{"ash_ketchum", "misty", "brock", "gary_oak"},
			messages: []messageSeed{
				{"ash_ketchum", "Guys I finally caught a Charizard!! 🔥"},
				{"misty", "You mean it finally stopped disobeying you?"},
				{"brock", "Ha! Reminds me of my Onix days."},
				{"gary_oak", "Cute. I have ten of them."},
				{"ash_ketchum", "It's not about how many, Gary, it's about the bond!"},
				{"brock", "Someone's gotta cook before this turns into a rivalry battle."},
			},
		},
		{
			name:    "Elite Four Prep",
			members: []string{"ash_ketchum", "may", "dawn"},
			messages: []messageSeed{
				{"dawn", "Ash, did you finish training for the Sinnoh League?"},
				{"may", "He's probably still stuck on type match-ups lol"},
				{"ash_ketchum", "Hey! I know Water beats Fire now!"},
				{"dawn", "...and Fire beats Grass, and Grass beats Water, right?"},
				{"ash_ketchum", "...right?"},
				{"may", "We have so much work to do."},
			},
		},
		{
			name:    "Rival Rundown",
			members: []string{"ash_ketchum", "gary_oak", "serena"},
			messages: []messageSeed{
				{"gary_oak", "Heard you lost to a Magikarp trainer in Kalos."},
				{"serena", "It evolved mid-battle, that's not really fair to leave out."},
				{"ash_ketchum", "IT BECAME A GYARADOS, GARY."},
				{"gary_oak", "Sure it did."},
				{"serena", "I filmed it, I can send proof."},
			},
		},
	}

	for _, g := range groups {
		creator := users[g.members[0]]

		grp, err := services.GroupService.Create(ctx, creator.ID, service.CreateGroupInput{
			Name: g.name,
			Type: model.GroupTypeGroup,
		})
		if err != nil {
			log.Fatalf("create group %s: %v", g.name, err)
		}

		for _, member := range g.members[1:] {
			u := users[member]
			if err := services.GroupService.AddMember(ctx, grp.ID, creator.ID, service.AddMemberInput{UserID: u.ID}); err != nil {
				log.Fatalf("add member %s to %s: %v", member, g.name, err)
			}
		}

		for _, m := range g.messages {
			sender := users[m.sender]
			if _, err := services.GroupService.SendMessage(ctx, grp.ID, sender.ID, service.SendMessageInput{Content: m.content}); err != nil {
				log.Fatalf("send message in %s: %v", g.name, err)
			}
		}

		log.Printf("✅ seeded group %q with %d members, %d messages", g.name, len(g.members), len(g.messages))
	}

	fmt.Println()
	fmt.Println("Seed complete. Log in as any trainer with password:", seedPassword)
	for _, t := range trainers {
		fmt.Printf("  %s\n", t.username)
	}
}

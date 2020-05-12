package main

import (
	"fmt"
	"log"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
)

func main() {
	db, err := bolt.Open("guardian.db", 0600, nil)
	if err != nil {
		log.Fatalf("failed to open db handler %v", err)
	}
	defer db.Close()
	dgo, err := discordgo.New("Bot " + "authentication token")
	if err != nil {
		log.Fatalf("failed to connect to discord api endpoint")
	}
	dgo.AddHandler(func(s *discordgo.Session, m *discordgo.MessageCreate) {
		fmt.Println(m.Author.Username)
	})

}

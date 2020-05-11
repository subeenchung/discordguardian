package main

import (
	"log"

	"github.com/boltdb/bolt"
	"github.com/bwmarrin/discordgo"
)

func main() {
	_, err := bolt.Open("guardian.db", 0600, nil)
	if err != nil {
		log.Fatalf("failed to open db handler %v", err)
	}

	_, err := discordgo.New("Bot " + "authentication token")
	if err != nil {
		log.Fatalf("failed to connect to discord api endpoint")
	}

}

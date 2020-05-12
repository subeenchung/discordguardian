package main

import (
	"fmt"
	"io/ioutil"
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

	token, err := ioutil.ReadFile("token.pem")
	if err != nil {
		log.Fatalf("failed to retrieve token file")
	}

	dgo, err := discordgo.New("Bot " + string(token))
	if err != nil {
		log.Fatalf("failed to connect to discord api endpoint")
	}
	dgo.AddHandler(messageHandler)

}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	fmt.Println(m.Author.Username)
}

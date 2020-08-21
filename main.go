package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type config struct {
	DBuser    string `json:"dbuser"`
	DBpass    string `json:"dbpass"`
	DBname    string `json:"dbname"`
	DBurl     string `json:"dburl"`
	DBoptions string `json:"dboptions"`
}

var today string

func main() {

	// Instantiate new config struct
	cfg := &config{}

	// Read config file
	fcfg, err := ioutil.ReadFile("config.json")
	errCheck(err)
	//Inject data into config struct
	err = json.Unmarshal(fcfg, &cfg)
	errCheck(err)

	// Connect to MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	log.Println("Connecting to DB backend...")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(
		fmt.Sprintf(
			"mongodb+srv://%s:%s@%s/%s?%s", cfg.DBuser, cfg.DBpass, cfg.DBurl, cfg.DBname, cfg.DBoptions),
	))

	errCheck(err)

	err = client.Ping(context.Background(), nil)
	errCheck(err)

	log.Println(" - DB backend connection success - ")

	// Instantiate today's date var with current date
	today = getDate()

	// Load bot token from file "token.pem"
	token, err := ioutil.ReadFile("token.pem")
	if err != nil {
		log.Fatalf("failed to retrieve token file\n")
	}
	log.Println("Connecting to Discord API and creating new session object...")
	// Discord Go Setup
	dgo, err := discordgo.New("Bot " + string(token))
	if err != nil {
		log.Fatalf("failed to create discordgo session object\n")
	}

	// Open connection with Discord API
	err = dgo.Open()
	if err != nil {
		log.Fatalf("failed to connect to discord api: %v\n", err)
	}
	log.Println("Connected to Discord API")
	fmt.Println("List of guilds that are using this bot:")
	for _, g := range getAllGuilds(dgo) {
		newgd, err := dgo.Guild(g.ID)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("ID: %s  -  Name: %s\n", g.ID, newgd.Name)
	}

	AddHandlers(dgo)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for {
			if strings.Compare(today, getDate()) != 0 {
				today = getDate()
			}
		}
	}()

	s := <-c
	fmt.Printf("Received signal: %v, closing Discord Bot\n", s)
	_ = dgo.Close()

	os.Exit(0)

}

//GetDate function returns the today's date in string format of YYYY-MM-DD
func getDate() string {
	t := time.Now()
	return t.Format("2006-01-02")
}

//AddHandlers function Adds all handlers
func AddHandlers(s *discordgo.Session) {
	s.AddHandler(messageHandler)
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	/*
		db.Update(func(tx *bolt.Tx) error {
			fmt.Printf("%s_%s\n", m.GuildID, today)
			b := tx.Bucket([]byte(getGuildDate(m.GuildID)))
			if b == nil {
				fmt.Println("bucket does not exist")
				newb, _ := tx.CreateBucket([]byte(getGuildDate(m.GuildID)))
				b = newb
			}
			r := b.Get([]byte(m.Author.Username))
			if r == nil {
				b.Put([]byte(m.Author.Username), []byte("1"))
			}
			return nil
		})
	*/
	fmt.Println("logged: " + m.Author.Username)
}

func getAllGuilds(s *discordgo.Session) []*discordgo.Guild {
	return s.State.Guilds
}

func getGuildDate(g string) string {
	return fmt.Sprintf("%s_%s", g, today)
}

func errCheck(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	// "github.com/goccy/go-yaml"
)

// Variables used for command line parameters
var (
	Token    string
	allVerbs []Verb
	Prefix   string
)

// type strings []string
type gifs []*Gif
type verbs []Verb

// Verb is a collection of gifs
type Verb struct {
	Images gifs
	Name   string
}

// Gif is a struct containing a url and tags
type Gif struct {
	URL  string
	Tags []string
}

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Prefix, "p", "+", "Prefix")
	flag.Parse()
}

func main() {
	AutoConfig()

	if Prefix == "+" {
		Prefix = viper.GetString("prefix")
	}

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// In this example, we only care about receiving message events.
	dg.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsGuildMessages)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	// Wait here until CTRL-C or other term signal is received.
	log.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	g := &Gif{
		URL:  "https://media1.tenor.com/images/4d89d7f963b41a416ec8a55230dab31b/tenor.gif?itemid=5166500",
		Tags: []string{"test", "t"},
	}
	// g2 := &Gif{
	// 	URL:  "https://uberi.fi",
	// 	Tags: []string{"test", "a"}}

	collection := gifs{g}
	allVerbs = []Verb{{collection, "pat"}} // Store(allVerbs)
	v := []Verb{}

	Load(&v)
	allVerbs = v
	Store(v)
	//fmt.Println(v)
	<-sc

	log.Println("Closing Bot.")

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, mCreate *discordgo.MessageCreate) {
	m := mCreate.Message

	log.Println("Inside message")
	out := discordgo.MessageSend{}

	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID {
		return
	}
	splitm := strings.Split(m.Content, " ")

	if m.Content == viper.GetString("cmdprefix")+"verbs" {
		out = ListVerbs(&allVerbs)
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == "pong" {
		s.ChannelMessageSendComplex(m.ChannelID, &out)
	}
	if (strings.HasPrefix(splitm[0], Prefix) && out.Content == discordgo.MessageSend{}.Content) {
		out = VerbCommand(m, s, &allVerbs)
	}

	s.ChannelMessageSendComplex(m.ChannelID, &out)

}

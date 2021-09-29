package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/Secret-Society-Blanket/goshleep/lib"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

// Variables used for command line parameters
var (
	Token    string
	allVerbs []goshleep.Verb
	Prefix   string
	responses []goshleep.Response
)

type gifs []goshleep.Gif
type verbs []goshleep.Verb

func init() {

	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&Prefix, "p", "+", "Prefix")
	flag.Parse()
}

func main() {
	goshleep.AutoConfig()

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

	allVerbs = []goshleep.Verb{{Images: []goshleep.Gif{}, Name: "pat"}} // Store(allVerbs)
	v := []goshleep.Verb{}

	goshleep.Load(&v)
	allVerbs = v
	goshleep.Store(v)
	<-sc

	log.Println("Closing bot.")

	// Cleanly close down the Discord session.
	dg.Close()
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the authenticated bot has access to.
func messageCreate(s *discordgo.Session, mCreate *discordgo.MessageCreate) {
	if (len(responses) > 0) {
		for _, r := range responses {
			if (goshleep.Contains(r.UserList, mCreate.Author.ID)){
				if (r.Check(mCreate.Content)) {
					r.Run(mCreate.Content, s, &allVerbs)
				}
			}
		}
	}
	if strings.HasPrefix(mCreate.Message.Content, "+") {
		m := goshleep.ConstructRequest(*mCreate.Message)

		log.Println("Inside message")
		out := discordgo.MessageSend{}

		// Ignore all messages created by the bot itself
		// This isn't required in this specific example but it's a good practice.
		if mCreate.Author.ID == s.State.User.ID {
			return
		}
		out = goshleep.ParseRequest(&m, s, &allVerbs)
		if (m.Resp != nil) {
			responses = append(responses, *m.Resp)
		}
		s.ChannelMessageSendComplex(mCreate.ChannelID, &out)
	}
}

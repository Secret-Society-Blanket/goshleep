package goshleep

import (
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

func Eightball(details *Request, s *discordgo.Session, _ *[]Verb) discordgo.MessageSend {

	var out string
	var question string
	if (len(details.SplitContent) <= 1){
		question = "That's not a question, bozo."
	} else {
		question = strings.Join(details.SplitContent[1:], " ")
	}

	out = strings.ReplaceAll(eightballTemplate, "QUESTION", question)

	answers := viper.GetStringSlice("eightballMessages")
	numAnswers := len(answers)

	out = strings.ReplaceAll(out, "ANSWER", answers[rand.Intn(numAnswers - 1)])

	 m := discordgo.MessageSend{
		 Content: out,
		 Reference: details.dMessage.Reference(),
	 }

	return m

}

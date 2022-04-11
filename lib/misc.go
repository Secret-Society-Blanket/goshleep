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
	if len(details.SplitContent) <= 1 {
		question = "There's no question!"
	} else {
		question = strings.Join(details.SplitContent[1:], " ")
	}

	out = strings.ReplaceAll(eightballTemplate, "QUESTION", question)

	answers := viper.GetStringSlice("eightballMessages")
	numAnswers := len(answers)

	out = strings.ReplaceAll(out, "ANSWER", answers[rand.Intn(numAnswers-1)])

	m := discordgo.MessageSend{
		Content:   out,
		Reference: details.dMessage.Reference(),
	}

	return m

}

func ChooseCommand(details *Request, s *discordgo.Session, _ *[]Verb) discordgo.MessageSend {
	out := baseMessage(details)
	choices := strings.Split(details.Content[len(details.SplitContent[0]):], "|")
	chosenOption := choices[rand.Intn(len(choices))]
	name := GetName(details.dMessage.Member)
	out.Content = strings.Replace(chooseTemplate, "NAME", name, 1)
	out.Content = strings.Replace(out.Content, "CHOICE",
		strings.TrimSpace(chosenOption), 1)
	return out
}



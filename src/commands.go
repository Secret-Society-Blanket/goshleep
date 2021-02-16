package main

import (
	"github.com/bwmarrin/discordgo"
	"strings"
)


type CommandFunction func (details Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend

type Request struct {

	dMessage discordgo.Message

	Content string

	SplitContent []string

	Type Command
}

type Command struct {

	Name string

	Description string

	HotStrings []string

	Function CommandFunction

	Admin bool

	Priority int // 0 = lowest

}

var (
	AllCommands []Command = []Command{



		// This should always be last, since it uses + as it's only prefix
		Command {
			Name: "Verb",
			Description: "Posts a gif based on the arguments the user gives",
			HotStrings: []string{"+"},
			Function: VerbCommand,
			Admin: false,
			Priority: 0,
		} }
)

func ConstructRequest (m discordgo.Message) Request {

	split := strings.Split(m.Content, " ")

	var cmd Command

	for i:= 0; i < len(AllCommands); i++ {
		if (IfInString(AllCommands[i].HotStrings, split[0])) {
			cmd = AllCommands[i]
		}
	}


	out := Request{

		Content: m.Content,
		SplitContent: split,
		dMessage: m,
		Type: cmd,

	}

	return out


}

func ParseRequest (r Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend  {
	return r.Type.Function(r, s, allVerbs)
}

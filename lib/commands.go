package goshleep

import (
	"strings"
	"log"
	"github.com/bwmarrin/discordgo"
)

type CommandFunction func(details *Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend

type Request struct {
	dMessage discordgo.Message

	Content string

	SplitContent []string

	Type Command

	Resp *Response `default: nil`
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

		// This should always be first, since it uses + as it's only prefix
		{
			Name:        "Verb",
			Description: "Posts a GIF based on the arguments the user gives",
			HotStrings:  []string{""},
			Function:    VerbCommand,
			Admin:       false,
			Priority:    0,
		},
		{
			Name:        "List Verbs",
			Description: "Lists all verbs",
			HotStrings:  []string{"verbs"},
			Function:    ListVerbs,
			Admin:       false,
			Priority:    0,
		},
		{
			Name:        "Eightball",
			Description: "Ask Shleepbot your most pressing questions...",
			HotStrings:  []string{"eightball"},
			Function:    ListVerbs,
			Admin:       false,
			Priority:    0,
		},
	}
)

func ConstructRequest(m discordgo.Message) Request {

	split := strings.Split(m.Content, " ")

	var cmd *Command

	for i := 0; i < len(AllCommands); i++ {
		if IfMatchHotStrings(AllCommands[i].HotStrings, split[0]) {
			cmd = &AllCommands[i]
		}
	}

	if cmd == nil {
		log.Println("Found no command")
	}

	out := Request{
		Content:      m.Content,
		SplitContent: split,
		dMessage:     m,
		Type:         *cmd,
		Resp:         nil,
	}

	return out
}

func ParseRequest(r *Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	log.Println(r.Type.Name)
	log.Println(r.Type.HotStrings)
	return r.Type.Function(r, s, allVerbs)
}

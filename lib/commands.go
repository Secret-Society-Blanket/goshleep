package goshleep

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

type CommandFunction func(details *Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend

type Request struct {
	dMessage     discordgo.Message
	Content      string
	SplitContent []string
	Type         Command
	Resp         *Response
}

type Command struct {
	Name        string
	Description string
	Template    string
	HotStrings  []string
	Function    CommandFunction
	Admin       bool
	Priority    int // 0 = lowest
}

func baseMessage(details *Request) discordgo.MessageSend {
	m := discordgo.MessageSend{}
	m.Reference = details.dMessage.Reference()
	m.AllowedMentions = nil
	return m
}

var (
	AllCommands []Command = []Command{

		// This should always be first, since it uses + as it's only prefix
		{
			Name:        "Verb",
			Description: "Posts a GIF based on the arguments the user gives",
			Template:    "+<verbname> @<user> [-t <tags>]",
			HotStrings:  []string{""},
			Function:    VerbCommand,
			Admin:       false,
			Priority:    0,
		},
		{
			Name:        "List Verbs",
			Description: "Lists all verbs",
			Template:    "+verbs",
			HotStrings:  []string{"verbs"},
			Function:    ListVerbs,
			Admin:       false,
			Priority:    0,
		},
		{
			Name:        "Eightball",
			Description: "Ask Shleepbot your most pressing questions...",
			Template:    "+eightball Should I eat some ice cream?",
			HotStrings:  []string{"eightball"},
			Function:    Eightball,
			Admin:       false,
			Priority:    0,
		},
		{
			Name:        "F",
			Description: "Pay your respects for something/someone",
			Template:    "+f Josh",
			HotStrings:  []string{"f"},
			Function:    FCommand,
			Admin:       false,
			Priority:    0,
		},
		{
			Name:        "Choose",
			Description: "Choose between a number of options!",
			Template:    "+choose <option 1> | <option 2> | <option 3>",
			HotStrings:  []string{"choose"},
			Function:    ChooseCommand,
			Admin:       false,
			Priority:    0,
		},
		{
			Name:        "Add Gif",
			Description: "Adds a gif, and creates a verb if needed.",
			Template:    "+add <verb> <url> [-t <tags>]",
			HotStrings:  []string{"add"},
			Function:    AddGifCommand,
			Admin:       true,
			Priority:    0,
		},
		{
			Name:        "Remove Gif",
			Description: "Removes a gif from a given verb.",
			Template:    "+remove <verb> <url>",
			HotStrings:  []string{"remove"},
			Function:    RemoveGifCommand,
			Admin:       true,
			Priority:    0,
		},
		{
			Name:        "AddAdmin",
			Description: "Adds a given discord ID or mention as an admin",
			Template:    "+addAdmin @<user>",
			HotStrings:  []string{"addAdmin", "aadd", "adminAdd"},
			Function:    AddAdminCommand,
			Admin:       true,
			Priority:    0,
		},
		{
			Name:        "Add Synonym",
			Description: "Adds a verb as a synonym for another",
			Template:    "+synonym <new verb> <base verb>",
			HotStrings:  []string{"synonym", "sadd", "synonymadd"},
			Function:    AddSynonymCommand,
			Admin:       true,
			Priority:    0,
		},
	}
)

func ConstructRequest(m discordgo.Message) Request {

	split := strings.Split(m.Content, " ")

	// This replaces AllCommands with AllCommands + Help. If we don't do this, it
	// causes an initialization loop
	AllCommands := append(AllCommands, (Command{
		Name:        "Help",
		Description: "Display this message!",
		Template:    "+help",
		HotStrings:  []string{"help"},
		Function:    HelpCommand,
		Admin:       false,
		Priority:    0,
	}))

	var cmd *Command

	for i := 0; i < len(AllCommands); i++ {
		if IfMatchHotStrings(AllCommands[i].HotStrings, split[0]) {
			cmd = &AllCommands[i]
		}
	}

	if cmd == nil {
		log.Println("Found no extra command, defaulting to verb.")
		/* This should always be the Verb command. The reason this is necessary is
		 in order to fix an error involving substrings of commands, the verb command would
		basically never be found. */
		cmd = &AllCommands[0]
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

// This is the function that actually runs a request, or command.
func ParseRequest(r *Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	log.Println(r.Type.Name)
	log.Println(r.Type.HotStrings)
	ReadViper()
	// If this command needs an admin, verify user is an admin. 
	if r.Type.Admin {
		if Contains(viper.GetStringSlice("admins"), r.dMessage.Author.ID) {
			return r.Type.Function(r, s, allVerbs)
		} else {
			return discordgo.MessageSend{
				Content:   "You tried to run an admin command, but you aren't an admin!",
				Reference: r.dMessage.Reference(),
			}
		}

	} else {
		// This is what actually runs the function.
		return r.Type.Function(r, s, allVerbs)
	}
}

func HelpCommand(details *Request, s *discordgo.Session, _ *[]Verb) discordgo.MessageSend {

	m := baseMessage(details)

	msg := ""
	IsAdmin := Contains(viper.GetStringSlice("admins"), details.dMessage.Author.ID)
	for _, cmd := range AllCommands {

		if cmd.Admin == true && !IsAdmin {
			continue
		}
		msg = msg + cmd.Name + ":\n\t"
		msg = msg + "Description: " + cmd.Description + "\n\t"
		msg = msg + "Template: ``" + cmd.Template + "``\n\n"
	}
	m.Content = msg

	return m

}

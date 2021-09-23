package goshleep

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

type ResponseFunction func(content string, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend
type CheckFunction func(content string) bool

// Response is a struct that lets goshleep ask a user questions and get
// respsonses.
type Response struct {
	// Users IDS that can respond
	UserList []string
	// Function to run the response through
	Run ResponseFunction
	// Is this response still alive?
	Alive bool
	// Verify if the given message satisfies requirements
	Check CheckFunction
	// When should this response die?
	timeout time.Time
	// When was this response created?
	Created time.Time
	Kill KillFunction
}

// KillFunction is a function type used to kill Responses
type KillFunction func(r Response)

// DefaultKillFunction is a minimalist kill function that simply flips alive to
// off.
func DefaultKillFunction (r Response) {
	r.Alive = false;
}

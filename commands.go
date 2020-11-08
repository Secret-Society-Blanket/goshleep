package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
)

// VerbCommand takes in the follow expected string template:
// +<verb> [recipient] [-t tags]
// and returns a discord message, including a gif reaction
//of the inputted verb, assuming one exists.
func VerbCommand(cmd []string, allVerbs *[]Verb) discordgo.MessageSend {
	m := discordgo.MessageSend{}

	v := strings.ToLower(cmd[0][len(Prefix):])

	log.Println("Found Verb: " + v)

	return m
}

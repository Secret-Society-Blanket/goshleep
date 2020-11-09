package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"log"
	"sort"
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

func getVerb(toFind string, allVerbs *[]Verb) (*Verb, bool) {
	var out *Verb
	fuzz := false
	for _, v := range *allVerbs {
		if fuzzy.MatchFold(toFind, v.Name) {
			if toFind != v.Name {
				fuzz = true
			}
			out = &v
		}
	}

	return out, fuzz

}

package main

import (
	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"log"
	"strings"
)

// VerbCommand takes in the follow expected string template:
// +<verb> [recipient] [-t tags]
// and returns a discord message, including a gif reaction
// of the inputted verb, assuming one exists.
func VerbCommand(ogMessage *discordgo.MessageCreate, allVerbs *[]Verb) discordgo.MessageSend {
	m := discordgo.MessageSend{}

	cmd := strings.Split(ogMessage.Content, " ")

	v, _ := getVerb(cmd[0][len(Prefix):], allVerbs)

	if v == nil {
		log.Println("Failed to find verb!")
		m.Content = "I couldn't find it..."
	} else {
		log.Println("Found Verb: " + v.Name)
		var i int
		var w string
		for i, w = range cmd {
			log.Println(w)
			if w == "-t" {
				break
			}
		}
		log.Println(i)
		if i > 0 {
			recipientArray := cmd[1 : i+1]
			log.Println(recipientArray)
			recipient := strings.Join(recipientArray, " ")
			title := "**RECIPIENT**, you got a **VERB** from **SENDER**"
			if len(ogMessage.Mentions) > 0 {
				// title = strings.ReplaceAll(title, "RECIPIENT", getName(*ogMessage.Mentions[0].))
			}
			title = strings.ReplaceAll(title, "RECIPIENT", recipient)
			title = strings.ReplaceAll(title, "VERB", v.Name)
			title = strings.ReplaceAll(title, "SENDER", getName(*ogMessage.Member))
			m.Embed = &discordgo.MessageEmbed{
				Title: title,
			}
		} else {
			title := "**SENDER sent a VERB**"
			title = strings.ReplaceAll(title, "VERB", v.Name)
			title = strings.ReplaceAll(title, "SENDER", getName(*ogMessage.Member))
			m.Embed = &discordgo.MessageEmbed{
				Title: title,
			}
		}
	}

	return m
}

func ListVerbs(allVerbs *[]Verb) discordgo.MessageSend {
	m := discordgo.MessageSend{}

	var verbNames []string
	log.Println(*allVerbs)
	for _, v := range *allVerbs {
		verbNames = append(verbNames, v.Name)
		// Temporary, change this later!
		m.Content = m.Content + v.Name + "\n"
		log.Println(v.Name)
	}

	return m
}

func getVerb(toFind string, allVerbs *[]Verb) (*Verb, bool) {
	var out *Verb
	out = nil
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

func getName(u discordgo.Member) string {
	if u.Nick != "" {
		return u.Nick
	}
	return u.User.Username
}
func getName(u discordgo.User, g discordgo.Guild) string {
	if u.Nick != "" {
		return u.Nick
	}
	return u.User.Username
}

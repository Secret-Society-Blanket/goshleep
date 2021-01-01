package main

import (
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// VerbCommand takes in the follow expected string template:
// +<verb> [recipient] [-t tags]
// and returns a discord message, including a gif reaction
// of the inputted verb, assuming one exists.
func VerbCommand(ogMessage *discordgo.Message, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	// Create the message to send out
	m := discordgo.MessageSend{}

	cmd := strings.Split(ogMessage.Content, " ")

	// This is finding the "verb" command
	v, _ := getVerb(cmd[0][len(Prefix):], allVerbs)

	// If there is no verb
	if v == nil {
		log.Println("Failed to find verb!")
		m.Content = "I couldn't find it..."
	} else {
		log.Println("Found Verb: " + v.Name)
		// These need to be declared early so they can be used outside the loop
		var i int
		var w string
		for i, w = range cmd {
			log.Println(w)
			// Only go through till u see -t
			if w == "-t" {
				break
			}
		}
		log.Println(i)
		// If there is a recipient
		if i > 0 {
			// Create an array from everything after the verb to the -t (assuming it exists)
			recipientArray := cmd[1 : i+1]
			log.Println("Found names", recipientArray)
			recipient := strings.Join(recipientArray, " ")

			// This is the format of the message
			title := "**RECIPIENT**, you got a **VERB** from **SENDER**"

			// If there are mentions, use them as the recipients
			if len(ogMessage.Mentions) > 0 {
				recipient = GetMentionNames(ogMessage, s)
				log.Println(recipient)
			}
			title = strings.ReplaceAll(title, "RECIPIENT", recipient)
			title = strings.ReplaceAll(title, "VERB", v.Name)
			title = strings.ReplaceAll(title, "SENDER", GetName(ogMessage.Member))
			m.Embed = &discordgo.MessageEmbed{
				Title: title,
			}
			// If there isn't a recipient
		} else {
			title := "**SENDER sent a VERB**"
			title = strings.ReplaceAll(title, "VERB", v.Name)
			title = strings.ReplaceAll(title, "SENDER", GetName(ogMessage.Member))
			m.Embed = &discordgo.MessageEmbed{
				Title: title,
			}
		}
		// Get the actual image
		log.Println("Images", len(v.Images))
		var num int
		if len(v.Images) > 1 {
			log.Println("I found a number of images")
			num = rand.Intn(len(v.Images) - 1)
		} else {
			log.Println("I only found 1 image")
			num = 0
		}
		m.Embed.Image = &discordgo.MessageEmbedImage{
			URL: v.Images[num].URL,
		}

	}

	return m
}

// ListVerbs lists all verbs in allVerbs
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
	var out Verb
	out = Verb{
		Images: gifs{
			&Gif{
				URL:  "https://animemotivation.com/wp-content/uploads/2020/06/cute-anime-cat-girl-confused-e1592069452432.jpg",
				Tags: []string{},
			},
		},
		Name: "unknown",
	}
	fuzz := false
	last := 10
	for _, v := range *allVerbs {
		// Get the levensthien distance if it matches, otherwise return -1
		i := fuzzy.LevenshteinDistance(strings.ToLower(toFind),
			strings.ToLower(v.Name))

		// If it matches, and the distance is less than 4
		if i != -1 && i < 4 {
			log.Println(v.Name, "matches", toFind)

			// Is the distance less than the last?
			// If not, ignore the result
			if i < last {
				if i != 0 {
					fuzz = true
				} else {
					fuzz = false
				}
				out = v
				last = i
			}

		}
	}

	return &out, fuzz
}
package goshleep

import (
	"log"
	"math/rand"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/spf13/viper"
)

// Verb is a collection of gifs
type Verb struct {
	Images []*Gif
	Name   string
}

// Gif is a struct containing a url and tags
type Gif struct {
	URL  string
	Tags []string
}

// VerbCommand takes in the follow expected string template:
// +<verb> [recipient] [-t tags]
// and returns a discord message, including a gif reaction
// of the inputted verb, assuming one exists.
func VerbCommand(myRequest Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	// Create the message to send out
	m := discordgo.MessageSend{}

	cmd := myRequest.SplitContent

	// This is finding the "verb" command
	v, _ := getVerb(cmd[0][len(viper.GetString("prefix")):], allVerbs)

	// If there is no verb
	if v == nil {
		log.Println("Failed to find verb!")
		m.Content = "I couldn't find it..."
	} else {
		log.Println("Found Verb: " + v.Name)
		// These need to be declared early so they can be used outside the loop
		var i int
		var w string
		var tags []string
		for i, w = range cmd {
			// Only go through till u see -t
			if w == "-t" {
				break
			}
		}
		tags = cmd[i+1:]
		// If there is a recipient
		if i > 0 {
			// Create an array from everything after the verb to the -t (assuming it exists)
			recipientArray := cmd[1:i]
			log.Println("Found names", recipientArray)
			recipient := strings.Join(recipientArray, " ")
			log.Println("Found names", recipient)

			// This is the format of the message
			title := "**RECIPIENT**, you got a **VERB** from **SENDER**"

			// If there are mentions, use them as the recipients
			if len(myRequest.dMessage.Mentions) > 0 {
				recipient = GetMentionNames(&myRequest.dMessage, s)
			}
			title = strings.ReplaceAll(title, "RECIPIENT", recipient)
			title = strings.ReplaceAll(title, "VERB", v.Name)
			title = strings.ReplaceAll(title, "SENDER", GetName(myRequest.dMessage.Member))
			m.Embed = &discordgo.MessageEmbed{
				Title: title,
			}
			// If there isn't a recipient
		} else {
			title := "**SENDER sent a VERB**"
			title = strings.ReplaceAll(title, "VERB", v.Name)
			title = strings.ReplaceAll(title, "SENDER", GetName(myRequest.dMessage.Member))
			m.Embed = &discordgo.MessageEmbed{
				Title: title,
			}
		}
		// Get the actual image
		img, _ := getImage(v, tags)
		m.Embed.Image = &discordgo.MessageEmbedImage{
			URL: img.URL,
		}
	}

	return m
}

// ListVerbs lists all verbs in allVerbs
func ListVerbs(_ Request, _ *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	m := discordgo.MessageSend{}

	m.Content = "```"
	var verbNames []string
	for _, v := range *allVerbs {
		verbNames = append(verbNames, v.Name)
		// Temporary, change this later!
		m.Content = m.Content + v.Name + "\n"
	}
	m.Content = m.Content + "\n```"
	return m
}

func getVerb(toFind string, allVerbs *[]Verb) (*Verb, bool) {
	var out Verb
	out = Verb{
		Images: []*Gif{
			{
				// Default picture, for failures
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

func getImage(v *Verb, tags []string) (*Gif, bool) {

	var num int
	tag := false
	var allGifs []*Gif

	if len(tags) > 0 {
		// Get all valid images
		for _, g := range v.Images {
			for _, t := range tags {
				if Contains(g.Tags, t) {
					allGifs = append(allGifs, g)
					tag = true
				}
			}
		}

	}
	if len(allGifs) == 0 {
		allGifs = v.Images
	}
	if len(allGifs) > 1 {
		log.Println("I found a number of images")
		num = rand.Intn(len(allGifs) - 1)
	} else {
		log.Println("I only found 1 image")
		num = 0
	}

	return allGifs[num], tag
}

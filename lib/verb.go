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
	Images []Gif
	Names  []string
}

// Gif is a struct containing a url and tags
type Gif struct {
	URL  string
	Tags []string
}

type pair struct {
	num  int
	name string
}

// VerbCommand takes in the follow expected string template:
// +<verb> [recipient] [-t tags]
// and returns a discord message, including a gif reaction
// of the inputted verb, assuming one exists.
func VerbCommand(myRequest *Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	// Create the message to send out
	m := discordgo.MessageSend{}

	cmd := myRequest.SplitContent

	verbName := cmd[0][len(viper.GetString("prefix")):]
	// This is finding the "verb" command
	v, f := getVerb(verbName, allVerbs)
	if f {
		verbName = v.Names[0]
	}

	// If there is no verb
	if v == nil {
		log.Println("Failed to find verb!")
	} else {
		log.Println("Found Verb: " + v.Names[0])
		// These need to be declared early so they can be used outside the loop
		var i int
		var w string
		var tags []string
		for i, w = range cmd {
			// Only go through till it sees -t
			if w == "-t" {
				break
			}
		}
		tags = cmd[i+1:]
		// If there is a recipient
		if i > 0 {
			// Create an array from everything after the verb to the -t (assuming it exists)
			var recipientArray []string
			if len(cmd) != i+1 {
				recipientArray = cmd[1:i]
			} else {
				recipientArray = cmd[1 : i+1]
			}
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
			title = strings.ReplaceAll(title, "VERB", verbName)
			title = strings.ReplaceAll(title, "SENDER", GetName(myRequest.dMessage.Author, myRequest.dMessage.GuildID, s))
			m.Embed = &discordgo.MessageEmbed{
				Title: title,
			}
			// If there isn't a recipient
		} else {
			title := "**SENDER sent a VERB**"
			title = strings.ReplaceAll(title, "VERB", verbName)
			title = strings.ReplaceAll(title, "SENDER", GetName(myRequest.dMessage.Author, myRequest.dMessage.GuildID, s))
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
func ListVerbs(_ *Request, _ *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	m := discordgo.MessageSend{}

	m.Content = "```"
	var verbNames []string
	all := allNames(allVerbs)
	lastNum := -1
	num := -1
	for _, v := range all {
		verbNames = append(verbNames, v.name)
		num = v.num
		if num != lastNum {
			m.Content = m.Content + "\n"
		} else {
			m.Content = m.Content + ", "
		}
		// Temporary, change this later!
		m.Content = m.Content + v.name
		lastNum = num
	}
	m.Content = m.Content + "\n```"
	return m
}

func getVerb(toFind string, allVerbs *[]Verb) (*Verb, bool) {
	var out *Verb
	out = &Verb{
		Images: []Gif{
			{
				// Default picture, for failures
				URL:  "https://animemotivation.com/wp-content/uploads/2020/06/cute-anime-cat-girl-confused-e1592069452432.jpg",
				Tags: []string{},
			},
		},
		Names: []string{"unknown"},
	}
	fuzz := false
	last := 10

	all := allNames(allVerbs)

	// If no verb (if it's an empty prefix), return nothing
	if toFind == "" {
		return nil, false
	}

	for _, n := range all {
		// Get the levensthien distance if it matches, otherwise return -1
		i := fuzzy.LevenshteinDistance(strings.ToLower(toFind),
			strings.ToLower(n.name))

		// If it matches, and the distance is less than 4
		if i != -1 && i < 4 {
			log.Println(n.name, "matches", toFind)
			// Is the distance less than the last?
			// If not, ignore the result
			if i < last {
				if i != 0 {
					fuzz = true
				} else {
					fuzz = false
				}
				out = &(*allVerbs)[n.num]
				last = i
			}
		}
	}

	return out, fuzz
}

func allNames(allVerbs *[]Verb) []pair {
	out := []pair{}

	for i, v := range *allVerbs {
		for _, n := range v.Names {
			p := pair{
				num:  i,
				name: n,
			}
			out = append(out, p)
		}
	}
	return out
}

func getImage(v *Verb, tags []string) (*Gif, bool) {
	var num int
	tag := false
	var allGifs []*Gif

	if len(tags) > 0 {
		// Get all valid images
		for i, g := range v.Images {
			for _, t := range tags {
				if Contains(g.Tags, t) {
					log.Println("Found tag: " + t)
					allGifs = append(allGifs, &v.Images[i])
					tag = true
				}
			}
		}

	}
	// If there are no tags, or if nothing is found
	if len(allGifs) == 0 {
		log.Println("I couldn't find any matching gifs, so I'm just using all of them")
		allGifs = makeReference(v.Images)
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

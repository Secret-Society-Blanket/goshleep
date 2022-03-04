package goshleep

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

// AddGifCommand is an admin command that allows admins to add a gif, and creates a verb if it does not exist.
func AddGifCommand(details *Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	log.Println("Backing up before modification...")
	Store(*allVerbs)
	m := baseMessage(details)
	vName := details.SplitContent[1]
	urlString := details.SplitContent[2]
	if len(details.SplitContent) == 3 {
		details.SplitContent = append(details.SplitContent, "")
	} else {
		log.Println(details.SplitContent[4:])
	}
	g := Gif{
		URL:  urlString,
		Tags: []string{},
	}
	for _, st := range details.SplitContent[4:] {
		g.Tags = append(g.Tags, st)

	}
	if AddGif(vName, g.Tags, g, allVerbs) {
		m.Content = "Succesfully added gif to " + vName + "!"
	} else {
		m.Content = "I ran into some error..."
	}
	Store(*allVerbs)
	return m
}

// AddGif will add a given gif to goshleep's database. If the verb isn't found,
// creates a new one.
func AddGif(verbName string, tags []string, g Gif, allVerbs *[]Verb) bool {
	v, fuzz := getVerb(verbName, allVerbs)
	if Contains(v.Names, "unknown") || fuzz {
		v.Names = []string{verbName}
		v.Images = []Gif{g}
		*allVerbs = append(*allVerbs, *v)
	} else {
		v.Images = append(v.Images, g)
	}

	return true
}

// RemoveGifCommand is an admin command that allows admins to remove a gif.
func RemoveGifCommand(details *Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	log.Println("Backing up before modification...")
	Store(*allVerbs)
	m := baseMessage(details)
	vName := details.SplitContent[1]
	urlString := details.SplitContent[2]
	if RemoveGif(vName, urlString, allVerbs) {
		m.Content = "Succesfully removed gif from " + vName + "!"
	} else {
		m.Content = "I ran into some error..."
	}

	Store(*allVerbs)
	return m
}

// AddGif will add a given gif to goshleep's database. If the verb isn't found,
// creates a new one.
func RemoveGif(verbName string, url string, allVerbs *[]Verb) bool {
	v, fuzz := getVerb(verbName, allVerbs)

	if !Contains(v.Names, "unknown") || fuzz {
		for i, gif := range v.Images {
			if (gif.URL == url) {
				v.Images = remove(v.Images, i)
				return true
			}
		}
	} else {
		return false
	}


	return true
}

func remove(s []Gif, i int) []Gif {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func AddSynonymCommand(details *Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	log.Println("Backing up before modification...")
	Store(*allVerbs)

	m := baseMessage(details)
	vName := details.SplitContent[1]
	bName := details.SplitContent[2]
	m.Content = "Creating a synonym of " + bName + ", " + vName
	log.Println(m.Content)
	success := AddSynonym(vName, bName, allVerbs)
	if !success {
		m.Content = "Something went wrong..."
	}
	Store(*allVerbs)
	log.Println("Backing up after modification...")
	return m
}

func AddSynonym(verbName string, baseName string, allVerbs *[]Verb) bool {
	v, fuzz := getVerb(verbName, allVerbs)
	b, bfuzz := getVerb(baseName, allVerbs)
	// If the verb name doesn't exist BUT the base name does
	if (Contains(v.Names, "unknown") || fuzz) && (!Contains(b.Names, "unknown") && !bfuzz) {
		b.Names = append(b.Names, verbName)
	} else {
		return false
	}
	return true
}

func AddAdminCommand(details *Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	ReadViper()
	added := []string{}
	m := discordgo.MessageSend{}
	log.Println(len(details.SplitContent))
	for _, id := range details.SplitContent[1:] {
		isId, _ := regexp.MatchString(`^\d+$`, id)
		log.Println(id)
		if isId {
			if addAdmin(id) {
				idUser, _ := s.User(id)
				added = append(added, idUser.Mention())
			}
		}
	}
	for _, u := range details.dMessage.Mentions {
		if addAdmin(u.ID) {
			added = append(added, u.Mention())
		}
	}
	if len(added) > 0 {
		if len(added) > 1 {
			m.Content = "Added admins "
		} else {
			m.Content = "Added admin "
		}
		m.Content = m.Content + strings.Join(added, " and ")

	} else {
		m.Content = "This seems to be formatted incorrectly, or they are already admins."
	}

	m.Reference = details.dMessage.Reference()
	m.AllowedMentions = nil

	return m
}

func addAdmin(id string) bool {
	admins := viper.GetStringSlice("admins")
	notPresent := !Contains(admins, id)
	if notPresent {
		viper.Set("admins", append(admins, id))
	}
	viper.WriteConfig()
	return notPresent
}

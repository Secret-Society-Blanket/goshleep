package goshleep

import (
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

// AddGif is an admin command that allows admins to add a gif, and creates a verb if it does not exist.
func AddGifCommand(details *Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	log.Println("Backing up before modification...")
	Store(*allVerbs);
	m :=discordgo.MessageSend{};
	return m;
}

// AddGif will add a given gif to goshleep's database. If the verb isn't found,
// creates a new one.
func AddGif(verbName string, tags []string, g Gif, allVerbs *[]Verb) bool {
	v, fuzz := getVerb(verbName, allVerbs)
	if v.Name == "unknown" || fuzz {
		v.Name = verbName
		v.Images = []Gif{g};
		*allVerbs = append(*allVerbs, *v)
	} else {
		v.Images = append(v.Images, g);
	}

	return true;
}


func AddAdminCommand(details *Request, s *discordgo.Session, allVerbs *[]Verb) discordgo.MessageSend {
	ReadViper();
 
	added := []string{};
	m := discordgo.MessageSend{}
	log.Println(len(details.SplitContent))
	for _, id := range details.SplitContent[1:] {
		isId, _ := regexp.MatchString(`^\d+$`, id)
		log.Println(id)
		if (isId) {
			if (addAdmin(id)) {
				idUser, _ := s.User(id)
				added = append(added,idUser.Mention());
			}
		}
	}
	for _, u := range details.dMessage.Mentions {
		if (addAdmin(u.ID)) {
			added = append(added, u.Mention());
		}
	}
	if (len(added) > 0) {
		if (len(added) > 1) {
			m.Content = "Added admins "
		} else {
			m.Content = "Added admin "
		}
		m.Content = m.Content + strings.Join(added, " and ")

	} else {
		m.Content = "This seems to be formatted incorrectly, or they are already admins."
	}

	m.Reference = details.dMessage.Reference();
	m.AllowedMentions = nil;

	return m;
}

func addAdmin(id string) bool {
	admins := viper.GetStringSlice("admins")
	notPresent := !Contains(admins, id);
	if notPresent {
		viper.Set("admins", append(admins,id));
	}
	viper.WriteConfig();
	return notPresent;
}

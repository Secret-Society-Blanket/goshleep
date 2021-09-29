package goshleep

import (
	"strings"

	"log"
	"github.com/spf13/viper"
	"github.com/bwmarrin/discordgo"
)

// GetName takes a member and returns their name
func GetName(u *discordgo.Member) string {
	if u.Nick != "" {
		return u.Nick
	}
	return u.User.Username
}

// GetMentionNames takes a message and a session, and returns a
// list of user nicks mentioned within.
func GetMentionNames(m *discordgo.Message, s *discordgo.Session) string {

	var names []string
	for _, user := range m.Mentions {
		member, _ := s.GuildMember(m.GuildID, user.ID)
		names = append(names, GetName(member))
	}
	return strings.Join(names, " and ")

}

// Checks
func IfMatchHotStrings (arr []string, check string) bool{

	out := false

	for i := 0; i < len(arr); i++ {

		log.Println(arr[i])
		if (strings.HasPrefix(check, viper.GetString("prefix") + arr[i])) {
			out = true
		}
	}

	return out
}

func Contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true;
		}
	}
	return false;
}


func makeReference(slice []Gif)  []*Gif {
	out := []*Gif{};
	for _, g := range slice {
		out = append(out, &g)
	}
	return out;
}

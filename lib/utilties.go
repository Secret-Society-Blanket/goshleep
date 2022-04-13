package goshleep

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
)

// GetName takes a member and returns their name
func GetName(u *discordgo.User, guildID string, s *discordgo.Session) string {

	m, _ := s.GuildMember(guildID, u.ID)

	if m.Nick != "" {
		return m.Nick
	}

	return u.Username
}

// GetMentionNames takes a message and a session, and returns a
// list of user nicks mentioned within.
func GetMentionNames(m *discordgo.Message, s *discordgo.Session) string {

	var names []string
	for _, user := range m.Mentions {
		names = append(names, GetName(user, m.GuildID, s))
	}
	return strings.Join(names, " and ")

}

// Checks
func IfMatchHotStrings(arr []string, check string) bool {

	out := false

	for i := 0; i < len(arr); i++ {

		if strings.HasPrefix(check, viper.GetString("prefix")+arr[i]) {
			out = true
		}
	}

	return out
}

func Contains(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}

func makeReference(slice []Gif) []*Gif {
	out := []*Gif{}
	for i := range slice {
		out = append(out, &slice[i])
	}
	return out
}

package main

import (
	"strings"

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

func IfInString (arr []string, check string) bool{

	out := false

	for i := 0; i < len(arr); i++ {

		if (arr[i] == check) {
			out = true
		}
	}

	return out
}

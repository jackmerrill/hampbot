package listeners

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jackmerrill/hampbot/internal/utils/config"
)

type MessageDeleteListener struct{}

const CHANNEL_ID = "1150951087323496450"

func (l *MessageDeleteListener) Exec(s *discordgo.Session, e *discordgo.MessageDelete) {
	var msg = discordgo.Message{}

	fields := []*discordgo.MessageEmbedField{}

	if m, ok := config.MessageLog[e.ID]; ok {
		msg = m
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Author",
			Value:  fmt.Sprintf("<@%s>", msg.Author.ID),
			Inline: true,
		})
	} else {
		msg.Content = "Unknown, untracked message."
		msg.ChannelID = e.ChannelID
	}

	if msg.Author.Bot {
		return
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "Channel",
		Value:  fmt.Sprintf("<#%s>", e.ChannelID),
		Inline: true,
	})

	_, err := s.ChannelMessageSendComplex(CHANNEL_ID, &discordgo.MessageSend{
		Content: "Message deleted",
		Embed: &discordgo.MessageEmbed{
			Title:       "Message deleted",
			Description: msg.Content,
			Fields:      fields,
			Color:       0xff0000,
		},
	})
	if err != nil {
		panic(err)
	}
}

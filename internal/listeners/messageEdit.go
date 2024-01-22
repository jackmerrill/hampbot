package listeners

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jackmerrill/hampbot/internal/utils/config"
)

type MessageEditListener struct{}

func (l *MessageEditListener) Exec(s *discordgo.Session, e *discordgo.MessageUpdate) {
	if e.Author.Bot {
		return
	}
	var old discordgo.Message
	fields := []*discordgo.MessageEmbedField{
		{
			Name:   "Author",
			Value:  fmt.Sprintf("<@%s>", e.Author.ID),
			Inline: true,
		},
		{
			Name:   "Channel",
			Value:  fmt.Sprintf("<#%s>", e.ChannelID),
			Inline: true,
		},
	}

	if o, ok := config.MessageLog[e.ID]; ok {
		old = o

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Old",
			Value:  old.Content,
			Inline: false,
		})
	}

	fields = append(fields, &discordgo.MessageEmbedField{
		Name:   "New",
		Value:  e.Content,
		Inline: false,
	})

	var image *discordgo.MessageEmbedImage

	if len(e.Attachments) > 0 {
		image = &discordgo.MessageEmbedImage{
			URL: e.Attachments[0].URL,
		}
	}

	_, err := s.ChannelMessageSendComplex(CHANNEL_ID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Title:  "Message edited",
			Fields: fields,
			Description: fmt.Sprintf(
				"[Jump to message](https://discordapp.com/channels/%s/%s/%s)",
				e.GuildID, e.ChannelID, e.ID,
			),
			Color: 0xffff00,
			Image: image,
		},
	})

	if err != nil {
		panic(err)
	}
}

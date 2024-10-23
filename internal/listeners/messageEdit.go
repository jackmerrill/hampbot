package listeners

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/sergi/go-diff/diffmatchpatch"
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

	diffText := ""

	if o, ok := config.MessageLog[e.ID]; ok {
		old = o

		dmp := diffmatchpatch.New()

		fileAdmp, fileBdmp, dmpStrings := dmp.DiffLinesToChars(old.Content, e.Content)
		diffs := dmp.DiffMain(fileAdmp, fileBdmp, false)
		diffs = dmp.DiffCharsToLines(diffs, dmpStrings)
		diffs = dmp.DiffCleanupSemantic(diffs)

		for _, d := range diffs {
			if d.Type == diffmatchpatch.DiffInsert {
				diffText += fmt.Sprintf("+ %s\n", d.Text)
			}
			if d.Type == diffmatchpatch.DiffDelete {
				diffText += fmt.Sprintf("- %s\n", d.Text)
			}
		}

	}

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
				"[Jump to message](https://discordapp.com/channels/%s/%s/%s)\n\n```diff\n%s```",
				e.GuildID, e.ChannelID, e.ID, diffText,
			),
			Color: 0xffff00,
			Image: image,
		},
	})

	if err != nil {
		panic(err)
	}
}

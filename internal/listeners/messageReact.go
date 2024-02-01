package listeners

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type MessageReactAddListener struct{}

const STARBOARD_CHANNEL_ID = "1202408436969766982"
const STARBOARD_THRESHOLD = 3

func (l *MessageReactAddListener) Exec(s *discordgo.Session, e *discordgo.MessageReactionAdd) {
	if e.MessageReaction.Emoji.Name == "⭐" {
		msg, err := s.ChannelMessage(e.MessageReaction.ChannelID, e.MessageReaction.MessageID)
		if err != nil {
			fmt.Println(err)
			return
		}

		for _, reaction := range msg.Reactions {
			if reaction.Emoji.Name == "⭐" {
				if reaction.Count >= STARBOARD_THRESHOLD {
					m, err := findMessageWithIDInFooter(s, STARBOARD_CHANNEL_ID, msg.ID)

					if err != nil {
						// not found

						content := fmt.Sprintf("**⭐ %d | <#%s>**", reaction.Count, msg.ChannelID)
						var image *discordgo.MessageEmbedImage

						if len(msg.Attachments) > 0 {
							image = &discordgo.MessageEmbedImage{
								URL: msg.Attachments[0].URL,
							}
						}

						embed := &discordgo.MessageEmbed{
							Author: &discordgo.MessageEmbedAuthor{
								Name:    msg.Author.Username,
								IconURL: msg.Author.AvatarURL(""),
							},
							Description: fmt.Sprintf("%s\n\n[Jump to message](https://discordapp.com/channels/%s/%s/%s)", msg.Content, msg.GuildID, msg.ChannelID, msg.ID),
							Image:       image,
							Footer: &discordgo.MessageEmbedFooter{
								Text: fmt.Sprintf("ID: %s", msg.ID),
							},
							Timestamp: msg.Timestamp.Format("2006-01-02T15:04:05-0700"),
						}

						_, err := s.ChannelMessageSendComplex(STARBOARD_CHANNEL_ID, &discordgo.MessageSend{
							Content: content,
							Embed:   embed,
						})

						if err != nil {
							fmt.Println(err)
						}
					} else {
						// found msg

						content := fmt.Sprintf("**⭐ %d | <#%s>**", reaction.Count, msg.ChannelID)

						_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
							Channel: STARBOARD_CHANNEL_ID,
							ID:      m.ID,
							Content: &content,
							Embeds:  m.Embeds,
						})

						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}
	}
}

type MessageReactRemoveListener struct{}

func (l *MessageReactRemoveListener) Exec(s *discordgo.Session, e *discordgo.MessageReactionRemove) {
	if e.MessageReaction.Emoji.Name == "⭐" {
		msg, err := s.ChannelMessage(e.MessageReaction.ChannelID, e.MessageReaction.MessageID)
		if err != nil {
			fmt.Println(err)
			return
		}

		// check if the message has any star reactions
		starReactions := 0

		for _, reaction := range msg.Reactions {
			if reaction.Emoji.Name == "⭐" {
				starReactions = reaction.Count
				if reaction.Count < STARBOARD_THRESHOLD {
					m, err := findMessageWithIDInFooter(s, STARBOARD_CHANNEL_ID, msg.ID)

					if err == nil {
						// found msg
						err = s.ChannelMessageDelete(STARBOARD_CHANNEL_ID, m.ID)
						if err != nil {
							fmt.Println(err)
						}
					}
				} else {
					m, err := findMessageWithIDInFooter(s, STARBOARD_CHANNEL_ID, msg.ID)

					if err == nil {
						// found msg
						content := fmt.Sprintf("**⭐ %d | <#%s>**", reaction.Count, msg.ChannelID)

						_, err := s.ChannelMessageEditComplex(&discordgo.MessageEdit{
							Channel: STARBOARD_CHANNEL_ID,
							ID:      m.ID,
							Content: &content,
							Embeds:  m.Embeds,
						})

						if err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}

		if starReactions < STARBOARD_THRESHOLD {
			m, err := findMessageWithIDInFooter(s, STARBOARD_CHANNEL_ID, msg.ID)

			if err == nil {
				// found msg
				err = s.ChannelMessageDelete(STARBOARD_CHANNEL_ID, m.ID)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}

func getAllMessages(session *discordgo.Session, channelID string) ([]*discordgo.Message, error) {
	var allMessages []*discordgo.Message
	var lastID string

	for {
		// Fetch a batch of 100 messages
		messages, err := session.ChannelMessages(channelID, 100, lastID, "", "")
		if err != nil {
			return nil, err
		}

		// Break the loop if there are no more messages
		if len(messages) == 0 {
			break
		}

		// Append the messages to the allMessages slice
		allMessages = append(allMessages, messages...)

		// Update lastID to the last message's ID in the batch
		lastID = messages[len(messages)-1].ID
	}

	return allMessages, nil
}

func findMessageWithIDInFooter(session *discordgo.Session, channelID string, searchID string) (*discordgo.Message, error) {
	messages, err := getAllMessages(session, channelID)
	if err != nil {
		return nil, err
	}

	for _, message := range messages {
		for _, embed := range message.Embeds {
			if embed.Footer != nil && strings.Contains(embed.Footer.Text, searchID) {
				return message, nil
			}
		}
	}

	return nil, fmt.Errorf("No message found with the specified ID in footer")
}

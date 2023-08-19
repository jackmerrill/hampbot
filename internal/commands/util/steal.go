package commands

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	"github.com/zekroTJA/shireikan"
)

type Steal struct {
}

func (c *Steal) GetInvokes() []string {
	return []string{"steal", "emote"}
}

func (c *Steal) GetDescription() string {
	return "Steal an emote from another server"
}

func (c *Steal) GetHelp() string {
	return "`steal [emote(s)]` - steal an emote from another server\n`steal [emote(s)]` - steal an emote (or two) from another server"
}

func (c *Steal) GetGroup() string {
	return config.GroupUtil
}

func (c *Steal) GetDomainName() string {
	return "hamp.util.steal"
}

func (c *Steal) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}
func (c *Steal) IsExecutableInDMChannels() bool {
	return false
}

func (c *Steal) Exec(ctx shireikan.Context) error {
	msg := ctx.GetMessage()

	emotes := msg.GetCustomEmojis()

	if len(emotes) == 0 {
		ctx.GetSession().ChannelMessageSendEmbed(msg.ChannelID, embed.NewErrorEmbed(ctx).SetTitle("No emotes found").SetDescription("No emotes were found in your message.").MessageEmbed)
		return nil
	}

	var newEmojis []string

	for _, emote := range emotes {
		// download the emote
		res, err := http.Get(fmt.Sprintf("https://cdn.discordapp.com/emojis/%s", emote.ID))

		if err != nil {
			ctx.GetSession().ChannelMessageSendEmbed(msg.ChannelID, embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while downloading your emote.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}

		defer res.Body.Close()

		data, err := ioutil.ReadAll(res.Body)

		if err != nil {
			ctx.GetSession().ChannelMessageSendEmbed(msg.ChannelID, embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while downloading your emote.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}

		// convert to base64
		base64Emote := base64.StdEncoding.EncodeToString(data)

		emoji := discordgo.EmojiParams{
			Name:  emote.Name,
			Image: "data:image/png;base64," + base64Emote,
			Roles: nil,
		}

		_, err = ctx.GetSession().GuildEmojiCreate(config.BotGuild, &emoji)

		if err != nil {
			ctx.GetSession().ChannelMessageSendEmbed(msg.ChannelID, embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while uploading your emote.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}

		newEmojis = append(newEmojis, emote.Name)
	}

	ctx.GetSession().ChannelMessageSendEmbed(msg.ChannelID, embed.NewSuccessEmbed(ctx).SetTitle("Emote(s) stolen!").SetDescription(fmt.Sprintf("The following emotes were stolen: `%s`", strings.Join(newEmojis, ", "))).MessageEmbed)

	return nil
}

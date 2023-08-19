package commands

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	"github.com/zekroTJA/shireikan"
)

// Ping is a command responding with a ping
// message in the commands channel.
type Ping struct {
}

// GetInvoke returns the command invokes.
func (c *Ping) GetInvokes() []string {
	return []string{"ping", "p"}
}

// GetDescription returns the commands description.
func (c *Ping) GetDescription() string {
	return "ping pong"
}

// GetHelp returns the commands help text.
func (c *Ping) GetHelp() string {
	return "`ping` - ping"
}

// GetGroup returns the commands group.
func (c *Ping) GetGroup() string {
	return config.GroupUtil
}

// GetDomainName returns the commands domain name.
func (c *Ping) GetDomainName() string {
	return "hamp.util.ping"
}

// GetSubPermissionRules returns the commands sub
// permissions array.
func (c *Ping) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

// IsExecutableInDMChannels returns whether
// the command is executable in DM channels.
func (c *Ping) IsExecutableInDMChannels() bool {
	return true
}

// Exec is the commands execution handler.
func (c *Ping) Exec(ctx shireikan.Context) error {
	start := time.Now()

	m, err := ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
		Reference: ctx.GetMessage().Reference(),
		Embed:     embed.NewWarningEmbed(ctx).SetTitle("Pinging... <a:whencat:993983100805730415>").MessageEmbed,
	})

	end := time.Now()
	diff := end.Sub(start)

	if err != nil {
		log.Error("Failed sending message: ", err)
	}

	_, err = ctx.GetSession().ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel: ctx.GetChannel().ID,
		ID:      m.ID,
		Embed:   embed.NewSuccessEmbed(ctx).SetTitle("Pong! :ping_pong:").SetDescription(fmt.Sprintf(":ping_pong: Gateway Ping: `%dms`\n:desktop: API Ping: `%dms`", ctx.GetSession().HeartbeatLatency().Milliseconds(), diff.Milliseconds())).MessageEmbed,
	})

	return err
}

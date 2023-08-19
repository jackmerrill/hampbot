package studentlife

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	"github.com/zekroTJA/shireikan"
)

type MachineType string

const (
	Washer MachineType = "Washer"
	Dryer  MachineType = "Dryer"
)

type Machine struct {
	Name          string      `json:"name"`
	Type          MachineType `json:"type"`
	Status        string      `json:"status"`
	Time          *time.Time  `json:"time"`
	EstimatedTime *time.Time  `json:"estimatedTime"`
}

type LaundryRoom struct {
	ID         string    `json:"id"`
	Name       string    `json:"name"`
	NextUpdate time.Time `json:"nextUpdate"`
	LastUpdate time.Time `json:"lastUpdate"`

	Machines []Machine `json:"machines"`

	updateChan chan bool
}

type Laundry struct {
}

func (c *Laundry) GetInvokes() []string {
	return []string{"laundry"}
}

func (c *Laundry) GetDescription() string {
	return "Get the laundry status of a residence hall"
}

func (c *Laundry) GetHelp() string {
	return "`laundry [dakin/merrill/prescott/enfield]` - Laundry help"
}

func (c *Laundry) GetGroup() string {
	return config.GroupStudentLife
}

func (c *Laundry) GetDomainName() string {
	return "hamp.studentlife.laundry"
}

func (c *Laundry) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}
func (c *Laundry) IsExecutableInDMChannels() bool {
	return true
}

func (c *Laundry) Exec(ctx shireikan.Context) error {
	building := ctx.GetArgs().Get(0).AsString()

	if building == "" {
		ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Embed:     embed.NewErrorEmbed(ctx).SetDescription("Please specify a building to get the laundry status of.\n\n**Usage:** `laundry [dakin/merrill/prescott/enfield]`").MessageEmbed,
			Reference: ctx.GetMessage().Reference(),
		})
		return nil
	}

	building = strings.ToLower(building)

	if building != "dakin" && building != "merrill" && building != "prescott" && building != "enfield" {
		ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Embed:     embed.NewErrorEmbed(ctx).SetDescription("Please specify a valid building to get the laundry status of.\n\n**Usage:** `laundry [dakin/merrill/prescott/enfield]`").MessageEmbed,
			Reference: ctx.GetMessage().Reference(),
		})
		return nil
	}

	res, err := http.Get(fmt.Sprintf("%s/api/utilities/laundry/%s", config.HampAPI, ctx.GetArgs().Get(0)))

	if err != nil {
		return err
	}

	var room LaundryRoom

	err = json.NewDecoder(res.Body).Decode(&room)

	if err != nil {
		return err
	}

	embed := embed.NewSuccessEmbed(ctx).SetTitle(fmt.Sprintf("%s Laundry Status :shirt:", room.Name)).SetDescription(fmt.Sprintf("Last update: %s", config.ConvertTimestampToDiscordTimestampWithFormat(room.LastUpdate, "T")))

	for _, machine := range room.Machines {
		if machine.Status == "Available" {
			embed.AddField(fmt.Sprintf("%s - %s", machine.Name, machine.Type), fmt.Sprintf(":green_circle: **Status:** %s", machine.Status), false)
		} else if machine.Status == "Not online" {
			embed.AddField(fmt.Sprintf("%s - %s", machine.Name, machine.Type), fmt.Sprintf(":red_circle: **Status:** %s", machine.Status), false)
		} else if machine.Status == "End of cycle" {
			embed.AddField(fmt.Sprintf("%s - %s", machine.Name, machine.Type), fmt.Sprintf(":blue_circle: **Status:** %s", machine.Status), false)
		} else {
			embed.AddField(fmt.Sprintf("%s - %s", machine.Name, machine.Type), fmt.Sprintf(":yellow_circle: **Status:** %s\n:alarm_clock: **Time Remaining:** %s", machine.Status, config.ConvertTimestampToDiscordTimestampWithFormat(*machine.EstimatedTime, "R")), false)
		}
	}

	ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
		Embed:     embed.MessageEmbed,
		Reference: ctx.GetMessage().Reference(),
	})

	return nil
}

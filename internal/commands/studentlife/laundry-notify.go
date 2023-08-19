package studentlife

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/charmbracelet/log"
	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	"github.com/zekroTJA/shireikan"
	"golang.org/x/net/websocket"
)

type LaundryNotify struct {
}

func (c *LaundryNotify) GetInvokes() []string {
	return []string{"laundrynotify", "ln"}
}

func (c *LaundryNotify) GetDescription() string {
	return "(Beta) Be notified when laundry machines are available or done"
}

func (c *LaundryNotify) GetHelp() string {
	return "`laundrynotify [building] [machine|any] <type>`\n\n**Example:** `laundrynotify dakin 01` `ln merrill any washer`"
}

func (c *LaundryNotify) GetGroup() string {
	return config.GroupStudentLife
}

func (c *LaundryNotify) GetDomainName() string {
	return "hamp.util.laundrynotify"
}

func (c *LaundryNotify) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}
func (c *LaundryNotify) IsExecutableInDMChannels() bool {
	return true
}

type NotifiedMachine struct {
	Building    string
	Machine     string
	MachineType *string
	WaitingFor  string // Available, Done
}

var NotifiedMachines map[string]NotifiedMachine = make(map[string]NotifiedMachine)

func NotifyUser(ctx *discordgo.Session, userId string, machine Machine, building string, waitingFor string) {
	var machineType string

	if machine.Type == Washer {
		machineType = "washer"
	} else {
		machineType = "dryer"
	}

	if waitingFor == "Available" && machine.Status == "Available" {
		ctx.ChannelMessageSendComplex(userId, &discordgo.MessageSend{
			Embed: embed.NewEmbed().SetTitle("Laundry Machine Available").SetDescription(fmt.Sprintf("The %s in %s is now available.", machineType, building)).SetColor(0x00ff00).MessageEmbed,
		})
	} else if waitingFor == "Done" && machine.Status == "End of cycle" {
		ctx.ChannelMessageSendComplex(userId, &discordgo.MessageSend{
			Embed: embed.NewEmbed().SetTitle("Laundry Machine Done").SetDescription(fmt.Sprintf("The %s in %s is now done.", machineType, building)).SetColor(0x00ff00).MessageEmbed,
		})
	}
}

func wsLoop(client *websocket.Conn, ctx *discordgo.Session) {
	for {
		msg := make([]byte, 4096) // 4kb to be safe lol
		n, err := client.Read(msg)
		if err != nil {
			log.Error(err)
			return
		}

		msg = msg[:n]

		var room LaundryRoom

		err = json.Unmarshal(msg, &room)

		if err != nil {
			log.Error(err)
			return
		}

		for _, machine := range room.Machines {
			for userId, notifiedMachine := range NotifiedMachines {
				if notifiedMachine.Building == room.Name {
					if notifiedMachine.Machine == "any" { // requesting any machine with specific type
						if notifiedMachine.MachineType != nil {
							if string(machine.Type) == *notifiedMachine.MachineType {
								NotifyUser(ctx, userId, machine, room.Name, notifiedMachine.WaitingFor)
							}
						}
					} else {
						if notifiedMachine.Machine == machine.Name { // doesn't matter, they're specifying a specific machine
							NotifyUser(ctx, userId, machine, room.Name, notifiedMachine.WaitingFor)
						}
					}
				}
			}
		}
	}
}

func InitDakin(ctx *discordgo.Session) error {
	u := url.URL{Scheme: "wss", Host: config.HampAPIHost, Path: "/api/utilities/laundry/dakin/live"}

	wsClient, err := websocket.Dial(u.String(), "", "http://localhost")

	if err != nil {
		log.Error(err)
		return err
	}

	go wsLoop(wsClient, ctx)
	return nil

}

func InitMerrill(ctx *discordgo.Session) error {
	u := url.URL{Scheme: "wss", Host: config.HampAPIHost, Path: "/api/utilities/laundry/merrill/live"}

	wsClient, err := websocket.Dial(u.String(), "", "http://localhost")

	if err != nil {
		log.Error(err)
		return err
	}

	go wsLoop(wsClient, ctx)
	return nil

}

func InitPrescott(ctx *discordgo.Session) error {
	u := url.URL{Scheme: "wss", Host: config.HampAPIHost, Path: "/api/utilities/laundry/prescott/live"}

	wsClient, err := websocket.Dial(u.String(), "", "http://localhost")

	if err != nil {
		log.Error(err)
		return err
	}

	go wsLoop(wsClient, ctx)
	return nil

}

func InitEnfield(ctx *discordgo.Session) error {
	u := url.URL{Scheme: "wss", Host: config.HampAPIHost, Path: "/api/utilities/laundry/enfield/live"}

	wsClient, err := websocket.Dial(u.String(), "", "http://localhost")

	if err != nil {
		log.Error(err)
		return err
	}

	go wsLoop(wsClient, ctx)
	return nil

}

func InitAll(ctx *discordgo.Session) error {
	err := InitDakin(ctx)

	if err != nil {
		return err
	}

	err = InitMerrill(ctx)

	if err != nil {
		return err
	}

	err = InitPrescott(ctx)

	if err != nil {
		return err
	}

	err = InitEnfield(ctx)

	if err != nil {
		return err
	}

	return nil
}

func (c *LaundryNotify) Exec(ctx shireikan.Context) error {
	building := ctx.GetArgs().Get(0).AsString()

	if building == "" {
		ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Embed:     embed.NewErrorEmbed(ctx).SetDescription("Please specify a building to get the laundry status of.\n\n**Usage:** `laundrynotify [building] [machine|any] <type>`").MessageEmbed,
			Reference: ctx.GetMessage().Reference(),
		})
		return nil
	}

	building = strings.ToLower(building)

	if building != "dakin" && building != "merrill" && building != "prescott" && building != "enfield" {
		ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Embed:     embed.NewErrorEmbed(ctx).SetDescription("Please specify a valid building to get the laundry status of.\n\n**Usage:** `laundrynotify [building] [machine|any] <type>`").MessageEmbed,
			Reference: ctx.GetMessage().Reference(),
		})
		return nil
	}

	machine := ctx.GetArgs().Get(1).AsString()

	if machine == "" {
		ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Embed:     embed.NewErrorEmbed(ctx).SetDescription("Please specify a machine to get the laundry status of.\n\n**Usage:** `laundrynotify [building] [machine|any] <type>`").MessageEmbed,
			Reference: ctx.GetMessage().Reference(),
		})
		return nil
	}

	machine = strings.ToLower(machine)

	if machine != "any" {
		if len(machine) == 1 {
			machine = fmt.Sprintf("0%s", machine)
		}
	}

	machineType := ctx.GetArgs().Get(2).AsString()

	notifiedMachine := NotifiedMachine{
		Building: building,
		Machine:  machine,
	}

	if machineType == "" {
		if machine == "any" {
			ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
				Embed:     embed.NewErrorEmbed(ctx).SetDescription("Please specify a machine type to get the laundry status of.\n\n**Usage:** `laundrynotify [building] [machine|any] <type>`").MessageEmbed,
				Reference: ctx.GetMessage().Reference(),
			})

			return nil
		} else {
			notifiedMachine.MachineType = nil
			notifiedMachine.WaitingFor = "Done"
		}
	} else {
		machineType = strings.ToLower(machineType)
		notifiedMachine.MachineType = &machineType
	}

	if machine == "any" {
		notifiedMachine.WaitingFor = "Available"
		machine = fmt.Sprintf("any %s", *notifiedMachine.MachineType)
	} else {
		notifiedMachine.WaitingFor = "Done"
	}

	NotifiedMachines[ctx.GetUser().ID] = notifiedMachine

	ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
		Embed: embed.NewSuccessEmbed(ctx).SetTitle("Laundry Notification Set").SetDescription(fmt.Sprintf("You will be notified when `%s` in %s is %s.", machine, building, notifiedMachine.WaitingFor)).MessageEmbed,
	})

	return nil
}

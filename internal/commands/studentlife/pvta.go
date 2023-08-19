package studentlife

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	pvtago "github.com/jackmerrill/pvta-go"
	"github.com/zekroTJA/shireikan"
)

type PVTA struct {
}

func (c *PVTA) GetInvokes() []string {
	return []string{"pvta", "bus"}
}

func (c *PVTA) GetDescription() string {
	return "Get the details and location of a PVTA bus"
}

func (c *PVTA) GetHelp() string {
	return "`pvta [route] <vehicle>` - Get the details and location of a PVTA bus"
}

func (c *PVTA) GetGroup() string {
	return config.GroupUtil
}

func (c *PVTA) GetDomainName() string {
	return "hamp.studentlife.pvta"
}

func (c *PVTA) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}

func (c *PVTA) IsExecutableInDMChannels() bool {
	return true
}

func getRouteByShortName(shortName string) (*pvtago.RouteDetail, error) {
	client := pvtago.NewClient()

	routes, err := client.GetRoutes()

	if err != nil {
		return nil, err
	}

	for _, route := range routes {
		if route.ShortName == shortName {
			return &route, nil
		}
	}

	return nil, nil
}

func (c *PVTA) Exec(ctx shireikan.Context) error {
	client := pvtago.NewClient()

	route := ctx.GetArgs().Get(0).AsString()

	if route == "" {
		routes, err := client.GetVisibleRoutes()

		if err != nil {
			return err
		}

		var routeIds []string

		for _, route := range routes {
			if route.ShortName != "" {
				routeIds = append(routeIds, route.ShortName)
			}
		}

		ctx.ReplyEmbed(embed.NewSuccessEmbed(ctx).SetTitle("Available PVTA routes").SetDescription(strings.Join(routeIds, ", ")).MessageEmbed)

		return nil
	}

	vehicle := ctx.GetArgs().Get(1).AsString()

	if vehicle == "" {
		route, err := getRouteByShortName(route)

		if err != nil {
			ctx.ReplyEmbed(embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while fetching the route.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}

		vehicles, err := client.GetVehiclesForRoute(route.RouteID)

		if err != nil {
			return err
		}

		e := embed.NewSuccessEmbed(ctx).SetTitle(":bus: PVTA Route Details").SetDescription(fmt.Sprintf("**Route:** %s\n**Name:** %s", route.ShortName, route.LongName))

		tileProvider := &sm.TileProvider{
			Name:        "thunderforest-transit",
			Attribution: "Maps (c) Thundeforest; Data (c) OSM and contributors, ODbL",
			TileSize:    256,
			URLPattern:  "https://%[1]s.tile.thunderforest.com/transport/%[2]d/%[3]d/%[4]d.png?apikey=" + os.Getenv("TF_API_KEY"),
			Shards:      []string{"a", "b", "c"},
		}

		mapCtx := sm.NewContext()
		mapCtx.SetSize(600, 400)
		mapCtx.SetTileProvider(tileProvider)

		for _, vehicle := range vehicles {
			e.AddField(fmt.Sprintf("Vehicle #`%d`", vehicle.VehicleID), fmt.Sprintf(":round_pushpin: **Destination:** %s (%s)\n:globe_with_meridians: **Lat/Long:** `%f` `%f`", vehicle.Destination, vehicle.Direction, vehicle.Latitude, vehicle.Longitude), false)

			rgb, err := config.Hex2RGB(config.Hex(route.Color))

			if err != nil {
				return err
			}

			mapCtx.AddObject(
				sm.NewMarker(
					s2.LatLngFromDegrees(vehicle.Latitude, vehicle.Longitude),
					color.RGBA{R: rgb.Red, G: rgb.Green, B: rgb.Blue, A: 255},
					16.0,
				),
			)
		}

		img, err := mapCtx.Render()

		if err != nil {
			return err
		}

		// convert to base64
		buf := new(bytes.Buffer)

		err = png.Encode(buf, img)

		if err != nil {
			return err
		}

		e.SetImage("attachment://map.png")

		_, err = ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Embed:     e.MessageEmbed,
			Reference: ctx.GetMessage().Reference(),
			Files: []*discordgo.File{
				{
					Name:        "map.png",
					ContentType: "image/png",
					Reader:      buf,
				},
			},
		})

		if err != nil {
			return err
		}

		return nil
	}

	return nil
}

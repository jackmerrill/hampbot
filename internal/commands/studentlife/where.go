package studentlife

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	sm "github.com/flopp/go-staticmaps"
	"github.com/golang/geo/s2"
	mapbox "github.com/jackmerrill/go-mapbox/lib"
	"github.com/jackmerrill/go-mapbox/lib/base"
	"github.com/jackmerrill/go-mapbox/lib/directions"
	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	"github.com/zekroTJA/shireikan"
)

// string is the name of the building in regex, []float64 is the lat/long
type Location map[string][]float64

var LOCATIONS = Location{
	"^(kern|\"kern kafe\"|\"kern cafe\"|\"kern center\"|\"r\\.w\\. kern center\"|\"rw kern center\"|admissions|\"financial aid\"|finaid)$": []float64{42.325490, -72.530425},
	"^(bridge|\"bridge cafe\"|the bridge)$": []float64{42.32560055268511, -72.53171782489403},
	"^(rcc|\"robert crown center\"|gym)$":   []float64{42.3259021736783, -72.53138568980566},
	"^(library|\"harold f. johnson library center\"|\"mail room\"|\"post office\"|hampstore|\"hamp store\"|\"campus store\"|duplications|it|\"art gallery\"|\"hampshire college art gallery\"|sparc)$": []float64{42.325503815980404, -72.53234665636477},
	"^(cole|\"the cole\"|\"cole science center\"|csc)$": []float64{42.325005132910455, -72.53264274547023},
	"^(\"mixed nuts\")$":             []float64{42.32438520935449, -72.53314741618625},
	"^(prescott)$":                   []float64{42.32370121438042, -72.53400755941487},
	"^(\"prescott tavern\"|tavern)$": []float64{42.32330551141666, -72.53410313088433},
	"^(c4d|\"center for design\"|\"lemelson building\")$": []float64{42.323653164859145, -72.53277277596145},
	"^(\"central records\"|casa)$":                        []float64{42.32400565840685, -72.53267685200223},
	"^(\"arts barn\"|\"art barn\"|art)$":                  []float64{42.323484930985096, -72.53270230871183},
	"^(jlc|\"jerome liebling center\")$":                  []float64{42.32343637951593, -72.53186281901529},
	"^(\"solar canopy\")$":                                []float64{42.3231160046793, -72.5323111913881},
	"^(\"music and dance building\"|mdb|music|dance)$":    []float64{42.32296880488658, -72.53257385728799},
	"^(ash|\"adele simmons hall\")$":                      []float64{42.32286199484466, -72.53189407069387},
	"^(carle|\"the carle\"|\"the eric carle museum\"|\"the eric carle museum of picture book art\")$": []float64{42.32109682435725, -72.53332586596198},
	"^(elc|early learning center)$":                                  []float64{42.32130800774271, -72.53485513102083},
	"^(\"multisport center\"|multisport|weight room|msc)$":           []float64{42.32138500151223, -72.53614340300827},
	"^(\"dakin student life center\"|dslc)$":                         []float64{42.32313128927387, -72.53022445236283},
	"^(\"merrill student life center\"|mslc)$":                       []float64{42.32347857352905, -72.53034306624653},
	"^(\"merrill pavillion\"|pavillion)$":                            []float64{42.323322471247494, -72.53030510980373},
	"^(dakin|hell|\"dakin house\")$":                                 []float64{42.322597743840625, -72.53027825479788},
	"^(merrill|\"merrill house\")$":                                  []float64{42.323732718386175, -72.52977892147766},
	"^(\"dining commons\"|dc|saga)$":                                 []float64{42.32319356233955, -72.52917576638322},
	"^(\"yiddish book center\"|ybc)$":                                []float64{42.32175452980474, -72.5276151025249},
	"^(\"franklin patterson hall\"|fph)$":                            []float64{42.324221538083805, -72.5306271019939},
	"^(\"the yurt\"|yurt|radio)$":                                    []float64{42.324071744881785, -72.53136719263956},
	"^(enfield|\"enfield mods\"|\"enfield house\")$":                 []float64{42.32646094303269, -72.52929170510444},
	"^(\"wellness center\")$":                                        []float64{42.32709812095117, -72.5291661607627},
	"^(\"spiritual life center\"|slc)$":                              []float64{42.32705614551552, -72.52973716786117},
	"^(edh|\"emily dickinson hall\"|theatre|theater)$":               []float64{42.327652942240896, -72.53062317253445},
	"^(\"writing center\")$":                                         []float64{42.32780110830725, -72.53112929268784},
	"^(greenwich|\"greenwich mods\"|\"greenwich house\")$":           []float64{42.32748781455265, -72.53192281696715},
	"^(soccer|\"soccer field\"|\"hampshire college soccer field\")$": []float64{42.32651025401564, -72.53454675256307},
	"^(\"cultural center\"|cc|\"lebron-wiggins-pran cultural center\"|\"Lebrón-Wiggins-Pran Cultural Center\")$": []float64{42.32487763771362, -72.5339769868106},
	"^(basketball|\"basketball courts\"|\"basketball court\"|\"tennis courts\"|\"tennis court\")$":               []float64{42.32557748709811, -72.53720546848393},
	"^(\"red barn\"|\"the red barn\"|barn)$":                                                     []float64{42.32642633315529, -72.52551730246282},
	"^(\"the hitchcock center\"|\"hitchcock center\"|\"hitchcock center for the environment\")$": []float64{42.32771760017527, -72.52570886249646},
	"^(csa|farm|\"hampshire college farm center\"|\"hampshire college farm\")$":                  []float64{42.32903799744424, -72.52573139895138},
	"^(\"health services\"|\"hampshire college health services\")$":                              []float64{42.32666929509394, -72.52524995366625},
	"^(atkins|\"atkins farms country market\")$":                                                 []float64{42.319355104827295, -72.52927125654328},
	"^(res|\"the res\"|the reservoir|reservoir)$":                                                []float64{42.317285834368995, -72.5406219720153},
}

type Where struct {
}

var mapboxToken = os.Getenv("MAPBOX_TOKEN")
var mapBox, _ = mapbox.NewMapbox(mapboxToken)

func (c *Where) GetInvokes() []string {
	return []string{"where"}
}

func (c *Where) GetDescription() string {
	return "Where is this building?"
}

func (c *Where) GetHelp() string {
	return "`where <building>` - Where is this building?"
}

func (c *Where) GetGroup() string {
	return config.GroupStudentLife
}

func (c *Where) GetDomainName() string {
	return "hamp.studentlife.where"
}

func (c *Where) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}
func (c *Where) IsExecutableInDMChannels() bool {
	return true
}

func (c *Where) Exec(ctx shireikan.Context) error {
	if len(ctx.GetArgs()) < 1 {
		// return all buildings
		e := embed.NewSuccessEmbed(ctx)

		e.SetTitle("Buildings")

		desc := ""
		locs := []string{}

		for regex := range LOCATIONS {
			l := strings.TrimLeft(regex, "^(")
			l = strings.TrimRight(l, ")$")

			ls := strings.Split(l, "|")

			for _, l := range ls {
				locs = append(locs, l)
			}

			desc += fmt.Sprintf("`%s`, ", strings.Join(ls, "`, `"))
		}

		e.SetDescription(desc)

		_, err := ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Embed:     e.MessageEmbed,
			Reference: ctx.GetMessage().Reference(),
		})

		return err
	}

	to := strings.ToLower(ctx.GetArgs().Get(0).AsString())
	from := strings.ToLower(ctx.GetArgs().Get(1).AsString())

	e := embed.NewSuccessEmbed(ctx)

	e.SetTitle("Where is " + to + "?")

	if from != "" {
		e.SetTitle("Directions from " + to + " to " + from)
	}

	e.SetDescription("It's right here!")

	mapCtx := sm.NewContext()
	mapCtx.SetSize(600, 600)
	mapCtx.SetTileProvider(sm.NewTileProviderOpenStreetMaps())

	if from == "" {
		mapCtx.SetZoom(17)
	}

	var toCoord []float64
	var fromCoord []float64

	for regex, coords := range LOCATIONS {
		toMatch, err := regexp.MatchString(regex, to)

		if err != nil {
			return err
		}

		fromMatch, err := regexp.MatchString(regex, from)

		if err != nil {
			return err
		}

		if toMatch || fromMatch {
			lat := coords[0]
			long := coords[1]

			if toMatch {
				toCoord = coords
				mapCtx.AddObject(
					sm.NewMarker(
						s2.LatLngFromDegrees(lat, long),
						color.RGBA{R: 0, G: 255, B: 0, A: 255},
						16.0,
					),
				)
			} else if fromMatch {
				fromCoord = coords
				mapCtx.AddObject(
					sm.NewMarker(
						s2.LatLngFromDegrees(lat, long),
						color.RGBA{R: 255, G: 0, B: 0, A: 255},
						16.0,
					),
				)
			}
		}
	}

	if len(toCoord) == 0 {
		errorEmbed := embed.NewErrorEmbed(ctx)

		errorEmbed.SetTitle("Error")
		errorEmbed.SetDescription(fmt.Sprintf("Building `%s` not found.", to))

		_, err := ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Embed:     errorEmbed.MessageEmbed,
			Reference: ctx.GetMessage().Reference(),
		})

		return err
	}

	if len(fromCoord) == 0 && from != "" {
		errorEmbed := embed.NewErrorEmbed(ctx)

		errorEmbed.SetTitle("Error")
		errorEmbed.SetDescription(fmt.Sprintf("Building `%s` not found.", from))

		_, err := ctx.GetSession().ChannelMessageSendComplex(ctx.GetChannel().ID, &discordgo.MessageSend{
			Embed:     errorEmbed.MessageEmbed,
			Reference: ctx.GetMessage().Reference(),
		})

		return err
	}

	if from != "" {
		// get mapbox directions
		dir, err := mapBox.Directions.GetDirections([]base.Location{
			{
				Latitude:  toCoord[0],
				Longitude: toCoord[1],
			},
			{
				Latitude:  fromCoord[0],
				Longitude: fromCoord[1],
			},
		}, directions.RoutingWalking, &directions.RequestOpts{
			Steps:        false,
			Alternatives: false,
		})

		if err != nil {
			return err
		}

		var latLngs []s2.LatLng

		// add route to map
		for _, coord := range dir.Routes[0].Geometry.Coordinates {
			latLngs = append(latLngs, s2.LatLngFromDegrees(coord[1], coord[0]))
		}

		mapCtx.AddObject(
			sm.NewPath(
				latLngs,
				color.RGBA{R: 0, G: 0, B: 255, A: 255},
				3.0,
			),
		)

		// set the zoom accordingly to the route length
		if dir.Routes[0].Distance > 1000 {
			mapCtx.SetZoom(15)
		}
	}

	img, err := mapCtx.Render()

	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)

	// convert to base64
	err = png.Encode(buf, img)

	if err != nil {
		return err
	}

	e.SetImage("attachment://map.png")

	if err != nil {
		return err
	}

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

	return err
}

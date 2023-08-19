package fun

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	"github.com/zekroTJA/shireikan"
)

type Xkcd struct {
	Month      string `json:"month"`
	Num        int64  `json:"num"`
	Link       string `json:"link"`
	Year       string `json:"year"`
	News       string `json:"news"`
	SafeTitle  string `json:"safe_title"`
	Transcript string `json:"transcript"`
	Alt        string `json:"alt"`
	Img        string `json:"img"`
	Title      string `json:"title"`
	Day        string `json:"day"`
}

type XKCD struct {
}

func (c *XKCD) GetInvokes() []string {
	return []string{"xkcd"}
}

func (c *XKCD) GetDescription() string {
	return "Get a random, specific, or the latest XKCD comic"
}

func (c *XKCD) GetHelp() string {
	return "`xkcd <number|latest|random>` - Get a random (default), specific, or the latest XKCD comic"
}

func (c *XKCD) GetGroup() string {
	return config.GroupFun
}

func (c *XKCD) GetDomainName() string {
	return "hamp.fun.xkcd"
}

func (c *XKCD) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}
func (c *XKCD) IsExecutableInDMChannels() bool {
	return true
}

func (c *XKCD) Exec(ctx shireikan.Context) error {
	id := ctx.GetArgs().Get(0).AsString()

	if id == "" {
		id = "random"
	}

	var comic Xkcd

	if id == "latest" {
		res, err := http.Get("https://xkcd.com/info.0.json")

		if err != nil {
			ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while fetching the latest comic.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}

		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&comic)

		if err != nil {
			ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while fetching the latest comic.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}
	} else if id == "random" {
		latest, err := http.Get("https://xkcd.com/info.0.json")

		if err != nil {
			ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while fetching a random comic.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}

		defer latest.Body.Close()

		var tempComic Xkcd

		err = json.NewDecoder(latest.Body).Decode(&tempComic)

		if err != nil {
			ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while fetching a random comic.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}

		id = fmt.Sprintf("%d", rand.Intn((int(tempComic.Num)-1)+1))

		res, err := http.Get(fmt.Sprintf("https://xkcd.com/%s/info.0.json", id))

		if err != nil {
			ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while fetching a random comic.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}

		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&comic)

		if err != nil {
			ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while fetching a random comic.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}
	} else {
		res, err := http.Get(fmt.Sprintf("https://xkcd.com/%s/info.0.json", id))

		if err != nil {
			ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while fetching the comic.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}

		defer res.Body.Close()

		err = json.NewDecoder(res.Body).Decode(&comic)

		if err != nil {
			ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, embed.NewErrorEmbed(ctx).SetTitle("Error").SetDescription("An error occured while fetching the comic.").AddField("Error", err.Error(), false).MessageEmbed)
			return err
		}
	}

	ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, embed.NewEmbed().
		SetTitle(fmt.Sprintf("%s (%d)", comic.Title, comic.Num)).
		SetDescription(comic.Alt).
		SetImage(comic.Img).
		SetColor(0x96A8C8).
		SetURL(fmt.Sprintf("https://xkcd.com/%d", comic.Num)).
		MessageEmbed)
	return nil
}

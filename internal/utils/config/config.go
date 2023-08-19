package config

import (
	"fmt"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	Version string
	Build   string
)

var (
	// BotPrefix is the prefix used for bot commands.
	BotPrefix = ">"

	// BotGuild is the ID of the guild the bot is running on.
	BotGuild = "936651575684915201"

	// HampAPI is the URL to the hamp API.
	// HampAPI = "https://api.hamp.sh"
	HampAPI = "http://localhost:1323"
)

var (
	GroupUtil = "Util"
	GroupFun  = "Fun"
	GroupInfo = "Info"
	GroupMod  = "Moderation"
	GroupDev  = "Dev"
)

var Statuses = []discordgo.Activity{
	{
		Name: "with frogs",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: fmt.Sprintf("PVTA Simulator %d", time.Now().Year()),
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "Stray",
		Type: discordgo.ActivityTypeGame,
	},
	{
		Name: "the Dakin fire alarm",
		Type: discordgo.ActivityTypeListening,
	},
	{
		Name: "at the Roos-Rhode house",
		Type: discordgo.ActivityTypeGame,
	},
}

func ConvertTimestampToDiscordTimestamp(t time.Time) string {
	// format: <t:1234567890> where 1234567890 is the unix timestamp

	u := t.Unix()

	return "<t:" + fmt.Sprint(u) + ">"
}

func ConvertTimestampToDiscordTimestampWithFormat(t time.Time, format string) string {
	// format: <t:1234567890:R> where 1234567890 is the unix timestamp and R is the format

	u := t.Unix()

	return "<t:" + fmt.Sprint(u) + ":" + format + ">"
}

type Hex string

type RGB struct {
	Red   uint8
	Green uint8
	Blue  uint8
}

func (h Hex) toRGB() (RGB, error) {
	return Hex2RGB(h)
}

func Hex2RGB(hex Hex) (RGB, error) {
	var rgb RGB
	values, err := strconv.ParseUint(string(hex), 16, 32)

	if err != nil {
		return RGB{}, err
	}

	rgb = RGB{
		Red:   uint8(values >> 16),
		Green: uint8((values >> 8) & 0xFF),
		Blue:  uint8(values & 0xFF),
	}

	return rgb, nil
}

package config

import (
	"fmt"
	"time"
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

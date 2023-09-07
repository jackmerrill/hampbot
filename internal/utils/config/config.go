package config

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	geojson "github.com/paulmach/go.geojson"
)

var (
	Version string
	Build   string
)

var (
	// BotPrefix is the prefix used for bot commands.
	BotPrefix = "!"

	// BotGuild is the ID of the guild the bot is running on.
	BotGuild = "936651575684915201"

	// HampAPI is the URL to the hamp API.
	HampAPI     = "https://api.hamp.sh"
	HampAPIHost = "api.hamp.sh"
	// HampAPIHost = "localhost:1323"
	// HampAPI = "http://localhost:1323"
)

var (
	GroupUtil        = "Util"
	GroupFun         = "Fun"
	GroupInfo        = "Info"
	GroupMod         = "Moderation"
	GroupStudentLife = "Student Life"
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

func long2tile(lon float64, zoom int) int {
	return int(math.Floor((lon + 180.0) / 360.0 * math.Pow(2.0, float64(zoom))))
}

func lat2tile(lat float64, zoom int) int {
	return int(math.Floor((1.0 - math.Log(math.Tan(lat*math.Pi/180.0)+1.0/math.Cos(lat*math.Pi/180.0))/math.Pi) / 2.0 * math.Pow(2.0, float64(zoom))))
}

func GetTileNumbers(lat, long float64, zoom int) (uint64, uint64) {
	return uint64(long2tile(long, zoom)), uint64(lat2tile(lat, zoom))
}

func decodeToken(bytes []byte) (pos int, value float64) {
	var token int64 = 0
	var shift uint = 0
	var result float64
	var factor float64 = 1e5

	for i, v := range bytes {
		current := int64(v) - 63
		token |= (current & 0x1f) << uint(shift)
		shift += 5

		if current&0x20 == 0 {
			pos = i + 1

			if token&1 != 0 {
				result = float64(^(token >> 1))
			} else {
				result = float64(token >> 1)
			}

			value = result / factor
			return
		}
	}

	pos = 0
	return
}

func DecodePolyline(bytes []byte) []byte {

	fc := geojson.NewFeatureCollection()
	coords := make([][]float64, 0)
	var pos int = 0
	var lat, lng float64

	for pos < len(bytes) {
		current, current_lat := decodeToken(bytes[pos:len(bytes)])
		pos += current
		lat += current_lat

		current, current_lng := decodeToken(bytes[pos:len(bytes)])
		pos += current
		lng += current_lng

		coord := []float64{lng, lat}
		coords = append(coords, coord)
	}

	line := geojson.NewLineStringFeature(coords)

	fc.AddFeature(line)

	json, _ := fc.MarshalJSON()

	return json
}

package fun

import (
	"fmt"
	"time"

	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	"github.com/zekroTJA/shireikan"
)

type MetricTime struct {
}

func (c *MetricTime) GetInvokes() []string {
	return []string{"metrictime", "mt"}
}

func (c *MetricTime) GetDescription() string {
	return "may lord have mercy on us all"
}

func (c *MetricTime) GetHelp() string {
	return "`metrictime` - help"
}

func (c *MetricTime) GetGroup() string {
	return config.GroupFun
}

func (c *MetricTime) GetDomainName() string {
	return "hamp.util.metrictime"
}

func (c *MetricTime) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}
func (c *MetricTime) IsExecutableInDMChannels() bool {
	return true
}

func (c *MetricTime) Exec(ctx shireikan.Context) error {
	now := time.Now()                                                               // Get current time for UTC
	midnight := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC) // Get midnight time for UTC
	secSinceMidnight := now.Sub(midnight).Seconds()                                 // Time duration since midnight in regular seconds
	metricSecOfDay := int(secSinceMidnight / 0.864)                                 // No. of seconds in the week in metric world

	metricHour := metricSecOfDay / 10000        // Each metric hour has 10,000 metric seconds
	metricMin := (metricSecOfDay % 10000) / 100 // Each metric minutes has 100 metric seconds
	metricSec := metricSecOfDay % 100           // Remaining are metric seconds

	// Format the time
	timeStr := fmt.Sprintf("%02d:%02d:%02d", metricHour, metricMin, metricSec)

	// Send the message
	e := embed.NewSuccessEmbed(ctx)

	e.SetTitle("Metric Time")
	e.SetDescription(timeStr)
	e.SetColor(0x00ff00)

	_, err := ctx.GetSession().ChannelMessageSendEmbed(ctx.GetChannel().ID, e.MessageEmbed)
	return err
}

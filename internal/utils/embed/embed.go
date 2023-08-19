package embed

import (
	"fmt"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Embed struct {
	*discordgo.MessageEmbed
}

func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{}}
}

func (e *Embed) SetTitle(title string) *Embed {
	e.MessageEmbed.Title = title
	return e
}

func (e *Embed) SetDescription(description string) *Embed {
	e.MessageEmbed.Description = description
	return e
}

func (e *Embed) SetColor(color int) *Embed {
	e.MessageEmbed.Color = color
	return e
}

func (e *Embed) SetColorRGB(r, g, b int) *Embed {
	return e.SetColor((r << 16) + (g << 8) + b)
}

func (e *Embed) SetColorHex(hex string) *Embed {
	var color int
	fmt.Sscanf(hex, "#%06x", &color)
	return e.SetColor(color)
}

func (e *Embed) SetURL(url string) *Embed {
	e.MessageEmbed.URL = url
	return e
}

func (e *Embed) SetTimestamp(timestamp time.Time) *Embed {
	e.MessageEmbed.Timestamp = timestamp.Format(time.RFC3339)
	return e
}

func (e *Embed) SetFooter(text, iconURL string) *Embed {
	e.MessageEmbed.Footer = &discordgo.MessageEmbedFooter{
		Text:    text,
		IconURL: iconURL,
	}
	return e
}

func (e *Embed) SetImage(url string) *Embed {
	e.MessageEmbed.Image = &discordgo.MessageEmbedImage{
		URL: url,
	}
	return e
}

func (e *Embed) SetThumbnail(url string) *Embed {
	e.MessageEmbed.Thumbnail = &discordgo.MessageEmbedThumbnail{
		URL: url,
	}
	return e
}

func (e *Embed) SetAuthor(name, url, iconURL string) *Embed {
	e.MessageEmbed.Author = &discordgo.MessageEmbedAuthor{
		Name:    name,
		URL:     url,
		IconURL: iconURL,
	}
	return e
}

func (e *Embed) AddField(name, value string, inline bool) *Embed {
	e.MessageEmbed.Fields = append(e.MessageEmbed.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
	return e
}

func (e *Embed) AddFields(fields ...*discordgo.MessageEmbedField) *Embed {
	e.MessageEmbed.Fields = append(e.MessageEmbed.Fields, fields...)
	return e
}

func (e *Embed) AddFieldsFromMap(fields map[string]string, inline bool) *Embed {
	for name, value := range fields {
		e.AddField(name, value, inline)
	}
	return e
}

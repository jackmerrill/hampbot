package fun

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	"github.com/sashabaranov/go-openai"
	"github.com/zekroTJA/shireikan"
)

type Dalle struct {
}

func (c *Dalle) GetInvokes() []string {
	return []string{"dalle", "ai-art"}
}

func (c *Dalle) GetDescription() string {
	return "Generate AI art with DALL-E"
}

func (c *Dalle) GetHelp() string {
	return "`dalle [prompt]` - Generate AI art with DALL-E"
}

func (c *Dalle) GetGroup() string {
	return config.GroupFun
}

func (c *Dalle) GetDomainName() string {
	return "hamp.fun.dalle"
}

func (c *Dalle) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}
func (c *Dalle) IsExecutableInDMChannels() bool {
	return true
}

func (c *Dalle) Exec(ctx shireikan.Context) error {
	openaiToken := os.Getenv("OPENAI_TOKEN")

	if openaiToken == "" {
		ctx.ReplyEmbed(embed.NewErrorEmbed(ctx).SetDescription("No OpenAI token set.").MessageEmbed)
		return fmt.Errorf("no openai token set")
	}

	prompt := strings.TrimPrefix(ctx.GetMessage().Content, fmt.Sprintf("%sdalle ", config.BotPrefix))

	client := openai.NewClient(openaiToken)

	ctx.GetSession().ChannelTyping(ctx.GetChannel().ID)

	resp, err := client.CreateImage(context.Background(), openai.ImageRequest{
		Prompt:         prompt,
		Size:           openai.CreateImageSize1024x1024,
		ResponseFormat: openai.CreateImageResponseFormatURL,
		N:              1,
	})

	if err != nil {
		ctx.ReplyEmbed(embed.NewErrorEmbed(ctx).SetDescription("Failed to generate image.").MessageEmbed)
		return err
	}

	ctx.ReplyEmbed(embed.NewSuccessEmbed(ctx).
		SetTitle("DALL-E").
		SetDescription(prompt).
		SetImage(resp.Data[0].URL).
		MessageEmbed)
	return nil
}

package fun

import (
	"context"
	"fmt"
	"os"

	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/jackmerrill/hampbot/internal/utils/embed"
	"github.com/sashabaranov/go-openai"
	"github.com/zekroTJA/shireikan"
)

type AI struct {
}

func (c *AI) GetInvokes() []string {
	return []string{"ai", "gpt", "ask", "question"}
}

func (c *AI) GetDescription() string {
	return "Ask AI anything!"
}

func (c *AI) GetHelp() string {
	return "`ai [prompt]` - Ask GPT-3.5 anything."
}

func (c *AI) GetGroup() string {
	return config.GroupFun
}

func (c *AI) GetDomainName() string {
	return "hamp.fun.AI"
}

func (c *AI) GetSubPermissionRules() []shireikan.SubPermission {
	return nil
}
func (c *AI) IsExecutableInDMChannels() bool {
	return true
}

func (c *AI) Exec(ctx shireikan.Context) error {
	openaiToken := os.Getenv("OPENAI_TOKEN")

	if openaiToken == "" {
		ctx.ReplyEmbed(embed.NewErrorEmbed(ctx).SetDescription("No OpenAI token set.").MessageEmbed)
		return fmt.Errorf("no openai token set")
	}

	client := openai.NewClient(openaiToken)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are HampBot. You are a bot on Discord that can answer any question. Multiple users may ask you questions at the same time. There will be context given. Don't add any decorations (such as `AI:` or `User:`) to your response, these will be added automatically.",
		},
	}

	if ctx.GetChannel().IsThread() {
		// get all messages in thread
		msgs, err := ctx.GetSession().ChannelMessages(ctx.GetChannel().ID, 100, "", "", "")
		if err != nil {
			return err
		}

		for _, msg := range msgs {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("%s: %s", msg.Author.Username, msg.Content),
			})
		}
	} else {
		msg := ctx.GetMessage()

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: fmt.Sprintf("%s: %s", msg.Author.Username, msg.Content),
		})
	}

	ctx.GetSession().ChannelTyping(ctx.GetChannel().ID)

	// get response from GPT-3.5
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest{
		Model:    openai.GPT3Dot5Turbo,
		Messages: messages,
	})

	if err != nil {
		ctx.ReplyEmbed(embed.NewErrorEmbed(ctx).SetDescription("An error occured while asking GPT-3.5.").AddField("Error", err.Error(), false).MessageEmbed)
		return err
	}

	// send response
	ctx.GetSession().ChannelMessageSend(ctx.GetChannel().ID, fmt.Sprintf(":robot: **AI:** %s", resp.Choices[0].Message.Content))

	return nil
}

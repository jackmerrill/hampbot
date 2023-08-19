package embed

import (
	"fmt"

	"github.com/jackmerrill/hampbot/internal/utils/config"
	"github.com/zekroTJA/shireikan"
)

func NewGenericEmbed(ctx shireikan.Context) *Embed {
	return NewEmbed().SetFooter(fmt.Sprintf("HampBot %s", config.Version), ctx.GetSession().State.User.AvatarURL("256"))
}

func NewSuccessEmbed(ctx shireikan.Context) *Embed {
	return NewGenericEmbed(ctx).SetColor(0x00ff00)
}

func NewWarningEmbed(ctx shireikan.Context) *Embed {
	return NewGenericEmbed(ctx).SetColor(0xffff00)
}

func NewErrorEmbed(ctx shireikan.Context) *Embed {
	return NewGenericEmbed(ctx).SetColor(0xff0000).SetTitle("Error")
}

package bot

import (
	"github.com/spf13/viper"
	tb "gopkg.in/telebot.v3"
	"strconv"
)

func isManageMiddleware(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		if isManage(c.Chat(), c.Sender().ID) {
			return next(c)
		}
		return c.Respond(&tb.CallbackResponse{
			Text:      "💢点什么，你又不是管理员",
			ShowAlert: true,
		})
	}
}

func PreGroupMiddleware(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		if !c.Message().FromGroup() {
			return nil
		}
		groups := viper.GetStringSlice("telegram.groups")
		for _, group := range groups {
			id, err := strconv.ParseInt(group, 10, 64)
			if err != nil {
				return err
			}
			if c.Chat().ID == id {
				return next(c)
			}
		}
		return nil
	}
}

func PreCmdMiddleware(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		if c.Message().FromGroup() && !isManage(c.Chat(), c.Sender().ID) {
			c.Delete()
		}
		return next(c)
	}
}

func CreatorCmdMiddleware(next tb.HandlerFunc) tb.HandlerFunc {
	return func(c tb.Context) error {
		if !isOwner(c.Sender().ID) {
			return nil
		}
		return next(c)
	}
}

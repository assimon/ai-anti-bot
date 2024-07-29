package bot

import (
	"fmt"
	"github.com/spf13/viper"
	tb "gopkg.in/telebot.v3"
	"time"
)

var (
	Bot *tb.Bot
)

func Start() error {
	var err error
	setting := tb.Settings{
		Token:   viper.GetString("telegram.token"),
		Updates: 100,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second, AllowedUpdates: []string{
			"message",
			"chat_member",
			"inline_query",
			"callback_query",
		}},
		OnError: func(err error, context tb.Context) {
			fmt.Printf("%+v\n", err)
		},
	}
	if viper.GetString("telegram.proxy") != "" {
		setting.URL = viper.GetString("telegram.proxy")
	}
	Bot, err = tb.NewBot(setting)
	if err != nil {
		return err
	}
	RegisterCommands()
	RegisterHandle()
	Bot.Start()
	return nil
}

func RegisterCommands() {
	_ = Bot.SetCommands(tb.Command{
		Text:        StartCmd,
		Description: "Hello🙌",
	})
}

func RegisterHandle() {
	Bot.Handle(StartCmd, func(c tb.Context) error {
		return c.Send("🙋hi,I am an AI anti-advertising robot. My father is Assimon. github.com/assimon")
	}, PreCmdMiddleware)

	Bot.Handle(AllAdCmd, AllAd, CreatorCmdMiddleware)
	Bot.Handle(AddAdCmd, AddAd, CreatorCmdMiddleware)
	Bot.Handle(DelAdCmd, DelAd, CreatorCmdMiddleware)

	Bot.Handle(tb.OnText, OnTextMessage, PreGroupMiddleware)
	Bot.Handle(tb.OnSticker, OnStickerMessage, PreGroupMiddleware)
	Bot.Handle(tb.OnPhoto, OnPhotoMessage, PreGroupMiddleware)
	Bot.Handle(tb.OnUserJoined, OnOnUserJoinedMessage, PreGroupMiddleware)
}

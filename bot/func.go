package bot

import (
	"fmt"
	"github.com/assimon/ai-anti-bot/adapter"
	"github.com/assimon/ai-anti-bot/ai"
	"github.com/assimon/ai-anti-bot/database"
	"github.com/spf13/viper"
	tb "gopkg.in/telebot.v3"
	"log"
	"strconv"
)

func Identifier() adapter.IModel {
	im := viper.GetString("identification_model")
	var option adapter.Option
	err := viper.UnmarshalKey(im, &option)
	if err != nil {
		log.Fatal(err)
	}
	model, err := ai.New(im, option)
	if err != nil {
		log.Fatal(err)
	}
	return model
}

func isManage(chat *tb.Chat, userId int64) bool {
	adminList, err := Bot.AdminsOf(chat)
	if err != nil {
		return false
	}
	for _, member := range adminList {
		if member.User.ID == userId {
			return true
		}
	}
	return false
}

func isOwner(userId int64) bool {
	owners := viper.GetStringSlice("telegram.owners")
	for _, owner := range owners {
		id, err := strconv.ParseInt(owner, 10, 64)
		if err != nil {
			return false
		}
		if id == userId {
			return true
		}
	}
	return false
}

func BanChatMember(c tb.Context, res *adapter.RecognizeResult) (err error) {
	userLink := fmt.Sprintf("tg://user?id=%d", c.Message().Sender.ID)
	userNickname := c.Message().Sender.LastName + c.Message().Sender.FirstName
	userId := c.Message().Sender.ID
	blockMessage := fmt.Sprintf(
		viper.GetString("message.block_hint"),
		userNickname,
		userLink,
		res.SpamScore,
		res.SpamMockText,
	)
	// ban user
	err = Bot.Restrict(c.Chat(), &tb.ChatMember{
		Rights:          tb.NoRights(),
		User:            c.Message().Sender,
		RestrictedUntil: tb.Forever(),
	})
	if err != nil {
		return err
	}
	manslaughterBtn := manslaughterMenu.Data("👮🏻Unblock", strconv.FormatInt(userId, 10))
	manslaughterMenu.Inline(manslaughterMenu.Row(manslaughterBtn))
	LoadAdMenuBtn(manslaughterMenu)
	Bot.Handle(&manslaughterBtn, func(c tb.Context) error {
		if err = Bot.Delete(c.Message()); err != nil {
			return err
		}
		err = Bot.Restrict(c.Chat(), &tb.ChatMember{
			User:   &tb.User{ID: userId},
			Rights: tb.NoRestrictions(),
		})
		if err != nil {
			return err
		}
		return c.Send(fmt.Sprintf("The administrator has unbanned user: [%s](%s)", userNickname, userLink), tb.ModeMarkdownV2)
	}, isManageMiddleware)
	if err = c.Reply(blockMessage, manslaughterMenu, tb.ModeMarkdownV2); err != nil {
		return err
	}
	return
}

func LoadAdMenuBtn(menu *tb.ReplyMarkup) {
	advertises, err := database.GetEfficientAdvertise()
	if err != nil {
		log.Println(err)
	} else {
		for _, advertise := range advertises {
			menu.InlineKeyboard = append(menu.InlineKeyboard, []tb.InlineButton{
				{
					Text: advertise.Title,
					URL:  advertise.Url,
				},
			})
		}
	}
}

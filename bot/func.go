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
	"time"
	"unicode/utf8"
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
		ProcessNickname(userNickname),
		userLink,
		res.SpamScore,
		res.SpamReason,
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
	msg, err := Bot.Send(c.Chat(), blockMessage, manslaughterMenu, tb.ModeMarkdownV2)
	if err != nil {
		return err
	}
	time.AfterFunc(time.Minute, func() {
		err = Bot.Delete(msg)
		if err != nil {
			log.Println(err)
		}
	})
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

func ProcessNickname(nickname string) string {
	length := utf8.RuneCountInString(nickname)
	switch length {
	case 1:
		return "||" + nickname + "||"
	case 2:
		runes := []rune(nickname)
		return string(runes[0]) + "||" + string(runes[1]) + "||"
	default:
		runes := []rune(nickname)
		if length > 2 {
			return string(runes[0]) + "||" + string(runes[1:length-1]) + "||" + string(runes[length-1])
		}
		return nickname
	}
}

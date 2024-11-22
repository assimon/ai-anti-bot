package bot

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/assimon/ai-anti-bot/database"
	"github.com/golang-module/carbon/v2"
	"github.com/spf13/viper"
	tb "gopkg.in/telebot.v3"
	"io"
	"log"
	"strconv"
	"strings"
	"time"
)

var (
	manslaughterMenu = &tb.ReplyMarkup{ResizeKeyboard: true}
)

func PreCheck(c tb.Context) (user *database.UserInfo, needCheck bool, err error) {
	if isManage(c.Chat(), c.Sender().ID) {
		return
	}
	first := database.UserInfo{
		TelegramUserId: c.Sender().ID,
		TelegramChatId: c.Chat().ID,
	}
	user, err = database.GetUserInfo(&first)
	if err != nil {
		return
	}
	if user.VerificationTimes > viper.GetInt64("strategy.verification_times") {
		return
	}
	diffDay := carbon.Now().DiffAbsInDays(carbon.Parse(user.JoinedTime.ToDateTimeString()))
	if diffDay > viper.GetInt64("strategy.joined_time") || user.NumberOfSpeeches > viper.GetInt64("strategy.number_of_speeches") {
		return
	}
	return user, true, nil
}

func OnTextMessage(c tb.Context) error {
	user, needCheck, err := PreCheck(c)
	if err != nil {
		return err
	}
	if c.Sender().IsBot && viper.GetBool("clean_bot_message") {
		time.AfterFunc(time.Minute, func() {
			err = c.Delete()
			if err != nil {
				log.Println(err)
			}
		})
	}
	if user != nil {
		defer func() {
			_ = database.IncrementNumberOfSpeeches(user)
		}()
	}
	checkFun := func() error {
		userInfo := fmt.Sprintf(
			viper.GetString("prompt.user_info"),
			c.Sender().LastName,
			c.Sender().FirstName,
			user.NumberOfSpeeches+1,
			user.JoinedTime.ToDateTimeString(),
		)
		res, err := Identifier().RecognizeTextMessage(context.Background(), userInfo, c.Message().Text)
		if err != nil {
			return err
		}
		if res.State == 0 {
			_ = database.IncrementVerificationTimes(user)
			return nil
		}
		if err = BanChatMember(c, &res); err != nil {
			return err
		}
		return c.Delete()
	}
	if needCheck {
		go func() {
			err = checkFun()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	return nil
}

func OnStickerMessage(c tb.Context) error {
	user, needCheck, err := PreCheck(c)
	if err != nil {
		return err
	}
	if user != nil {
		defer func() {
			_ = database.IncrementNumberOfSpeeches(user)
		}()
	}
	checkFun := func() error {
		rc, err := Bot.File(&c.Message().Sticker.File)
		defer rc.Close()
		if err != nil {
			return err
		}
		buf := &bytes.Buffer{}
		if _, err = io.Copy(buf, rc); err != nil {
			return err
		}
		baseEncStr := MediaTypeWebp + base64.StdEncoding.EncodeToString(buf.Bytes())
		userInfo := fmt.Sprintf(
			viper.GetString("prompt.user_info"),
			c.Sender().LastName,
			c.Sender().FirstName,
			user.NumberOfSpeeches,
			user.JoinedTime.ToDateTimeString(),
		)
		res, err := Identifier().RecognizeImageMessage(context.Background(), userInfo, baseEncStr)
		if res.State == 0 {
			_ = database.IncrementVerificationTimes(user)
			return nil
		}
		if err = BanChatMember(c, &res); err != nil {
			return err
		}
		return c.Delete()
	}
	if needCheck {
		go func() {
			err = checkFun()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	return nil
}

func OnPhotoMessage(c tb.Context) error {
	user, needCheck, err := PreCheck(c)
	if err != nil {
		return err
	}
	if user != nil {
		defer func() {
			_ = database.IncrementNumberOfSpeeches(user)
		}()
	}
	checkFun := func() error {
		rc, err := Bot.File(&c.Message().Photo.File)
		defer rc.Close()
		if err != nil {
			return err
		}
		buf := &bytes.Buffer{}
		if _, err = io.Copy(buf, rc); err != nil {
			return err
		}
		baseEncStr := MediaTypeJpg + base64.StdEncoding.EncodeToString(buf.Bytes())
		userInfo := fmt.Sprintf(
			viper.GetString("prompt.user_info"),
			c.Sender().LastName,
			c.Sender().FirstName,
			user.NumberOfSpeeches,
			user.JoinedTime.ToDateTimeString(),
		)
		res, err := Identifier().RecognizeImageMessage(context.Background(), userInfo, baseEncStr)
		if res.State == 0 {
			_ = database.IncrementVerificationTimes(user)
			return nil
		}
		if err = BanChatMember(c, &res); err != nil {
			return err
		}
		return c.Delete()
	}
	if needCheck {
		go func() {
			err = checkFun()
			if err != nil {
				fmt.Println(err)
			}
		}()
	}
	return nil
}

func OnChatMemberMessage(c tb.Context) error {
	user := database.UserInfo{
		TelegramUserId:    c.Sender().ID,
		TelegramChatId:    c.Chat().ID,
		JoinedTime:        carbon.Now().ToDateTimeStruct(),
		NumberOfSpeeches:  0,
		VerificationTimes: 0,
	}
	err := database.SaveUserInfo(&user)
	return err
}

func AddAd(c tb.Context) error {
	payload := c.Message().Payload
	payloadSlice := strings.Split(payload, "|")
	if len(payloadSlice) != 4 {
		return c.Send("❌格式错误")
	}
	title := payloadSlice[0]
	url := payloadSlice[1]
	validityPeriod := payloadSlice[2]
	sort, _ := strconv.Atoi(payloadSlice[3])
	ad := database.Advertise{
		Title:          title,
		Url:            url,
		Sort:           sort,
		ValidityPeriod: carbon.Parse(validityPeriod).ToDateTimeStruct(),
	}
	err := database.AddAdvertise(ad)
	if err != nil {
		return c.Send("❌无法添加推广:" + err.Error())
	}
	if err = c.Send("✅成功添加推广"); err != nil {
		fmt.Println("[AddAd] send success message err:", err)
	}
	return AllAd(c)
}

func AllAd(c tb.Context) error {
	adList, err := database.AllAdvertise()
	if err != nil {
		return c.Send("❌无法获取推广:" + err.Error())
	}
	table := "💾全部推广：\n"
	for _, advertise := range adList {
		table += fmt.Sprintf("Id:%d\n推广名:%s 链接:%s 排序:%d 到期时间:%s 创建时间:%s \n",
			advertise.ID,
			advertise.Title,
			advertise.Url,
			advertise.Sort,
			advertise.ValidityPeriod.ToDateTimeString(),
			advertise.CreatedAt.ToDateTimeString(),
		)
	}
	return c.Send(table)
}

func DelAd(c tb.Context) error {
	payload := c.Message().Payload
	if payload == "" {
		return c.Send("❌格式错误")
	}
	id, err := strconv.ParseInt(payload, 10, 64)
	if err != nil {
		return c.Send(err.Error())
	}
	if err = database.DeleteAdvertise(id); err != nil {
		return c.Send(err.Error())
	}
	if err = c.Send("✅成功删除推广"); err != nil {
		fmt.Println("[DelAd] send success message err:", err)
	}
	return AllAd(c)
}

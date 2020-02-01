package menu

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"libretaxi/context"
	"libretaxi/objects"
	"libretaxi/validation"
	"log"
	"strings"
)

type PostMenuHandler struct {
	user *objects.User
	context *context.Context
}

func (handler *PostMenuHandler) informUsersAround(lon float64, lat float64, text string, postId int64) {
	userIds := handler.context.Repo.UserIdsAround(lon, lat)

	textWithContacts := ""

	if len(handler.user.Username) == 0 {
		userTextContact := fmt.Sprintf("[%s %s](tg://user?id=%d)", handler.user.FirstName, handler.user.LastName, handler.user.UserId)
		textWithContacts = fmt.Sprintf("%s\n\nvia %s", text, userTextContact)
	} else {
		textWithContacts = fmt.Sprintf("%s\n\nvia @%s", text, handler.user.Username)
	}

	for i, _ := range userIds {
		userId := userIds[i]
		msg := tgbotapi.NewMessage(userId, textWithContacts)
		msg.ParseMode = "MarkdownV2"

		reportKeyboard := tgbotapi.NewInlineKeyboardMarkup(
			tgbotapi.NewInlineKeyboardRow(
				tgbotapi.NewInlineKeyboardButtonData("☝️️Report ⚠️",fmt.Sprintf("{\"Action\":\"REPORT_POST\",\"Id\":%d}", postId)),
			),
		)
		msg.ReplyMarkup = reportKeyboard
		_, err := handler.context.Bot.Send(msg)

		if err != nil {
			log.Println(err)
		}
	}
}

func (handler *PostMenuHandler) Handle(user *objects.User, context *context.Context, message *tgbotapi.Message) {
	log.Println("Post menu")

	handler.user = user
	handler.context = context

	if context.Repo.UserPostedRecently(user.UserId) {

		msg := tgbotapi.NewMessage(user.UserId, "🕙 Wait for 5 minutes")
		context.Bot.Send(msg)

		user.MenuId = objects.Menu_Feed
		context.Repo.SaveUser(user)

	} else if len(message.Text) == 0 {

		msg := tgbotapi.NewMessage(user.UserId, "Copy & paste text starting with 🚗 or 👋 in the following format (you can use your own language), or /cancel, examples:")
		context.Bot.Send(msg)

		msg = tgbotapi.NewMessage(user.UserId, `🚗 Driver looking for passenger(s)
Pick up: foobar square
Drop off: airport
Date: today
Time: now
Payment: cash, venmo`)
		context.Bot.Send(msg)

		msg = tgbotapi.NewMessage(user.UserId, `👋🏻 Passenger looking for driver
Pick up: foobar st, 42
Drop off: downtown
Date: today
Time: now
Pax: 1`)
		context.Bot.Send(msg)

	} else {

		textValidation := validation.NewTextValidation()
		error := textValidation.Validate(message.Text)

		if len(error) > 0 {
			msg := tgbotapi.NewMessage(user.UserId, error)
			context.Bot.Send(msg)
			return
		}

		cleanText := strings.TrimSpace(message.Text)

		post := &objects.Post{
			UserId: user.UserId,
			Text: cleanText,
			Lon: user.Lon,
			Lat: user.Lat,
			ReportCnt: 0,
		}

		context.Repo.SavePost(post);

		handler.informUsersAround(post.Lon, post.Lat, cleanText, post.PostId)

		msg := tgbotapi.NewMessage(user.UserId, "✅ Sent to users around you (25km)")
		context.Bot.Send(msg)

		user.MenuId = objects.Menu_Feed
		context.Repo.SaveUser(user)
	}
}

func NewPostMenu() *PostMenuHandler {
	handler := &PostMenuHandler{}
	return handler
}

package controller

import (
	"fmt"
	"net/http"
	"os"

	"github.com/jgvkmea/go-sort-tabelog/interface/controller/middleware/logger"
	"github.com/jgvkmea/go-sort-tabelog/interface/gateway"
	"github.com/jgvkmea/go-sort-tabelog/usecase"
	"github.com/line/line-bot-sdk-go/linebot"
)

func TabelogSearchHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.FromContext(ctx)

	lineClient, err := linebot.New(getChannelSecret(), getChannelToken())
	if err != nil {
		log.Errorln("failed to create linebot: ", err)
		return
	}

	alert := gateway.NewAlertBot()
	events, err := lineClient.ParseRequest(req)
	if err != nil {
		errMessage := fmt.Sprintf("failed to parse request: %s", err)
		log.Errorln(errMessage)
		alert.AlertMessage(errMessage)
		w.WriteHeader(500)
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				u := usecase.NewTabelogInteractor(logger.FromContext)
				inputdata := usecase.NewGetShopsOrderByRatingInputData(
					gateway.NewLineClient(lineClient),
					gateway.NewWebDriver(),
					message.Text,
					event.ReplyToken,
					event.Source.UserID,
				)

				err := u.GetShopsOrderByRating(ctx, inputdata)
				if err != nil {
					eMsg := fmt.Sprintf("failed to usecase.GetShopsOrderByRating")
					log.Errorln(eMsg)

					err = alert.AlertMessage(eMsg)
					if err != nil {
						log.Errorln("failed to alert message: ", err)
					}

					if err == usecase.BadRequestError {
						w.WriteHeader(400)
						return
					}
					w.WriteHeader(500)
					return
				}
			default:
				log.Errorln("received message is not TextMessage")
				w.WriteHeader(400)
			}
		}
	}
	return
}

func getChannelSecret() string {
	return os.Getenv("TABELOG_SORT_CHANNEL_SECRET")
}

func getChannelToken() string {
	return os.Getenv("TABELOG_SORT_CHANNEL_TOKEN")
}

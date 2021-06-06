package service

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/sirupsen/logrus"
)

const replyText = "今から調べるから数分待っててね〜\n(多めに言ってるわけじゃなくてちゃんと数分かかります)"

func TabelogSearchHandler(w http.ResponseWriter, req *http.Request) {
	// TODO: contextから取得するようにかきかえる
	log := logrus.New()

	lineClient, err := linebot.New(
		os.Getenv("CHANNEL_SECRET"),
		os.Getenv("CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Errorln("failed to create linebot: ", err)
		return
	}

	events, err := lineClient.ParseRequest(req)
	if err != nil {
		log.Errorln("failed to parse request: ", err)
		w.WriteHeader(500)
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if _, err := lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do(); err != nil {
					log.Errorln("echo reply failed: ", err)
					return
				}

				conditions := strings.Split(message.Text, " ")
				if len(conditions) > 2 {
					log.Errorln("received search condition is too many")
					w.WriteHeader(400)
				}

				area := conditions[0]
				keyword := ""
				if len(conditions) == 2 {
					keyword = conditions[1]
				}
				log.Infof("are: %s, keyword: %s", area, keyword)

				shops, err := GetShopsOrderByRating(area, keyword)
				if err != nil {
					log.Errorln("failed to get shops order by rating: ", err)
					w.WriteHeader(500)
				}

				count := getOutputCount(shops)
				var pushMessages []string
				for i := 0; i < count; i++ {
					pushMessages = append(pushMessages, fmt.Sprintf(
						"%d位 %s\nRating: %g\n%s",
						i+1,
						shops[i].Name,
						float64(shops[i].Rating)/100,
						shops[i].Url,
					))
				}

				_, err = lineClient.PushMessage(
					event.Source.UserID,
					linebot.NewTextMessage(strings.Join(pushMessages, "\n\n")),
				).Do()
				if err != nil {
					log.Errorln("failed to reply message: ", err)
					w.WriteHeader(500)
				}
			default:
				log.Errorln("received message is not TextMessage")
				w.WriteHeader(400)
			}
		}
	}
	return
}

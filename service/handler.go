package service

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/jgvkmea/go-sort-tabelog/middleware/logger"
	"github.com/jgvkmea/go-sort-tabelog/utils"
	"github.com/line/line-bot-sdk-go/linebot"
)

const ()

var (
	replyText                = fmt.Sprintf("今から調べるから数分待っててね%c", navigateMan)
	littleConditionErrorText = fmt.Sprintf("%c入力エラー%c\n検索条件が足りないよ%c%c\n「エリア名 キーワード」で検索してね！", warning, warning, tiredFace, tiredFace)
	manyConditionErrorText   = fmt.Sprintf("%c入力エラー%c\n検索条件が多すぎるよ%c%c\n「エリア名 キーワード」で検索してね！", warning, warning, tiredFace, tiredFace)
)

const (
	emojiStar   = 0x2B50
	tiredFace   = 0x1F62B
	warning     = 0x1F6AB
	navigateMan = 0x1F481
)

func TabelogSearchHandler(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()
	log := logger.FromContext(ctx)

	lineClient, err := linebot.New(
		os.Getenv("TABELOG_SORT_CHANNEL_SECRET"),
		os.Getenv("TABELOG_SORT_CHANNEL_TOKEN"),
	)
	if err != nil {
		log.Errorln("failed to create linebot: ", err)
		return
	}

	events, err := lineClient.ParseRequest(req)
	if err != nil {
		errMessage := fmt.Sprintf("failed to parse request: %s", err)
		log.Errorln(errMessage)
		utils.AlertByLinebot(errMessage)
		w.WriteHeader(500)
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:

				text := strings.ReplaceAll(message.Text, "　", " ")
				conditions := strings.Split(text, " ")
				if len(conditions) < 2 {
					log.Errorln("received search condition is too little")
					if _, err := lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(littleConditionErrorText)).Do(); err != nil {
						err = fmt.Errorf("echo reply failed: %s", err)
						log.Errorln(err)
						return
					}
					w.WriteHeader(400)
					return
				} else if len(conditions) > 2 {
					log.Errorln("received search condition is too many")
					if _, err := lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(manyConditionErrorText)).Do(); err != nil {
						err = fmt.Errorf("echo reply failed: %s", err)
						log.Errorln(err)
						return
					}
					w.WriteHeader(400)
					return
				}

				if _, err := lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyText)).Do(); err != nil {
					err = fmt.Errorf("echo reply failed: %s", err)
					log.Errorln(err)
					return
				}

				area := conditions[0]
				keyword := ""
				if len(conditions) == 2 {
					keyword = conditions[1]
				}
				log.Infof("are: %s, keyword: %s", area, keyword)

				shops, err := GetShopsOrderByRating(area, keyword)
				if err != nil {
					eMsg := fmt.Sprintf("failed to get shops order by rating: %s", err)
					log.Errorln(eMsg)
					err = utils.AlertByLinebot(eMsg)
					if err != nil {
						log.Errorln("failed to push message by linebot")
					}
					w.WriteHeader(500)
					return
				}

				count := getOutputCount(shops)
				var pushMessages []string
				for i := 0; i < count; i++ {
					pushMessages = append(pushMessages, fmt.Sprintf(
						"%d位 %s\n%c%g\n%s",
						i+1,
						shops[i].Name,
						emojiStar,
						float64(shops[i].Rating)/100,
						shops[i].Url,
					))
				}

				_, err = lineClient.PushMessage(
					event.Source.UserID,
					linebot.NewTextMessage(strings.Join(pushMessages, "\n\n")),
				).Do()
				if err != nil {
					eMsg := fmt.Sprintf("failed to reply message: %s", err)
					log.Errorln(eMsg)
					err = utils.AlertByLinebot(eMsg)
					if err != nil {
						log.Errorln("failed to push message by linebot")
					}
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

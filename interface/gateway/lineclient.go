package gateway

import (
	"fmt"
	"strings"

	"github.com/jgvkmea/go-sort-tabelog/entity"
	"github.com/jgvkmea/go-sort-tabelog/usecase"
	"github.com/line/line-bot-sdk-go/linebot"
)

const (
	EMOJI_STAR = "⭐️"
)

type LineClient struct {
	lc *linebot.Client
}

func NewLineClient(l *linebot.Client) usecase.SNSClient {
	return &LineClient{l}
}

func (lc *LineClient) ReplyMessage(replyToken string, message string) error {
	_, err := lc.lc.ReplyMessage(replyToken, linebot.NewTextMessage(message)).Do()
	if err != nil {
		return fmt.Errorf("failed to reply message: %s", err)
	}
	return nil
}

func (lc *LineClient) PushShopsMessage(destinationUserID string, shops entity.Shops) error {
	message, err := createShopListMessage(shops)
	if err != nil {
		return fmt.Errorf("failed to create shop list message: %s", err)
	}

	_, err = lc.lc.PushMessage(destinationUserID, linebot.NewTextMessage(message)).Do()
	if err != nil {
		return fmt.Errorf("failed to push message: %s", err)
	}
	return nil
}

func createShopListMessage(shops []entity.Shop) (string, error) {
	if len(shops) == 0 {
		return "", fmt.Errorf("shops length is 0")
	}

	var pushMessages []string
	for i, shop := range shops {
		pushMessages = append(pushMessages, fmt.Sprintf(
			"%d位 %s\n%s%g\n%s",
			i+1,
			shop.GetName(),
			EMOJI_STAR,
			float64(shop.GetRating())/100,
			shop.GetURL(),
		))
	}
	return strings.Join(pushMessages, "\n\n"), nil
}

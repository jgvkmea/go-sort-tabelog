package usecase

import (
	"context"

	"github.com/sirupsen/logrus"
)

type TabelogInteractor struct {
	FromContext func(ctx context.Context) *logrus.Logger
	snsClient   SNSClient
}

func NewTabelogInteractor(FromContext func(ctx context.Context) *logrus.Logger, snsClient SNSClient) IUseCase {
	return TabelogInteractor{FromContext, snsClient}
}

func (t TabelogInteractor) GetShopsOrderByRating(ctx context.Context, inputdata *GetShopsOrderByRatingInputData) error {
	log := t.FromContext(ctx)
	log.Infoln("start to GetShopsOrderByRating")

	conditions := getSearchCondition(inputdata.message)
	if len(conditions) < 1 {
		log.Errorln("received search condition is too little")
		if err := t.snsClient.ReplyMessage(inputdata.replyToken, littleConditionErrorText); err != nil {
			log.Errorf("failed to reply message: %s", err)
			return InternalError
		}
		return BadRequestError
	} else if len(conditions) > 2 {
		log.Errorln("received search condition is too many")
		if err := t.snsClient.ReplyMessage(inputdata.replyToken, manyConditionErrorText); err != nil {
			log.Errorf("failed to reply message: %s", err)
			return InternalError
		}
		return BadRequestError
	}

	if err := t.snsClient.ReplyMessage(inputdata.replyToken, replyText); err != nil {
		log.Errorf("failed to reply message: %s", err)
		return InternalError
	}

	area := conditions[0]
	keyword := ""
	if len(conditions) == 2 {
		keyword = conditions[1]
	}
	log.Infof("area: %s, keyword: %s", area, keyword)

	shops, err := inputdata.ds.GetShopList(area, keyword)
	if err != nil {
		log.Errorf("failed to get shop list: %s", err)
		return err
	}

	shops.SortByRating()
	count := shops.GetOutputCount()
	shops = shops.GetTopShopList(count)
	if err != nil {
		log.Errorln("failed to get top shop list: %s", err)
		return err
	}
	log.Infoln("shops: ", shops)

	if err = t.snsClient.PushShopsMessage(inputdata.destinationUserID, shops); err != nil {
		log.Errorln("failed to push message: ", err)
		return err
	}
	return nil
}

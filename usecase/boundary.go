package usecase

import (
	"context"

	"github.com/jgvkmea/go-sort-tabelog/entity"
)

type IUseCase interface {
	GetShopsOrderByRating(context.Context, *GetShopsOrderByRatingInputData) error
}

type ShopDataSource interface {
	GetShopList(area string, keyword string) (entity.Shops, error)
}

type SNSClient interface {
	ReplyMessage(replyToken string, message string) error
	PushShopsMessage(userID string, shops entity.Shops) error
}

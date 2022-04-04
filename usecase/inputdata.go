package usecase

type GetShopsOrderByRatingInputData struct {
	ds                ShopDataSource
	message           string
	replyToken        string
	destinationUserID string
}

func NewGetShopsOrderByRatingInputData(ds ShopDataSource, message string, replyToken string, destinationUserID string) *GetShopsOrderByRatingInputData {
	return &GetShopsOrderByRatingInputData{
		ds,
		message,
		replyToken,
		destinationUserID,
	}
}

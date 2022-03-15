package usecase

type GetShopsOrderByRatingInputData struct {
	snsClient         SNSClient
	ds                ShopDataSource
	message           string
	replyToken        string
	destinationUserID string
}

func NewGetShopsOrderByRatingInputData(snsClient SNSClient, ds ShopDataSource, message string, replyToken string, destinationUserID string) *GetShopsOrderByRatingInputData {
	return &GetShopsOrderByRatingInputData{
		snsClient,
		ds,
		message,
		replyToken,
		destinationUserID,
	}
}

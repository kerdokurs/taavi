package external

import (
	"context"
	"net/http"
)

var meditationQuotesURL string

type meditationResponse struct {
	Quote string `json:"quote"`
}

func GetMeditationQuote(ctx context.Context) (string, error) {
	res, err := DoRequest[any, meditationResponse](ctx, http.MethodGet, meditationQuotesURL, nil)
	if err != nil {
		return "", err
	}
	return res.Quote, nil
}

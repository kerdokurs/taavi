package external

import (
	"context"
	"net/http"
)

var sleepQuotesURL string

type sleepQuoteResponse struct {
	Quote string `json:"quote"`
}

func GetSleepQuote(ctx context.Context) (string, error) {
	res, err := DoRequest[any, sleepQuoteResponse](ctx, http.MethodGet, sleepQuotesURL, nil)
	if err != nil {
		return "", err
	}
	return res.Quote, nil
}

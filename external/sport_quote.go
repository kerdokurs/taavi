package external

import (
	"context"
	"net/http"
)

var sportQuotesURL string

type sportQuoteResponse struct {
	Quote string `json:"quote"`
}

func GetSportQuote(ctx context.Context) (string, error) {
	res, err := DoRequest[any, sportQuoteResponse](ctx, http.MethodGet, sportQuotesURL, nil)
	if err != nil {
		return "", err
	}
	return res.Quote, nil
}

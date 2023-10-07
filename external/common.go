package external

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"time"

	"kerdo.dev/taavi/logger"
)

const requestTimeout = time.Second * 10

func Init() {
	meditationQuotesURL = os.Getenv("MEDITATION_URL")
	sleepQuotesURL = os.Getenv("SLEEP_URL")

	// ChatGPT
	dailyTaskListPrompt = os.Getenv("DAILY_DASKLIST_PROMPT")
}

func DoRequest[Req any, Res any](ctx context.Context, method string, url string, data *Req) (*Res, error) {
	ctx, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	var body io.Reader
	if data != nil {
		bytesData, err := json.Marshal(data)
		if err != nil {
			logger.Errorw("error marshalling data", logger.M{
				"err": err.Error(),
				"url": url,
			})
			return nil, err
		}
		body = bytes.NewBuffer(bytesData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	res := new(Res)
	if err := json.NewDecoder(resp.Body).Decode(res); err != nil {
		return nil, err
	}

	return res, nil
}

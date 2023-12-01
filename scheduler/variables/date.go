package variables

import (
	"context"
	"time"

	"gitlab.com/metakeule/fmtdate"
)

const defaultFormat = "DD/MM/YYYY hh:mm:ss"

const DateFormatKey = "date_format"

func GetDate(ctx context.Context) (string, error) {
	format := func() string {
		fmtAny := ctx.Value(DateFormatKey)
		if fmtAny == nil {
			return defaultFormat
		}

		if fmt, ok := fmtAny.(string); ok {
			return fmt
		}

		return defaultFormat
	}()

	now := time.Now()
	return fmtdate.Format(format, now), nil
}

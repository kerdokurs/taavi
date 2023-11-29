package variables

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/metakeule/fmtdate"
)

// For custom formatting check https://stackoverflow.com/a/20234207
const defaultFormat = "DD/MM/YYYY hh:mm:ss"

const DateFormatKey = "date_format"

func GetDate(ctx context.Context) (string, error) {
	format := func() string {
		fmtAny := ctx.Value(DateFormatKey)
		fmt.Println(fmtAny)
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

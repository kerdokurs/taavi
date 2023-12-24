package scheduler

import (
	"context"
	"regexp"
	"strings"

	"kerdo.dev/taavi/pkg/external"
	"kerdo.dev/taavi/pkg/logger"
	"kerdo.dev/taavi/pkg/scheduler/variables"
)

type GetVariable func(ctx context.Context) (string, error)

var variableTable = map[string]GetVariable{
	"meditation_quote": external.GetMeditationQuote,
	"sleep_quote":      external.GetSleepQuote,
	"sport_quote":      external.GetSportQuote,
	"daily_tasks":      external.GetDailyTaskList,
	"why_library":      external.GetLibraryMessage,
	"date":             variables.GetDate,
}

var variableRegex = regexp.MustCompile(`\{(?P<VariableName>[\w_-]+)\}`)

func replaceVariables(ctx context.Context, msg string) string {
	subMatch := variableRegex.FindStringSubmatch(msg)
	if len(subMatch) == 0 {
		return msg
	}

	if fn, ok := variableTable[subMatch[1]]; ok {
		result, err := fn(ctx)
		if err == nil {
			return strings.ReplaceAll(msg, subMatch[0], result)
		}
		logger.Errorw("error replacing variable", logger.M{
			"err":      err.Error(),
			"msg":      msg,
			"variable": subMatch[0],
		})
	}

	return strings.ReplaceAll(msg, subMatch[0], "")
}

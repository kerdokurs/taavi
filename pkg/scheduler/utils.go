package scheduler

import (
	"strconv"
	"strings"

	"kerdo.dev/taavi/pkg/logger"
)

// ab:cd, ef:gh -> [ab, cd, ef, gh]
func parseStartEnd(start, end string) [4]int {
	result := [4]int{}

	tmp := start + ":" + end
	var err error
	for i, part := range strings.Split(tmp, ":") {
		if result[i], err = strconv.Atoi(part); err != nil {
			logger.Errorw("error converting start/end time number", logger.M{
				"err": err.Error(),
			})
			result[i] = 0
		}
	}

	return result
}

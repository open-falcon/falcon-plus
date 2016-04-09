package utils

import (
	"fmt"
)

func Counter(metric string, tags map[string]string) string {
	if tags == nil || len(tags) == 0 {
		return metric
	}
	return fmt.Sprintf("%s/%s", metric, SortedTags(tags))
}

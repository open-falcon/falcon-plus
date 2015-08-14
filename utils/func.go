package utils

import (
	"fmt"
)

func PK(endpoint, metric string, tags map[string]string) string {
	if tags == nil || len(tags) == 0 {
		return fmt.Sprintf("%s/%s", endpoint, metric)
	}
	return fmt.Sprintf("%s/%s/%s", endpoint, metric, SortedTags(tags))
}

func UUID(endpoint, metric string, tags map[string]string, dstype string, step int) string {
	if tags == nil || len(tags) == 0 {
		return fmt.Sprintf("%s/%s/%s/%d", endpoint, metric, dstype, step)
	}
	return fmt.Sprintf("%s/%s/%s/%s/%d", endpoint, metric, SortedTags(tags), dstype, step)
}

func Checksum(endpoint string, metric string, tags map[string]string) string {
	pk := PK(endpoint, metric, tags)
	return Md5(pk)
}

func ChecksumOfUUID(endpoint, metric string, tags map[string]string, dstype string, step int64) string {
	return Md5(UUID(endpoint, metric, tags, dstype, int(step)))
}

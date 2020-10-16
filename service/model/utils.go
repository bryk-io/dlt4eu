package model

import (
	"fmt"
	"time"
)

func formatDate(value int64, format DateFormat) string {
	stamp := time.Unix(value, 0)
	switch format {
	case DateFormatRfc822:
		return stamp.Format(time.RFC822)
	case DateFormatRfc3339:
		return stamp.Format(time.RFC3339)
	case DateFormatUnix:
		return fmt.Sprintf("%d", stamp.Unix())
	default:
		return stamp.Format(time.RFC3339)
	}
}

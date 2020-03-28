package goodreads

import "time"

const goodreadsTimeLayout = "Mon Jan 2 15:04:05 -0700 2006"

func ParseDate(dateStr string) time.Time {
	if len(dateStr) > 0 {
		dateParsed, _ := time.Parse(goodreadsTimeLayout, dateStr)
		return dateParsed
	}
	return time.Time{}
}

package utils

import "time"

const layout = "2006-01-02"

func ParseDate(s string) time.Time {
	if s != "" && len(s) == 10 {
		t, _ := time.Parse(layout, s)
		return t
	}
	return time.Time{}
}

func FormatDate(t time.Time) string {
	return t.Format(layout)
}


func GetYear(d1 string) int {
	layout := "2006-01-02"
	t1, err1 := time.Parse(layout, d1)
	if err1 == nil {
		return t1.Year()
	}
	return -1
}

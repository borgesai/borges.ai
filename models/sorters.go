package models

type ByStatusDesc []Book

func (a ByStatusDesc) Len() int { return len(a) }
func (a ByStatusDesc) Less(i, j int) bool {
	if a[i].Review.ID > 0 &&  a[j].Review.ID >0 {
		return a[i].Review.StatusDate.Unix() > a[j].Review.StatusDate.Unix()
	}
	return true
}
func (a ByStatusDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type ByReadingDateDesc []Book

func (a ByReadingDateDesc) Len() int { return len(a) }
func (a ByReadingDateDesc) Less(i, j int) bool {
	if a[i].Reading.ID > 0 &&  a[j].Reading.ID >0 {
		return a[i].Reading.FinishDate.Unix() > a[j].Reading.FinishDate.Unix()
	}
	return true
}
func (a ByReadingDateDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type ByRatingDesc []Book

func (a ByRatingDesc) Len() int { return len(a) }
func (a ByRatingDesc) Less(i, j int) bool {
	return a[i].Review.Rating > a[j].Review.Rating
}
func (a ByRatingDesc) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

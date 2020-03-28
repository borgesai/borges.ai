package models

func AverageRating(books []Book) float32 {
	totalRating := 0
	totalRatedCount := 0
	for _, book := range books {
		if book.Review.Rating > 0 {
			totalRating += book.Review.Rating
			totalRatedCount += 1
		}
	}
	if totalRatedCount == 0 {
		return 0
	}
	return float32(totalRating) / float32(totalRatedCount)
}

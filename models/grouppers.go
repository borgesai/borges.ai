package models

type IndexedData struct {
	TotalPages int
	Items    []Book
}

func GroupBooksByFinishedYear(books []Book) map[int]IndexedData{
	byYear := make(map[int]IndexedData, 0)
	for _, book:=range books {
       if book.Reading.ID >0 && !book.Reading.FinishDate.IsZero(){
		   year := book.Reading.FinishDate.Year()
		   if originalData, ok := byYear[year]; ok {
			   if !ok {
				   data := IndexedData{}
				   data.TotalPages = book.GetNumPages()
				   yearBooks := make([]Book, 0)
				   yearBooks = append(yearBooks, book)
				   data.Items = yearBooks
				   byYear[year] = data
			   } else {
				   originalData.TotalPages = originalData.TotalPages + book.GetNumPages()
				   originalData.Items = append(originalData.Items, book)
				   byYear[year] = originalData
			   }
		   } else {
			   data := IndexedData{}

			   data.TotalPages = book.GetNumPages()
			   yearBooks := make([]Book, 0)
			   yearBooks = append(yearBooks, book)
			   data.Items = yearBooks
			   byYear[year] = data
		   }
	   }
	}
	return byYear
}

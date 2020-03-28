// Package goodreads provides a REST client for the public goodreads.com API.
//
// https://www.goodreads.com/api
package goodreads

import (
	"fmt"
	"net/url"
	"strconv"
)

// Client wraps the public Goodreads API.
type Client struct {
	ApiKey  string
	Decoder Decoder
}

// NewClient initializes a Client with default parameters.
func NewClient(key string) *Client {
	return &Client{
		ApiKey:  key,
		Decoder: DefaultDecoder,
	}
}

// AuthorBooks returns a list of books by a particular author.
// https://www.goodreads.com/api/index "Find books by title, author, or ISBN"
func (c *Client) SearchBooks(query string, page int) (*[]Work, error) {
	v := c.defaultValues()
	if query != "" {
		v.Set("q", query)
	}
	if page > 0 {
		v.Set("page", strconv.Itoa(page))
	}

	var r struct {
		Books []Work `xml:"search>results>work"`
	}
	err := c.Decoder.Decode(fmt.Sprintf("search/index.xml"), v, &r)
	if err != nil {
		return nil, err
	}
	return &r.Books, nil
}

// AuthorBooks returns a list of books by a particular author.
// https://www.goodreads.com/api/index#author.books
func (c *Client) AuthorBooks(authorID string, page int) (*Author, error) {
	v := c.defaultValues()
	if page > 0 {
		v.Set("page", strconv.Itoa(page))
	}

	var r struct {
		Author Author `xml:"author"`
	}
	err := c.Decoder.Decode(fmt.Sprintf("author/list/%s", authorID), v, &r)
	if err != nil {
		return nil, err
	}
	return &r.Author, nil
}

// AuthorShow returns the full details of an author.
// https://www.goodreads.com/api/index#author.show
func (c *Client) AuthorShow(authorID string) (*Author, error) {
	var r struct {
		Author Author `xml:"author"`
	}
	err := c.Decoder.Decode(fmt.Sprintf("author/show/%s", authorID), c.defaultValues(), &r)
	if err != nil {
		return nil, err
	}
	return &r.Author, nil
}

// ReviewList returns the books on a members shelf.
// https://www.goodreads.com/api/index#reviews.list
func (c *Client) ReviewList(userID int, shelf, sort, search, order string, page, perPage int) ([]Review, error) {
	v := c.defaultValues()
	userIDStr := strconv.Itoa(userID)
	v.Set("id", userIDStr)
	v.Set("v", "2")
	if shelf != "" {
		v.Set("shelf", shelf)
	}
	if sort != "" {
		v.Set("sort", sort)
	}
	if search != "" {
		v.Set("search", search)
	}
	if order != "" {
		v.Set("order", order)
	}
	if page > 0 {
		v.Set("page", strconv.Itoa(page))
	}
	if perPage > 0 {
		v.Set("per_page", strconv.Itoa(perPage))
	}

	var r struct {
		Reviews []Review `xml:"reviews>review"`
	}
	err := c.Decoder.Decode("review/list", v, &r)
	if err != nil {
		return nil, err
	}
	return r.Reviews, nil
}

// ShelvesList returns the list of shelves belonging to a user.
// https://www.goodreads.com/api/index#shelves.list
func (c *Client) ShelvesList(userID string) ([]UserShelf, error) {
	v := c.defaultValues()
	v.Set("user_id", userID)
	var r struct {
		Shelves []UserShelf `xml:"shelves>user_shelf"`
	}
	err := c.Decoder.Decode("shelves/list", v, &r)
	if err != nil {
		return nil, err
	}
	return r.Shelves, nil
}

// UserShow returns the public information about a given Goodreads user.
// https://www.goodreads.com/api/index#user.show
func (c *Client) UserShow(userID int) (*User, error) {
	var r struct {
		User User `xml:"user"`
	}
	userIDStr := strconv.Itoa(userID)
	err := c.Decoder.Decode(fmt.Sprintf("user/show/%s.xml", userIDStr), c.defaultValues(), &r)
	if err != nil {
		return nil, err
	}
	return &r.User, nil
}

// BookByTitle 
func (c *Client) BookByTitle(title string) (*Book, error) {
	var r struct {
		Book Book `xml:"book"`
	}
	err := c.Decoder.Decode("book/title.xml", c.bookByTitleValues(title), &r)
	if err != nil {
		return nil, err
	}
	return &r.Book, nil
}

// BookByGoodreadsID
func (c *Client) BookByGoodreadsID(goodreadsID string) (*Book, error) {
	var r struct {
		Book Book `xml:"book"`
	}
	err := c.Decoder.Decode(fmt.Sprintf("book/show/%s.xml", goodreadsID), c.defaultValues(), &r)
	if err != nil {
		return nil, err
	}
	return &r.Book, nil
}

// BookByGoodreadsID
func (c *Client) BookByISBN(isbn string) (*Book, error) {
	var r struct {
		Book Book `xml:"book"`
	}
	err := c.Decoder.Decode(fmt.Sprintf("book/isbn/%s", isbn), c.defaultValues(), &r)
	if err != nil {
		return nil, err
	}
	return &r.Book, nil
}

func (c *Client) defaultValues() url.Values {
	v := url.Values{}
	v.Set("key", c.ApiKey)
	return v
}

func (c *Client) bookByTitleValues(title string) url.Values {
	v := url.Values{}
	v.Set("key", c.ApiKey)
	v.Set("title", title)
	return v
}

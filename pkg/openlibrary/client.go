package openlibrary

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"library-backend/config"
	"library-backend/internal/models"
	"net/http"
	"sync"
)

type Client struct{}

func (c *Client) FetchBooksBySubject(subject string, limit, page int) ([]models.Book, error) {
	offset := page * limit
	url := fmt.Sprintf(config.Config.OpenLibBaseURL, subject, limit, offset)
	fmt.Println("Fetching URL:", url)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error making HTTP request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("non-200 response: %d - %s, body: %s", resp.StatusCode, http.StatusText(resp.StatusCode), string(body))
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding JSON: %w", err)
	}

	works, ok := result["works"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("no works found")
	}

	var books []models.Book
	var wg sync.WaitGroup
	var mu sync.Mutex

	for _, work := range works {
		wg.Add(1)
		go func(w interface{}) {
			defer wg.Done()
			workMap, ok := w.(map[string]interface{})
			if !ok {
				return
			}

			id, _ := workMap["key"].(string)
			title, _ := workMap["title"].(string)

			var authors []string
			if authorsList, ok := workMap["authors"].([]interface{}); ok {
				for _, author := range authorsList {
					if authorMap, ok := author.(map[string]interface{}); ok {
						if name, ok := authorMap["name"].(string); ok {
							authors = append(authors, name)
						}
					}
				}
			}

			editionCount := 0
			if ec, ok := workMap["edition_count"].(float64); ok {
				editionCount = int(ec)
			}

			firstPublishYear := 0
			if fpy, ok := workMap["first_publish_year"].(float64); ok {
				firstPublishYear = int(fpy)
			}

			coverImage := ""
			if coverID, ok := workMap["cover_id"].(float64); ok {
				coverImage = fmt.Sprintf(config.Config.OpenLibImageURL, int(coverID))
			}

			book := models.Book{
				ID:               id,
				Title:            title,
				Authors:          authors,
				EditionCount:     editionCount,
				FirstPublishYear: firstPublishYear,
				CoverImage:       coverImage,
				Genre:            subject,
			}

			mu.Lock()
			books = append(books, book)
			mu.Unlock()
		}(work)
	}

	wg.Wait()
	return books, nil
}

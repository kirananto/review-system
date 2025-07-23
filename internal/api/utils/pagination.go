package utils

import (
	"fmt"
	"net/http"
)

// GetPaginationLinks returns the previous and next page URLs for pagination
func GetPaginationLinks(r *http.Request, offset, limit, total int) (prevURL, nextURL *string) {
	var prev, next string

	// If not on the first page, add prev link
	if offset > 0 {
		prevOffset := offset - limit
		if prevOffset < 0 {
			prevOffset = 0
		}
		prev = generatePageURL(r, prevOffset, limit)
		prevURL = &prev
	}

	// If not on the last page, add next link
	if offset+limit < total {
		nextOffset := offset + limit
		next = generatePageURL(r, nextOffset, limit)
		nextURL = &next
	}

	return prevURL, nextURL
}

// generatePageURL creates a new URL with the given offset and limit
func generatePageURL(r *http.Request, offset, limit int) string {
	// Create a new copy of the URL to avoid modifying the original
	urlCopy := *r.URL

	// Get a copy of the query parameters
	q := urlCopy.Query()

	// Set the pagination parameters
	q.Set("offset", fmt.Sprintf("%d", offset))
	q.Set("limit", fmt.Sprintf("%d", limit))

	// Update the URL with the new query string
	urlCopy.RawQuery = q.Encode()

	return urlCopy.String()
}

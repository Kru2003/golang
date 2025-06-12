package utils

import "fmt"

func Paginate[T any](items []T, page, limit int) ([]T, error) {
	if page <= 0 {
		page = 1
	}
	if limit == 0 {
		limit = 10
	}
	if limit < 0 {
		return nil, fmt.Errorf("limit cannot be negative")
	}

	// Apply pagination
	start := (page - 1) * limit
	end := start + limit
	if start > len(items) {
		return nil, fmt.Errorf("page %d is out of range", page)
	}
	if end > len(items) {
		end = len(items)
	}

	return items[start:end], nil

}

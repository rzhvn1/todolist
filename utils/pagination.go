package utils

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type PaginationParams struct {
	Page   int
	Limit  int
	SortBy string
	Order  string
	Offset int
}

type PaginatedResponse struct {
	Page       int         `json:"page"`
	Limit      int         `json:"limit"`
	Total      int         `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

func ParsePaginationParams(r *http.Request, allowedSortFields []string) (PaginationParams, error) {
	query := r.URL.Query()

	// default values
	page := 1
	limit := 10
	sortBy := "createdAt"
	order := "asc"

	// parse page
	if pageStr := query.Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	// parse limit
	if limitStr := query.Get("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 {
			limit = l
		}
	}

	// parse sort_by
	if sortByQuery := query.Get("sort_by"); sortByQuery != "" {
		valid := false
		for _, field := range allowedSortFields {
			if sortByQuery == field {
				valid = true
				break
			}
		}
		if !valid {
			return PaginationParams{}, fmt.Errorf("invalid sort_by field: %s", sortByQuery)
		}

		sortBy = sortByQuery
	}

	// parse order
	orderQuery := strings.ToLower(query.Get("order"))
	if orderQuery == "desc" {
		order = "desc"
	} else if orderQuery != "" && orderQuery != "asc" {
		return PaginationParams{}, fmt.Errorf("invalid order: must be 'asc' or 'desc'")
	}

	// offset calculation
	offset := (page - 1) * limit

	return PaginationParams{
		Page: page,
		Limit: limit,
		SortBy: sortBy,
		Order: order,
		Offset: offset,
	}, nil
}

func WritePaginatedResponse(w http.ResponseWriter, page, limit, total int, data interface{}) {
	totalPages := (total + limit - 1)/limit

	response := PaginatedResponse{
		Page: page,
		Limit: limit,
		Total: total,
		TotalPages: totalPages,
		Data: data,
	}

	WriteJson(w, http.StatusOK, response)
}

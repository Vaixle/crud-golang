package httpquery

import (
	"strconv"
	"strings"
)

type FilterOption struct {
	Operator string
	Field    string
	Value    string
}

type Pagination struct {
	Limit int
	Page  int
}

func (p *Pagination) GetOffset() int {
	return p.Limit * p.Page
}

const (
	defaultLimit = 100
	defaultPage  = 0
)

func ParseQueryParams(values map[string][]string) ([]FilterOption, Pagination, error) {
	var filterOptions []FilterOption

	for field, value := range values {
		for _, v := range value {
			if strings.Contains(v, ":") {
				parts := strings.Split(v, ":")
				if len(parts) == 2 {
					operator := parts[0]
					queryValue := parts[1]
					switch operator {
					case "gt", "lt", "ge", "le", "eq", "ne", "order_by", "like":
						filterOptions = append(filterOptions, FilterOption{Operator: operator, Field: field, Value: queryValue})
						continue
					default:
						//TODO err
					}
				}
				//TODO err
			}
			filterOptions = append(filterOptions, FilterOption{Field: field, Value: v})
		}
	}

	pagination, err := getPagination(values)
	if err != nil {
		return nil, pagination, err
	}

	return filterOptions, pagination, nil
}

func getPagination(values map[string][]string) (Pagination, error) {
	var pagination Pagination

	if limitValues, ok := values["limit"]; ok == false {
		pagination.Limit = defaultLimit
	} else {
		limitStr := limitValues[0]
		if limitStr == "" {
			pagination.Limit = defaultLimit
		} else {
			limitInt, err := strconv.Atoi(limitStr)
			if err != nil {
				return pagination, err
			}

			if limitInt < 0 {
				pagination.Limit = defaultLimit
			} else {
				pagination.Limit = limitInt
			}
		}
	}

	if pageStrValues, ok := values["page"]; ok == false {
		pagination.Page = defaultPage
	} else {
		pageStr := pageStrValues[0]
		if pageStr == "" {
			pagination.Limit = defaultPage
		} else {
			pageInt, err := strconv.Atoi(pageStr)
			if err != nil {
				return pagination, err
			}

			if pageInt < 0 {
				pagination.Page = defaultPage
			} else {
				pagination.Page = pageInt
			}
		}
	}

	return pagination, nil
}

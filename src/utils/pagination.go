package utils

import (
	"errors"
	"fmt"
	"github.com/valyala/fasthttp"
	"gorm.io/gorm"
	"math"
	"strings"
)

// PaginationModel struct is used to return paginated data.
type PaginationModel struct {
	Limit     int         `json:"limit"`
	Page      int         `json:"page"`
	PageCount int         `json:"pageCount"`
	Total     int         `json:"total"`
	Result    interface{} `json:"result"`
}

// PaginationQuery builds a pagination query with the provided values
// and checks the input columns against the allowedColumns list.
// Returns a gorm query to be used in the function or an error.
func PaginationQuery(args *fasthttp.Args, allowedColumns map[string]bool) func(*gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		db = parseSearchLike(args.Peek("searchLike"), db, allowedColumns)
		db = parseSearchEq(args.Peek("searchEq"), db, allowedColumns)
		db = parseSearchLikeOr(args.Peek("searchLikeOr"), db, allowedColumns)
		db = parseSearchEqOr(args.Peek("searchEqOr"), db, allowedColumns)
		db = parseSearchIn(args.Peek("searchIn"), db, allowedColumns)
		db = parseSearchBetween(args.Peek("searchBetween"), db, allowedColumns)
		db = parseSortBy(args.Peek("sortBy"), db, allowedColumns)

		return db
	}
}

// PaginationCount calculates the page count with the given resultCount of a pagination query and a page limit.
func PaginationCount(resultCount int, limit int) int {
	return int(math.Ceil(float64(resultCount) / float64(limit)))
}

// PaginationOffset calculates the offset with the page and limit params
func PaginationOffset(page int, limit int) int {
	return (page - 1) * limit
}

// CreatePaginationModel is a helper to be able to return a pagination model in a single line
func CreatePaginationModel(limit int, page int, pageCount int, total int, result interface{}) PaginationModel {
	return PaginationModel{
		Limit:     limit,
		Page:      page,
		PageCount: pageCount,
		Total:     total,
		Result:    result,
	}
}

// search_like: for |where ... LIKE ... AND| query = search_like=column:value,column:value =>
// search_like=firstname:john,lastname:doe
func parseSearchLike(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	var statements []string
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		statements = append(statements, fmt.Sprintf("%s ILIKE '%%%s%%'", key, value))
	}

	/*
		for key, value := range paramMap {
			db = db.Where(fmt.Sprintf("%s ILIKE ?", key), fmt.Sprintf("%%%s%%", value))
		}
	*/

	if len(statements) > 0 {
		db = db.Where(strings.Join(statements, " AND "))
	}

	return db
}

// search_eq: for |where ... = ... AND| query = search_eq=column:value,column:value =>
// search_eq=firstname:john,lastname:doe
func parseSearchEq(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		db = db.Where(fmt.Sprintf("%s = ?", key), value)
	}

	return db
}

// search_like_or: for |where ... like ... OR| query = search_like_or=column:value,column:value =>
// search_or_like=firstname:john,lastname:doe
func parseSearchLikeOr(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	var statements []string
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		statements = append(statements, fmt.Sprintf("%s ILIKE '%%%s%%'", key, value))
	}

	if len(statements) > 0 {
		db = db.Where(strings.Join(statements, " OR "))
	}

	return db
}

// search_eq_or: for |where ... = ... OR| query = search_eq_or=column:value,column:value =>
// search_or_eq=firstname:john,lastname:doe
// TODO: Same issues as parseSearchLikeOr.
func parseSearchEqOr(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		db = db.Or(fmt.Sprintf("%s = ?", key), value)
	}

	return db
}

// search_in: for |where IN| query = search_in=column:value.value.value => search_in=is_online:true.false
func parseSearchIn(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseMultiValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		db = db.Where(fmt.Sprintf("%s IN (?)", key), value)
	}

	return db
}

// search_between  for |where ... between ... AND ...| query = search_between=column:value1.value2 =>
// search_between=created_at:2020-08-03.2020-09-03
func parseSearchBetween(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseMultiValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		if len(value) != 2 {
			err := db.AddError(errors.New("not exactly two values for between query"))
			if err != nil {
				return nil
			}
		}
		db = db.Where(fmt.Sprintf("%s BETWEEN ? AND ?", key), value[0], value[1])
	}

	return db
}

// sort_by: for |ORDER BY| query = sort_by=column:value,column:value => sort_by=firstname:asc,lastname:desc
func parseSortBy(params []byte, db *gorm.DB, allowedColumns map[string]bool) *gorm.DB {
	paramMap := parseSingleValueParams(db, string(params), allowedColumns)

	for key, value := range paramMap {
		if value == "desc" {
			db = db.Order(fmt.Sprintf("%s DESC", key))
		} else if value == "asc" {
			db = db.Order(fmt.Sprintf("%s ASC", key))
		} else {
			err := db.AddError(errors.New("order not asc or desc"))
			if err != nil {
				return nil
			}
		}
	}

	return db
}

// parseSingleValueParams parses the query string for single value params.
// The query string should be in the format of key:value,key:value
func parseSingleValueParams(db *gorm.DB, params string, allowedColumns map[string]bool) map[string]string {
	paramMap := make(map[string]string)

	if params != "" {
		paramSearchParts := strings.Split(params, ",")
		for _, paramSearchPart := range paramSearchParts {
			valuePairs := strings.Split(paramSearchPart, ":")
			canParse := len(valuePairs) == 2 && valuePairs[0] != "" && valuePairs[1] != ""
			isAllowed := allowedColumns[valuePairs[0]]

			if !canParse {
				err := db.AddError(errors.New("cannot parse invalid format"))
				if err != nil {
					return nil
				}
			}
			if !isAllowed {
				err := db.AddError(errors.New("column not allowed"))
				if err != nil {
					return nil
				}
			}
			if isAllowed && canParse {
				paramMap[valuePairs[0]] = valuePairs[1]
			}
		}
	}

	return paramMap
}

// parseMultiValueParams parses the query string for multi value params.
// The query string should be in the format of key:value.value.value,key:value.value.value
func parseMultiValueParams(db *gorm.DB, params string, allowedColumns map[string]bool) map[string][]string {
	paramMap := make(map[string][]string)

	if params != "" {
		paramSearchParts := strings.Split(params, ",")
		for _, paramSearchPart := range paramSearchParts {
			valuePairs := strings.SplitN(paramSearchPart, ":", 2)
			canParse := len(valuePairs) == 2 && valuePairs[0] != "" && valuePairs[1] != ""
			isAllowed := allowedColumns[valuePairs[0]]

			if !canParse {
				err := db.AddError(errors.New("cannot parse invalid format"))
				if err != nil {
					return nil
				}
			}
			if !isAllowed {
				err := db.AddError(errors.New("column not allowed"))
				if err != nil {
					return nil
				}
			}
			if isAllowed && canParse {
				paramMap[valuePairs[0]] = strings.Split(valuePairs[1], ".")
			}
		}
	}

	return paramMap
}
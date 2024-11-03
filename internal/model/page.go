package model

import (
	"fmt"

	"github.com/kataras/iris/v12/context"
)

const (
	PageSize = "pageSize" // Constant for the page size parameter
	PageNum  = "pageNum"  // Constant for the page number parameter
)

// Pagination constants
const (
	PageSizeD = 10 // Default page size
	PageNumD  = 1  // Default page number
)

// Page struct for pagination information
type Page struct {
	Total    int         `json:"total"`    // Total data amount
	PageSize int         `json:"pageSize"` // Size of each page
	PageNum  int         `json:"pageNum"`  // Page number
	Items    interface{} `json:"items"`    // Data list
}

// PageParam extracts pagination parameters from the context
func PageParam(ctx *context.Context) *Page {
	pageNum, err1 := ctx.URLParamInt(PageNum)
	if err1 != nil {
		pageNum = PageNumD // Use default page number if not provided
	}
	pageSize, err2 := ctx.URLParamInt(PageSize)
	if err2 != nil {
		pageSize = PageSizeD // Use default page size if not provided
	}

	return &Page{
		PageSize: pageSize,
		PageNum:  pageNum,
	}
}

// PageFilter filters the data based on page number and size
func PageFilter(num, size int, data []interface{}) (int, []interface{}, error) {
	total := len(data) // Total number of items in the data
	result := make([]interface{}, 0)

	if num < 1 {
		return 0, result, fmt.Errorf("page number must be greater than 0") // Error if page number is less than 1
	}
	if num*size < total {
		result = data[(num-1)*size : (num * size)] // Return the specified page of data
	} else {
		if (num-1)*size < total {
			result = data[(num-1)*size:] // Return remaining data if there are not enough items for a full page
		}
	}
	return total, result, nil
}

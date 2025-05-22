package abstraction

import (
	"fmt"
	"reflect"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"gorm.io/gorm"
)

var (
	DEFAULT_PAGE_SIZE       = 100
	DEFAULT_SORT_BY         = "modified_at"
	DEFAULT_ORDER_BY        = "desc"
	DEFAULT_CURSOR_FIELD    = "modified_at"
	DEFAULT_CURSOR_SORT_BY  = "modified_at"
	DEFAULT_CURSOR_ORDER_BY = "desc"
)

func ConvertSnakeToCamel(str string) string {
	// Split the string by underscores
	parts := strings.Split(str, "_")
	for i := range parts {
		// Use cases.Title to properly capitalize the first letter
		caser := cases.Title(language.English)
		parts[i] = caser.String(parts[i])

	}
	// Join the parts back together to form CamelCase
	return strings.Join(parts, "")
}

func (p *PaginationCursor) GetNextCursor(dataSlice reflect.Value) string {
	if dataSlice.Len() == 0 {
		return ""
	}
	lastData := dataSlice.Index(p.Limit() - 1) // Pastikan index tidak out of range
	if lastData.Kind() == reflect.Ptr {
		lastData = lastData.Elem()
	}
	fieldName := ConvertSnakeToCamel(p.Field)
	f := lastData.FieldByName(fieldName)
	if !f.IsValid() {
		return ""
	}

	if f.Kind() == reflect.Ptr {
		f = f.Elem()
	}
	var nextCursor string
	switch f.Kind() {
	case reflect.Int64:
		val := f.Int()
		nextCursor = fmt.Sprintf("%v", val)
	case reflect.String:
		val := f.String()
		nextCursor = val
	}
	return nextCursor
}

type PaginationCursor struct {
	PageSize int `json:"page_size" query:"page_size" default:"100"`
	Field    string
	Cursor   string  `json:"cursor"    query:"cursor"`
	SortBy   *string `json:"sort_by"   query:"sort_by"                 example:"id"`
	OrderBy  *string `json:"order_by"  query:"order_by"                             enums:"asc,desc"`
}

func (p *PaginationCursor) Limit() int {
	if p == nil || p.PageSize <= 0 || p.PageSize > DEFAULT_PAGE_SIZE {
		return DEFAULT_PAGE_SIZE
	}

	return p.PageSize
}

func (p *PaginationCursor) SetDefault() {
	if p.Field == "" {
		p.Field = DEFAULT_CURSOR_FIELD
	}
	if p.PageSize <= 0 || p.PageSize > DEFAULT_PAGE_SIZE {
		p.PageSize = DEFAULT_PAGE_SIZE
	}

	if IsStringBlank(p.SortBy) {
		p.SortBy = &DEFAULT_CURSOR_SORT_BY
	}

	if IsStringBlank(p.OrderBy) {
		p.OrderBy = &DEFAULT_CURSOR_ORDER_BY
	}
	var sort string
	switch *p.OrderBy {
	case "asc":
		sort = "asc"
	case "desc":
		sort = "desc"
	default:
		sort = "desc"
	}
	p.OrderBy = &sort
}

func (p *PaginationCursor) GetPageSize() int {
	if p == nil || p.PageSize <= 0 || p.PageSize > DEFAULT_PAGE_SIZE {
		return DEFAULT_PAGE_SIZE
	}
	return p.PageSize
}

type Pagination struct {
	Page     int     `json:"page" query:"page" default:"1"`
	PageSize int     `json:"page_size" query:"page_size" default:"100"`
	Cursor   string  `json:"cursor" query:"cursor"`
	SortBy   *string `json:"sort_by" query:"sort_by" example:"id"`
	OrderBy  *string `json:"order_by" query:"order_by" enums:"asc,desc"`
}

// set default before further processing
// do it after binding the pagination object in handler layer
func (p *Pagination) SetDefault() {
	if p == nil || p.PageSize <= 0 {
		p.PageSize = DEFAULT_PAGE_SIZE
	}
	if p == nil || p.Page <= 0 {
		p.Page = 1
	}

	var sort string

	if p == nil {
		return
	}

	if IsStringBlank(p.SortBy) {
		p.SortBy = &DEFAULT_SORT_BY
	}

	if *p.SortBy == "order" {
		*p.SortBy = fmt.Sprintf("\"%s\"", "order")
	}

	if IsStringBlank(p.OrderBy) {
		p.OrderBy = &DEFAULT_ORDER_BY
	}

	switch *p.OrderBy {
	case "asc":
		sort = "asc"
	case "desc":
		sort = "desc"
	default:
		sort = "desc"
	}
	p.OrderBy = &sort
}

func (p *Pagination) Limit() int {
	if p == nil || p.PageSize <= 0 {
		return DEFAULT_PAGE_SIZE
	}
	return p.PageSize
}

func (p *Pagination) Offset() int {
	if p == nil || p.Page <= 0 {
		return 0
	}
	return (p.Page - 1) * p.Limit()
}

func (p *Pagination) GetPage() int {
	if p == nil || p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p *Pagination) GetPageSize() int {
	if p == nil || p.PageSize <= 0 {
		return DEFAULT_PAGE_SIZE
	}
	return p.PageSize
}

func (p *Pagination) GetOrderBy() string {
	if p == nil || IsStringBlank(p.OrderBy) {
		return DEFAULT_ORDER_BY
	}
	return *p.OrderBy
}

func (p *Pagination) Apply(db *gorm.DB) {
	if p != nil {
		db.Limit(p.Limit())
		db.Offset(p.Offset())
		// if p.OrderBy != nil {
		db.Order(p.GetSortBy())
		// }
	}
}

func IsStringBlank(s *string) bool {
	if s == nil || strings.TrimSpace(*s) == "" {
		return true
	}
	return false
}

func (p *Pagination) ChangeDefaultSortingClause(sortBy string, orderBy *string) {
	// because this function is to populate the sorting
	// sortBy can't be empty string
	// orderBy might be empty
	// DEFAULT_ORDER_BY := "desc"
	if sortBy == "" {
		return
	}
	if IsStringBlank(p.SortBy) {
		p.SortBy = &DEFAULT_SORT_BY
	}

	if *p.SortBy == DEFAULT_SORT_BY {
		p.SortBy = &sortBy
	}

	if !IsStringBlank(orderBy) && IsStringBlank(p.OrderBy) {
		p.OrderBy = orderBy
	}
}

// Return SortBy value if not nil with OrderBy value with default asc
func (p *Pagination) GetSortBy() string {
	return *p.SortBy + " " + *p.OrderBy
}

func (p *Pagination) GetSorting() *Sorting {
	if IsStringBlank(p.OrderBy) {
		p.OrderBy = &DEFAULT_ORDER_BY
	}
	if p.SortBy != nil {
		return &Sorting{
			*p.SortBy,
			*p.OrderBy,
		}
	}

	return nil
}

type Sorting struct {
	SortBy  string `query:"sort_by"`
	OrderBy string `query:"order_by" enums:"asc,desc" default:"asc"`
}

func NewSorting(sortBy, sort string) *Sorting {
	if sort != "desc" {
		sort = "asc"
	}
	return &Sorting{
		SortBy:  sortBy,
		OrderBy: sort,
	}
}

type PaginationInfo struct {
	*Pagination
	*Sorting
	Count         int    `json:"count"`
	MoreRecords   bool   `json:"more_records"`
	NextCursor    string `json:"next_cursor"`
	TotalPageSize int    `json:"total_page"`
	TotalCount    int    `json:"total_count,omitempty"`
}

func (p *Pagination) CreatePageInfo(count int64) *PaginationInfo {
	return &PaginationInfo{
		Pagination:    p,
		Count:         int(count),
		TotalPageSize: int((count + int64(p.GetPageSize()) - 1) / int64(p.GetPageSize())),
	}
}

package tests

import (
	"reflect"
	"testing"
	"wakuwaku_nihongo/internals/abstraction"

	"github.com/stretchr/testify/assert"
)

func TestSetDefault(t *testing.T) {
	// Test 1: Default values when all fields are blank
	t.Run("Default values when fields are blank", func(t *testing.T) {
		p := &abstraction.PaginationCursor{}

		p.SetDefault()

		assert.Equal(t, abstraction.DEFAULT_CURSOR_FIELD, p.Field)
		assert.Equal(t, abstraction.DEFAULT_PAGE_SIZE, p.PageSize)
		assert.Equal(t, &abstraction.DEFAULT_CURSOR_SORT_BY, p.SortBy)
		assert.Equal(t, &abstraction.DEFAULT_CURSOR_ORDER_BY, p.OrderBy)
		assert.Equal(t, "desc", *p.OrderBy) // Default is "desc"
	})

	// Test 2: When Field is provided, it should not change
	t.Run("Field is provided", func(t *testing.T) {
		p := &abstraction.PaginationCursor{Field: "custom_field"}

		p.SetDefault()

		assert.Equal(t, "custom_field", p.Field)
	})

	// Test 3: When PageSize is valid (between 1 and abstraction.DEFAULT_PAGE_SIZE)
	t.Run("PageSize is valid", func(t *testing.T) {
		p := &abstraction.PaginationCursor{PageSize: 10}

		p.SetDefault()

		assert.Equal(t, 10, p.PageSize)
	})

	// Test 4: When PageSize is invalid (<= 0)
	t.Run("PageSize is invalid (<= 0)", func(t *testing.T) {
		p := &abstraction.PaginationCursor{PageSize: -5}

		p.SetDefault()

		assert.Equal(t, abstraction.DEFAULT_PAGE_SIZE, p.PageSize) // Should fall back to abstraction.DEFAULT_PAGE_SIZE
	})

	// Test 5: When PageSize is greater than abstraction.DEFAULT_PAGE_SIZE
	t.Run("PageSize is greater than abstraction.DEFAULT_PAGE_SIZE", func(t *testing.T) {
		p := &abstraction.PaginationCursor{PageSize: abstraction.DEFAULT_PAGE_SIZE + 1}

		p.SetDefault()

		assert.Equal(t, abstraction.DEFAULT_PAGE_SIZE, p.PageSize) // Should fall back to abstraction.DEFAULT_PAGE_SIZE
	})

	// Test 6: When SortBy is set
	t.Run("SortBy is set", func(t *testing.T) {
		p := &abstraction.PaginationCursor{SortBy: strPtr("custom_sort")}

		p.SetDefault()

		assert.Equal(t, strPtr("custom_sort"), p.SortBy) // Should remain custom_sort
	})

	// Test 7: When OrderBy is set to "asc"
	t.Run("OrderBy is set to asc", func(t *testing.T) {
		p := &abstraction.PaginationCursor{OrderBy: strPtr("asc")}

		p.SetDefault()

		assert.Equal(t, "asc", *p.OrderBy) // Should remain "asc"
	})

	// Test 8: When OrderBy is set to "desc"
	t.Run("OrderBy is set to desc", func(t *testing.T) {
		p := &abstraction.PaginationCursor{OrderBy: strPtr("desc")}

		p.SetDefault()

		assert.Equal(t, "desc", *p.OrderBy) // Should remain "desc"
	})

	// Test 9: When OrderBy is invalid (e.g. "unknown")
	t.Run("OrderBy is invalid", func(t *testing.T) {
		p := &abstraction.PaginationCursor{OrderBy: strPtr("unknown")}

		p.SetDefault()

		assert.Equal(t, "desc", *p.OrderBy) // Should fall back to "desc"
	})
}

func TestConvertSnakeToCamel(t *testing.T) {
	// Test 1: Convert simple snake_case to CamelCase
	t.Run("Convert simple snake_case to CamelCase", func(t *testing.T) {
		input := "simple_test_case"
		expected := "SimpleTestCase"

		result := abstraction.ConvertSnakeToCamel(input)

		assert.Equal(t, expected, result)
	})

	// Test 2: Convert snake_case with single word (no underscores)
	t.Run("Convert snake_case with single word", func(t *testing.T) {
		input := "word"
		expected := "Word"

		result := abstraction.ConvertSnakeToCamel(input)

		assert.Equal(t, expected, result)
	})

	// Test 3: Convert snake_case with multiple words
	t.Run("Convert snake_case with multiple words", func(t *testing.T) {
		input := "convert_this_example"
		expected := "ConvertThisExample"

		result := abstraction.ConvertSnakeToCamel(input)

		assert.Equal(t, expected, result)
	})

	// Test 4: Convert snake_case with all upper case letters
	t.Run("Convert snake_case with all upper case letters", func(t *testing.T) {
		input := "ALL_UPPER_CASE"
		expected := "AllUpperCase"

		result := abstraction.ConvertSnakeToCamel(input)

		assert.Equal(t, expected, result)
	})

	// Test 5: Convert an empty string (should return empty string)
	t.Run("Empty string returns empty string", func(t *testing.T) {
		input := ""
		expected := ""

		result := abstraction.ConvertSnakeToCamel(input)

		assert.Equal(t, expected, result)
	})

	// Test 6: Convert snake_case with leading/trailing underscores
	t.Run("Snake_case with leading/trailing underscores", func(t *testing.T) {
		input := "_leading_trailing_"
		expected := "LeadingTrailing"

		result := abstraction.ConvertSnakeToCamel(input)

		assert.Equal(t, expected, result)
	})

	// Test 7: Single word with leading underscores
	t.Run("Single word with leading underscores", func(t *testing.T) {
		input := "_word"
		expected := "Word"

		result := abstraction.ConvertSnakeToCamel(input)

		assert.Equal(t, expected, result)
	})
}

func int64Ptr(i int64) *int64 {
	return &i
}
func TestGetNextCursor(t *testing.T) {
	// Helper to create a pointer to int64
	type SampleInt struct {
		Id *int64
	}

	type SampleStr struct {
		Code *string
	}

	t.Run("Valid int64 field", func(t *testing.T) {
		slice := []*SampleInt{
			{Id: int64Ptr(1)},
			{Id: int64Ptr(99)},
		}
		data := reflect.ValueOf(slice)
		cursor := &abstraction.PaginationCursor{Field: "id", PageSize: 2}
		cursor.SetDefault()

		result := cursor.GetNextCursor(data)
		assert.Equal(t, "99", result)
	})

	t.Run("Valid string field", func(t *testing.T) {
		slice := []*SampleStr{
			{Code: strPtr("abc")},
			{Code: strPtr("xyz")},
		}
		data := reflect.ValueOf(slice)
		cursor := &abstraction.PaginationCursor{Field: "code", PageSize: 2}
		cursor.SetDefault()

		result := cursor.GetNextCursor(data)
		assert.Equal(t, "xyz", result)
	})

	t.Run("Last element field is nil (int64)", func(t *testing.T) {
		slice := []*SampleInt{
			{Id: int64Ptr(1)},
			{Id: nil},
		}
		data := reflect.ValueOf(slice)
		cursor := &abstraction.PaginationCursor{Field: "id", PageSize: 2}
		cursor.SetDefault()

		result := cursor.GetNextCursor(data)
		assert.Equal(t, "", result)
	})

	t.Run("Last element field is nil (string)", func(t *testing.T) {
		slice := []*SampleStr{
			{Code: strPtr("hello")},
			{Code: nil},
		}
		data := reflect.ValueOf(slice)
		cursor := &abstraction.PaginationCursor{Field: "code", PageSize: 2}
		cursor.SetDefault()

		result := cursor.GetNextCursor(data)
		assert.Equal(t, "", result)
	})

	t.Run("Empty slice returns empty string", func(t *testing.T) {
		slice := []*SampleInt{}
		data := reflect.ValueOf(slice)
		cursor := &abstraction.PaginationCursor{Field: "id", PageSize: 1}
		cursor.SetDefault()

		result := cursor.GetNextCursor(data)
		assert.Equal(t, "", result)
	})

	t.Run("Invalid field name returns empty string", func(t *testing.T) {
		slice := []*SampleInt{
			{Id: int64Ptr(42)},
		}
		data := reflect.ValueOf(slice)
		cursor := &abstraction.PaginationCursor{Field: "nonexistent_field", PageSize: 1}
		cursor.SetDefault()

		result := cursor.GetNextCursor(data)
		assert.Equal(t, "", result)
	})
}

func TestPagination_Limit(t *testing.T) {
	type fields struct {
		Page     int
		PageSize int
	}
	tests := []struct {
		name   string
		fields fields
		want   int
	}{
		{
			name: "PageSize = 0",
			fields: fields{
				PageSize: 0,
			},
			want: abstraction.DEFAULT_PAGE_SIZE,
		},
		{
			name: "PageSize = 1",
			fields: fields{
				PageSize: 1,
			},
			want: 1,
		},
		{
			name: "PageSize = 15",
			fields: fields{
				PageSize: 7,
			},
			want: 7,
		},
		{
			name: "PageSize = negative",
			fields: fields{
				PageSize: -12,
			},
			want: abstraction.DEFAULT_PAGE_SIZE,
		},
		{
			name: "PageSize = some big number",
			fields: fields{
				PageSize: 1000,
			},
			want: 1000,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &abstraction.Pagination{
				Page:     tt.fields.Page,
				PageSize: tt.fields.PageSize,
			}
			if got := p.Limit(); got != tt.want {
				t.Errorf("Limit() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_Offset(t *testing.T) {
	type fields struct {
		Page     int
		PageSize int
	}
	tests := []struct {
		name string
		data fields
		want int
	}{
		{
			name: "page = neg and pageSize = neg",
			data: fields{
				Page:     -1,
				PageSize: -1,
			},
			want: 0,
		},
		{
			name: "page = 0 and pageSize = neg",
			data: fields{
				Page:     0,
				PageSize: -1,
			},
			want: 0,
		},
		{
			name: "page = pos and pageSize = neg",
			data: fields{
				Page:     3,
				PageSize: -1,
			},
			want: 200,
		},
		{
			name: "page = neg and pageSize = zero",
			data: fields{
				Page:     -1,
				PageSize: 0,
			},
			want: 0,
		},
		{
			name: "page = 0 and pageSize = zero",
			data: fields{
				Page:     0,
				PageSize: 0,
			},
			want: 0,
		},
		{
			name: "page = pos and pageSize = zero",
			data: fields{
				Page:     3,
				PageSize: 0,
			},
			want: 200,
		},
		{
			name: "page = neg and pageSize = pos",
			data: fields{
				Page:     -1,
				PageSize: 4,
			},
			want: 0,
		},
		{
			name: "page = 0 and pageSize = pos",
			data: fields{
				Page:     0,
				PageSize: 4,
			},
			want: 0,
		},
		{
			name: "page = pos and pageSize = pos",
			data: fields{
				Page:     3,
				PageSize: 4,
			},
			want: 8,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &abstraction.Pagination{
				Page:     tt.data.Page,
				PageSize: tt.data.PageSize,
			}
			if got := p.Offset(); got != tt.want {
				t.Errorf("Offset() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPagination_SetDefault_pagetest(t *testing.T) {
	type fields struct {
		Page     int
		PageSize int
	}
	test := []struct {
		name string
		data fields
		want fields
	}{
		{
			name: "page = neg and pageSize = neg",
			data: fields{
				Page:     -1,
				PageSize: -1,
			},
			want: fields{
				Page:     1,
				PageSize: abstraction.DEFAULT_PAGE_SIZE,
			},
		},
		{
			name: "page = 0 and pageSize = neg",
			data: fields{
				Page:     0,
				PageSize: -1,
			},
			want: fields{
				Page:     1,
				PageSize: abstraction.DEFAULT_PAGE_SIZE,
			},
		},
		{
			name: "page = pos and pageSize = neg",
			data: fields{
				Page:     3,
				PageSize: -1,
			},
			want: fields{
				Page:     3,
				PageSize: abstraction.DEFAULT_PAGE_SIZE,
			},
		},
		{
			name: "page = neg and pageSize = zero",
			data: fields{
				Page:     -1,
				PageSize: 0,
			},
			want: fields{
				Page:     1,
				PageSize: abstraction.DEFAULT_PAGE_SIZE,
			},
		},
		{
			name: "page = 0 and pageSize = zero",
			data: fields{
				Page:     0,
				PageSize: 0,
			},
			want: fields{
				Page:     1,
				PageSize: abstraction.DEFAULT_PAGE_SIZE,
			},
		},
		{
			name: "page = pos and pageSize = zero",
			data: fields{
				Page:     3,
				PageSize: 0,
			},
			want: fields{
				Page:     3,
				PageSize: abstraction.DEFAULT_PAGE_SIZE,
			},
		},
		{
			name: "page = neg and pageSize = pos",
			data: fields{
				Page:     -1,
				PageSize: 4,
			},
			want: fields{
				Page:     1,
				PageSize: 4,
			},
		},
		{
			name: "page = 0 and pageSize = pos",
			data: fields{
				Page:     0,
				PageSize: 4,
			},
			want: fields{
				Page:     1,
				PageSize: 4,
			},
		},
		{
			name: "page = pos and pageSize = pos",
			data: fields{
				Page:     3,
				PageSize: 4,
			},
			want: fields{
				Page:     3,
				PageSize: 4,
			},
		},
	}
	for _, tt := range test {
		t.Run(tt.name, func(t *testing.T) {
			p := &abstraction.Pagination{
				Page:     tt.data.Page,
				PageSize: tt.data.PageSize,
			}
			p.SetDefault()
			if p.Page != tt.want.Page || p.PageSize != tt.want.PageSize {
				t.Errorf("PAGE: get  %d want %d; PAGESIZE: get %d want %d",
					p.Page, tt.want.Page, p.PageSize, tt.want.PageSize)
			}
		})
	}
}

func TestPagiation_IsStringBlank(t *testing.T) {
	// map
	type field struct {
		name string
		data *string
		want bool
	}
	emptystring := ""
	notemptystring := "not empty"
	tests := []field{
		{
			name: "nil",
			want: true,
		},
		{
			name: "empty",
			data: &emptystring,
			want: true,
		},
		{
			name: "not empty",
			data: &notemptystring,
			want: false,
		},
	}

	for _, val := range tests {
		t.Run(val.name, func(t *testing.T) {
			if got := abstraction.IsStringBlank(val.data); got != val.want {
				t.Errorf("get: %v, want %v, data %v", got, val.want, val.data)
			}
		})
	}
}

func TestPagination_GetSortBy(t *testing.T) {
	// combination of sortBy and orderBy field
	// each can be nil, empty string, and ordinary string

	// var nullString *string
	emptyString := ""
	normalString := "normalString"
	randomString := "randomString"
	asc := "asc"
	desc := "desc"

	type field struct {
		name    string
		sortBy  *string
		orderBy *string
		want    string
	}

	tests := []field{
		{
			name:    "",
			want:    abstraction.DEFAULT_SORT_BY + " " + desc,
			orderBy: &randomString,
			sortBy:  nil,
		},
		{
			name:    "",
			want:    "modified_at" + " " + asc,
			orderBy: &asc,
			sortBy:  nil,
		},
		{
			name:    "",
			want:    "modified_at" + " " + asc,
			orderBy: &asc,
			sortBy:  &emptyString,
		},
		{
			name:    "",
			want:    normalString + " " + asc,
			orderBy: &asc,
			sortBy:  &normalString,
		},

		{
			name:    "",
			want:    "modified_at" + " " + desc,
			orderBy: &emptyString,
			sortBy:  nil,
		},
		{
			name:    "",
			want:    "modified_at" + " " + desc,
			orderBy: &emptyString,
			sortBy:  &emptyString,
		},
		{
			name:    "",
			want:    normalString + " " + desc,
			orderBy: &emptyString,
			sortBy:  &normalString,
		},

		{
			name:    "",
			want:    "modified_at" + " " + desc,
			orderBy: nil,
			sortBy:  nil,
		},
		{
			name:    "",
			want:    "modified_at" + " " + desc,
			orderBy: nil,
			sortBy:  &emptyString,
		},
		{
			name:    "",
			want:    normalString + " " + desc,
			orderBy: nil,
			sortBy:  &normalString,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			p := new(abstraction.Pagination)
			p.OrderBy = test.orderBy
			p.SortBy = test.sortBy
			p.SetDefault()
			got := p.GetSortBy()
			if got != test.want {
				t.Errorf("get %s, want %s", got, test.want)
			}
		})
	}
}

func strPtr(s string) *string {
	return &s
}

func TestChangeDefaultSortingClause_Updated(t *testing.T) {
	tests := []struct {
		name         string
		initialSort  *string
		initialOrder *string
		sortBy       string
		orderBy      *string
		wantSort     *string
		wantOrder    *string
	}{
		{
			name:         "Sets to abstraction.DEFAULT_SORT_BY then overrides with sortBy",
			initialSort:  nil,
			initialOrder: nil,
			sortBy:       abstraction.DEFAULT_SORT_BY,
			orderBy:      strPtr("asc"),
			wantSort:     strPtr(abstraction.DEFAULT_SORT_BY),
			wantOrder:    strPtr("asc"),
		},
		{
			name:         "Change SortBy when it equals abstraction.DEFAULT_SORT_BY",
			initialSort:  strPtr(abstraction.DEFAULT_SORT_BY),
			initialOrder: nil,
			sortBy:       abstraction.DEFAULT_SORT_BY,
			orderBy:      nil,
			wantSort:     strPtr(abstraction.DEFAULT_SORT_BY),
			wantOrder:    nil,
		},
		{
			name:         "Does not change sortBy if already set and not default",
			initialSort:  strPtr("email"),
			initialOrder: nil,
			sortBy:       abstraction.DEFAULT_SORT_BY,
			orderBy:      strPtr("desc"),
			wantSort:     strPtr("email"),
			wantOrder:    strPtr("desc"),
		},
		{
			name:         "Sets default first, then applies sortBy if default",
			initialSort:  nil,
			initialOrder: strPtr("asc"),
			sortBy:       "name",
			orderBy:      nil,
			wantSort:     strPtr("name"),
			wantOrder:    strPtr("asc"),
		},
		{
			name:         "Skips everything when sortBy is empty",
			initialSort:  nil,
			initialOrder: nil,
			sortBy:       "",
			orderBy:      strPtr("desc"),
			wantSort:     nil,
			wantOrder:    nil,
		},
		{
			name:         "Does not override OrderBy if already set",
			initialSort:  strPtr(abstraction.DEFAULT_SORT_BY),
			initialOrder: strPtr("asc"),
			sortBy:       "updated_at",
			orderBy:      strPtr("desc"),
			wantSort:     strPtr("updated_at"),
			wantOrder:    strPtr("asc"),
		},
		{
			name:         "Skips OrderBy if orderBy is blank",
			initialSort:  strPtr(abstraction.DEFAULT_SORT_BY),
			initialOrder: nil,
			sortBy:       "email",
			orderBy:      strPtr(" "),
			wantSort:     strPtr("email"),
			wantOrder:    nil,
		},
		{
			name:         "Only changes OrderBy if it was blank",
			initialSort:  strPtr("name"),
			initialOrder: nil,
			sortBy:       "email",
			orderBy:      strPtr("asc"),
			wantSort:     strPtr("name"), // unchanged
			wantOrder:    strPtr("asc"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &abstraction.Pagination{
				SortBy:  tt.initialSort,
				OrderBy: tt.initialOrder,
			}

			p.ChangeDefaultSortingClause(tt.sortBy, tt.orderBy)

			if tt.wantSort == nil {
				assert.Nil(t, p.SortBy)
			} else {
				assert.NotNil(t, p.SortBy)
				assert.Equal(t, *tt.wantSort, *p.SortBy)
			}

			if tt.wantOrder == nil {
				assert.Nil(t, p.OrderBy)
			} else {
				assert.NotNil(t, p.OrderBy)
				assert.Equal(t, *tt.wantOrder, *p.OrderBy)
			}
		})
	}
}

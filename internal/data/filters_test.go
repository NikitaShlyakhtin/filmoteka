package data

import (
	"filmoteka/internal/validator"
	"testing"
)

func TestValidateFilters(t *testing.T) {
	v := validator.New()

	t.Run("ValidSortValue", func(t *testing.T) {
		f := Filters{
			Sort:         "title",
			SortSafelist: []string{"title", "year", "rating", "-title", "-year", "-rating"},
		}

		ValidateFilters(v, f)

		if !v.Valid() {
			t.Error(v.Errors)
		}
	})

	t.Run("InvalidSortValue", func(t *testing.T) {
		f := Filters{
			Sort:         "invalid",
			SortSafelist: []string{"title", "year", "rating", "-title", "-year", "-rating"},
		}

		ValidateFilters(v, f)

		if v.Valid() {
			t.Error("expected error for invalid sort value")
		}
	})
}

func TestSortColumn(t *testing.T) {
	t.Run("ValidSortColValue", func(t *testing.T) {
		f := Filters{
			Sort:         "title",
			SortSafelist: []string{"title", "year", "rating", "-title", "-year", "-rating"},
		}

		expected := "title"
		actual := f.sortColumn()

		if actual != expected {
			t.Errorf("Expected sort column to be %s, but got %s", expected, actual)
		}
	})

	t.Run("InvalidSortColValue", func(t *testing.T) {
		f := Filters{
			Sort:         "invalid",
			SortSafelist: []string{"title", "year", "rating", "-title", "-year", "-rating"},
		}

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic for invalid sort value")
			}
		}()

		f.sortColumn()
	})
}

func TestSortDirection(t *testing.T) {
	t.Run("SortDirectionASC", func(t *testing.T) {
		f := Filters{
			Sort:         "title",
			SortSafelist: []string{"title", "year", "rating", "-title", "-year", "-rating"},
		}

		expected := "ASC"
		actual := f.sortDirection()

		if actual != expected {
			t.Errorf("Expected sort direction to be %s, but got %s", expected, actual)
		}
	})

	t.Run("SortDirectionDESC", func(t *testing.T) {
		f := Filters{
			Sort:         "-title",
			SortSafelist: []string{"title", "year", "rating", "-title", "-year", "-rating"},
		}

		expected := "DESC"
		actual := f.sortDirection()

		if actual != expected {
			t.Errorf("Expected sort direction to be %s, but got %s", expected, actual)
		}
	})
}

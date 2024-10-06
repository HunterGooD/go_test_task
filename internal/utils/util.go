package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
)

func MergeSongParams(query *entity.SongListQueryParams, filter *entity.SongFilters) *entity.SongFilters {
	// ID
	if filter.ID == 0 {
		filter.ID = query.ID
	}

	// Name
	if filter.Name == "" {
		filter.Name = query.Name
	}

	// Link
	if filter.Link == "" {
		filter.Link = query.Link
	}

	// Text
	if filter.Text == "" {
		filter.Text = query.Text
	}

	// ReleaseDate (проверяем, что значение не нулевое)
	if filter.ReleaseDate.IsZero() {
		filter.ReleaseDate = query.ReleaseDate
	}

	// GroupName
	if filter.GroupName == "" {
		filter.GroupName = query.GroupName
	}

	return filter
}

func GetFilterString(filters *entity.SongFilters) string {
	conditions := []string{}

	// get value struct without pointer
	refStruct := reflect.ValueOf(filters).Elem()

	for i := 0; i < refStruct.NumField(); i++ {
		field := refStruct.Type().Field(i) // get field info
		value := refStruct.Field(i)        // Get value field

		if isZero(value) { // ignore zero values
			continue
		}

		fieldName := field.Tag.Get("db") // Используем имя из тега `db`

		switch value.Kind() {
		case reflect.String:
			conditions = append(conditions, fmt.Sprintf("%s LIKE '%s%%'", fieldName, value.String()))
		case reflect.Int, reflect.Int64:
			conditions = append(conditions, fmt.Sprintf("%s = %d", fieldName, value.Int()))
		case reflect.Struct: // Работаем с временем
			if field.Type == reflect.TypeOf(time.Time{}) {
				conditions = append(conditions, fmt.Sprintf("%s = '%s'", fieldName, value.Interface().(time.Time).Format("2006-01-02")))
			}
		default:
			// TODO: for anything data
		}
	}

	// Join to where statement with AND
	if len(conditions) > 0 {
		return "AND " + strings.Join(conditions, " AND ")
	}
	return ""
}

func isZero(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr && v.IsNil() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

package utils

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
)

func MergeSongParams(query *entity.SongListQueryParams, filter *entity.SongFilters) *entity.SongFilters {

	if filter == nil {
		filter = &entity.SongFilters{}
	}

	if query == nil {
		query = &entity.SongListQueryParams{}
	}

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
	if filter.ReleaseDate == nil {
		filter.ReleaseDate = query.ReleaseDate
	}

	// GroupName
	if filter.GroupName == "" {
		filter.GroupName = query.GroupName
	}

	return filter
}

func GetFilterString(startPlaceholders int, filters *entity.SongFilters) (string, []any) {
	conditions := []string{}
	args := make([]any, 0)
	argCounter := startPlaceholders

	if filters == nil {
		return "", args
	}

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
			// Добавляем placeholder и значение в args
			conditions = append(conditions, fmt.Sprintf("%s LIKE $%d", fieldName, argCounter))
			args = append(args, value.String()+"%")
			argCounter++
		case reflect.Int, reflect.Int64:
			conditions = append(conditions, fmt.Sprintf("%s = $%d", fieldName, argCounter))
			args = append(args, value.Int())
			argCounter++
		case reflect.Struct: // Работаем с временем
			if field.Type == reflect.TypeOf(time.Time{}) {
				conditions = append(conditions, fmt.Sprintf("%s = $%d", fieldName, argCounter))
				args = append(args, value.Interface().(time.Time))
				argCounter++
			}
		case reflect.Pointer:
			if value.Type() == reflect.TypeOf(&time.Time{}) {
				timeValue := value.Interface().(*time.Time)
				if timeValue != nil {
					conditions = append(conditions, fmt.Sprintf("%s = $%d", fieldName, argCounter))
					args = append(args, *timeValue) // Разыменовываем указатель
					argCounter++
				}
			}
		default:
			// TODO: обработка других типов данных
		}
	}

	// Join to where statement with AND
	if len(conditions) > 0 {
		return "AND " + strings.Join(conditions, " AND "), args
	}
	return "", args
}

func isZero(v reflect.Value) bool {
	return v.Kind() == reflect.Ptr && v.IsNil() || reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}

package cli

import (
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
	"time"
	"unicode"
)

func DisplayList(items interface{}) error {
	sliceValue := reflect.ValueOf(items)
	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("input must be a slice")
	}
	if sliceValue.Len() == 0 {
		fmt.Println("No items found.")
		return nil
	}

	// Check if the elements are structs
	elemType := sliceValue.Type().Elem()
	if elemType.Kind() != reflect.Struct {
		return fmt.Errorf("slice elements must be structs, got %s", elemType.Kind())
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	headers := []string{}
	for i := range elemType.NumField() {
		field := elemType.Field(i)
		if !field.IsExported() {
			continue
		}
		headers = append(headers, humanizeFieldName(field.Name))
	}
	fmt.Fprintf(w, "%s\n", strings.Join(headers, "\t"))

	for i := range sliceValue.Len() {
		elemValue := sliceValue.Index(i)
		values := []string{}
		for j := range elemType.NumField() {
			field := elemType.Field(j)
			if !field.IsExported() {
				continue
			}
			fieldValue := elemValue.Field(j).Interface()
			fieldName := elemType.Field(j).Name
			values = append(values, formatField(fieldName, fieldValue))
		}
		fmt.Fprintf(w, "%s\n", strings.Join(values, "\t"))
	}

	if err := w.Flush(); err != nil {
		log.Printf("Failed to display list: %v", err)
		return err
	}
	return nil
}

func humanizeFieldName(name string) string {
	if name == "ID" {
		return name
	}

	var result []rune
	for i, r := range name {
		if i > 0 && unicode.IsUpper(r) {
			result = append(result, ' ')
		}
		result = append(result, r)
	}
	return string(result)
}

func formatField(name string, value interface{}) string {
	if value == nil {
		return ""
	}

	ts, ok := value.(int64)
	if !ok {
		return fmt.Sprintf("%v", value) // Default to %v if not int64
	}

	// Check if the field name suggests a timestamp
	if strings.HasSuffix(name, "At") {
		return time.Unix(ts, 0).Format("2006-01-02 15:04")
	}

	return fmt.Sprintf("%d", ts)
}

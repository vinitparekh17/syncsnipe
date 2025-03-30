package cli

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"text/tabwriter"
	"time"
	"unicode"
)

func DisplayList(items any) error {
	sliceValue := reflect.ValueOf(items)
	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("input must be a slice")
	}
	if sliceValue.Len() == 0 {
		fmt.Println("No items found.")
		return nil
	}

	// Use the first element to extract headers.
	firstElem := sliceValue.Index(0).Interface()
	headers, _, err := extractHeadersAndValues(firstElem)
	if err != nil {
		return err
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "%s\n", strings.Join(headers, "\t"))

	// Process each struct element.
	for i := range sliceValue.Len() {
		_, values, err := extractHeadersAndValues(sliceValue.Index(i).Interface())
		if err != nil {
			return err
		}
		fmt.Fprintf(w, "%s\n", strings.Join(values, "\t"))
	}

	if err := w.Flush(); err != nil {
		log.Printf("Failed to display list: %v", err)
		return err
	}
	return nil
}

func DisplayStruct(item any) error {
	headers, values, err := extractHeadersAndValues(item)
	if err != nil {
		return err
	}
	for i, header := range headers {
		fmt.Printf("%s: %s\n", header, values[i])
	}
	return nil
}

func extractHeadersAndValues(item interface{}) (headers, values []string, err error) {
	v := reflect.ValueOf(item)
	if v.Kind() != reflect.Struct {
		return nil, nil, fmt.Errorf("expected a struct, got %s", v.Kind())
	}
	t := v.Type()
	for i := range t.NumField() {
		field := t.Field(i)
		if !field.IsExported() {
			continue
		}
		headers = append(headers, humanizeFieldName(field.Name))
		values = append(values, formatField(field.Name, v.Field(i).Interface()))
	}
	return headers, values, nil
}

func humanizeFieldName(name string) string {
	if strings.HasSuffix(name, "ID") {
		if name == "ID" {
			return name
		}
		return strings.TrimSuffix(name, "ID") + " ID"
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

func formatField(name string, value any) string {
	if value == nil {
		return ""
	}

	// Handle sql.Null types
	switch v := value.(type) {
	case sql.NullBool:
		if !v.Valid {
			return "null"
		}
		return fmt.Sprintf("%t", v.Bool)
	case sql.NullString:
		if !v.Valid {
			return "null"
		}
		return v.String
	case sql.NullInt64:
		if !v.Valid {
			return "null"
		}
		return fmt.Sprintf("%d", v.Int64)
	case sql.NullFloat64:
		if !v.Valid {
			return "null"
		}
		return fmt.Sprintf("%f", v.Float64)
	}

	// Handle timestamps
	ts, ok := value.(int64)
	if ok {
		if strings.HasSuffix(name, "At") {
			return time.Unix(ts, 0).Format("2006-01-02 03:04 PM")
		}
		return fmt.Sprintf("%d", ts)
	}

	return fmt.Sprintf("%v", value) // Default case
}

package router

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
)


func BindJSON(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		return errors.New("content type is not application/json")
	}

	if r.Body == nil {
		return errors.New("request body is empty")
	}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		return fmt.Errorf("failed to decode JSON: %v", err)
	}

	return nil
}

func BindXML(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/xml") && !strings.Contains(contentType, "text/xml") {
		return errors.New("content type is not xml")
	}

	if r.Body == nil {
		return errors.New("request body is empty")
	}

	decoder := xml.NewDecoder(r.Body)
	if err := decoder.Decode(v); err != nil {
		return err
	}

	return nil
}
func BindForm(r *http.Request, v interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if !strings.Contains(contentType, "application/x-www-form-urlencoded") {
		return errors.New("content type is not application/x-www-form-urlencoded")
	}

	if err := r.ParseForm(); err != nil {
		return err
	}

	val := reflect.ValueOf(v).Elem()
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		formValue := r.FormValue(fieldType.Name)
		if formValue == "" {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			field.SetString(formValue)
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			intValue, err := strconv.ParseInt(formValue, 10, 64)
			if err != nil {
				return err
			}
			field.SetInt(intValue)
		case reflect.Float32, reflect.Float64:
			floatValue, err := strconv.ParseFloat(formValue, 64)
			if err != nil {
				return err
			}
			field.SetFloat(floatValue)
		case reflect.Bool:
			boolValue, err := strconv.ParseBool(formValue)
			if err != nil {
				return err
			}
			field.SetBool(boolValue)
		// Add other types as needed
		default:
			return errors.New("unsupported field type")
		}
	}

	return nil
}
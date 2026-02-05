package middleware

import (
	"mime/multipart"
	"reflect"

	"github.com/labstack/echo/v4"
)

func BodyValidationMiddleware(schemaFactory func() interface{}) echo.MiddlewareFunc {
	return ValidationMiddleware(schemaFactory, ValidationOptions{
		Source:     ValidationSourceBody,
		ContextKey: "validatedBody",
	})
}

// MultipartValidationMiddleware validates multipart/form-data requests with files
func MultipartValidationMiddleware(schemaFactory func() interface{}) echo.MiddlewareFunc {
	return ValidationMiddleware(schemaFactory, ValidationOptions{
		Source:     ValidationSourceMultipart,
		ContextKey: "validatedBody",
	})
}

// bindMultipartForm binds multipart form data to struct
func bindMultipartForm(c echo.Context, schema interface{}, form *multipart.Form) error {
	val := reflect.ValueOf(schema)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			continue
		}

		// Get form tag
		formTag := fieldType.Tag.Get("form")
		if formTag == "" {
			formTag = fieldType.Name
		}

		// Handle file fields
		if field.Type() == reflect.TypeOf(&multipart.FileHeader{}) {
			if files, ok := form.File[formTag]; ok && len(files) > 0 {
				field.Set(reflect.ValueOf(files[0]))
			} else {
				println("DEBUG: No file found for field:", fieldType.Name, "form tag:", formTag)
			}
		} else if field.Type() == reflect.TypeOf([]*multipart.FileHeader{}) {
			if files, ok := form.File[formTag]; ok {
				field.Set(reflect.ValueOf(files))
			}
		} else {
			// Handle regular form values
			if values, ok := form.Value[formTag]; ok && len(values) > 0 {
				// Set string field
				if field.Kind() == reflect.String {
					field.SetString(values[0])
				}
			}
		}
	}

	return nil
}

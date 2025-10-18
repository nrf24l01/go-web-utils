package middleware

import (
	"mime/multipart"
	"net/http"
	"reflect"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func BodyValidationMiddleware(schemaFactory func() interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			schema := schemaFactory()

			if err := c.Bind(schema); err != nil {
				// If the bind error comes from validation, return formatted validation errors
				if _, ok := err.(validator.ValidationErrors); ok {
					return c.JSON(http.StatusUnprocessableEntity, FormatValidationErrors(err))
				}

				// For other bind errors return 400 Bad Request
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
			}

			if err := c.Validate(schema); err != nil {
				// Return 422 for validation errors with formatted message
				return c.JSON(http.StatusUnprocessableEntity, FormatValidationErrors(err))
			}

			c.Set("validatedBody", schema)
			return next(c)
		}
	}
}

// MultipartValidationMiddleware validates multipart/form-data requests with files
func MultipartValidationMiddleware(schemaFactory func() interface{}) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			schema := schemaFactory()

			// Parse multipart form
			form, err := c.MultipartForm()
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid multipart form"})
			}

			// Bind form values and files to schema
			if err := bindMultipartForm(c, schema, form); err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
			}

			println("DEBUG: Before validation, schema:", schema)
			// Validate schema
			if err := c.Validate(schema); err != nil {
				println("DEBUG: Validation FAILED:", err.Error())
				return c.JSON(http.StatusUnprocessableEntity, FormatValidationErrors(err))
			}
			println("DEBUG: Validation PASSED")

			c.Set("validatedBody", schema)
			return next(c)
		}
	}
}

// bindMultipartForm binds multipart form data to struct
func bindMultipartForm(c echo.Context, schema interface{}, form *multipart.Form) error {
	val := reflect.ValueOf(schema)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	println("DEBUG: Binding multipart form, fields count:", val.NumField())
	println("DEBUG: Available files in form:", len(form.File))
	for key := range form.File {
		println("DEBUG: Form file key:", key)
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanSet() {
			println("DEBUG: Field", fieldType.Name, "cannot be set (unexported)")
			continue
		}

		// Get form tag
		formTag := fieldType.Tag.Get("form")
		if formTag == "" {
			formTag = fieldType.Name
		}

		println("DEBUG: Processing field:", fieldType.Name, "form tag:", formTag, "type:", field.Type())

		// Handle file fields
		if field.Type() == reflect.TypeOf(&multipart.FileHeader{}) {
			if files, ok := form.File[formTag]; ok && len(files) > 0 {
				println("DEBUG: Setting file for field:", fieldType.Name, "filename:", files[0].Filename)
				field.Set(reflect.ValueOf(files[0]))
			} else {
				println("DEBUG: No file found for field:", fieldType.Name, "form tag:", formTag)
			}
		} else if field.Type() == reflect.TypeOf([]*multipart.FileHeader{}) {
			if files, ok := form.File[formTag]; ok {
				println("DEBUG: Setting", len(files), "files for field:", fieldType.Name)
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

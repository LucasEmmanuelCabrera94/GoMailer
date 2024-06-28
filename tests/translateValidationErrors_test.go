package tests

import (
	"goMailer/internal/utils"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestTranslateValidationErrors(t *testing.T) {
	validate := validator.New()
	type TestStruct struct {
		RequiredField string `validate:"required"`
		EmailField    string `validate:"email"`
		NefieldField  string `validate:"nefield=OtherField"`
		OtherField    string
	}

	tests := []struct {
		name  string
		input TestStruct
		want  []string
	}{
		{
			name: "required field error",
			input: TestStruct{
				RequiredField: "",
				EmailField:    "test@example.com",
				NefieldField:  "value1",
				OtherField:    "value2",
			},
			want: []string{"The field 'RequiredField' is required."},
		},
		{
			name: "email field error",
			input: TestStruct{
				RequiredField: "value",
				EmailField:    "invalid-email",
				NefieldField:  "value1",
				OtherField:    "value2",
			},
			want: []string{"The field 'EmailField' must be a valid email address."},
		},
		{
			name: "nefield error",
			input: TestStruct{
				RequiredField: "value",
				EmailField:    "test@example.com",
				NefieldField:  "value1",
				OtherField:    "value1",
			},
			want: []string{"The field 'NefieldField' must not be the same as field 'OtherField'."},
		},
		{
			name: "multiple errors",
			input: TestStruct{
				RequiredField: "",
				EmailField:    "invalid-email",
				NefieldField:  "value1",
				OtherField:    "value1",
			},
			want: []string{
				"The field 'RequiredField' is required.",
				"The field 'EmailField' must be a valid email address.",
				"The field 'NefieldField' must not be the same as field 'OtherField'.",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.input)
			if err != nil {
				validationErrors := err.(validator.ValidationErrors)
				errorMessages := utils.TranslateValidationErrors(validationErrors)
				assert.Equal(t, tt.want, errorMessages)
			} else {
				t.Fatalf("Expected validation errors but got none")
			}
		})
	}
}

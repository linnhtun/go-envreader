package envreader

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"testing"
)

func TestReadEnv(t *testing.T) {
	// Clean up environment variables after each test to ensure isolation
	t.Cleanup(func() {
		os.Unsetenv("TEST_INT")
		os.Unsetenv("TEST_INT64")
		os.Unsetenv("TEST_STRING")
		os.Unsetenv("TEST_BOOL")
		os.Unsetenv("TEST_FLOAT")
		os.Unsetenv("TEST_EMPTY")
		os.Unsetenv("TEST_INVALID_INT")
		os.Unsetenv("TEST_INVALID_INT64")
		os.Unsetenv("TEST_INVALID_BOOL")
		os.Unsetenv("TEST_INVALID_FLOAT")
	})

	tests := []struct {
		name         string
		envKey       string
		envValue     string // Value to set for the environment variable
		setEnv       bool   // Whether to set the environment variable at all
		defaultValue interface{}
		expectedVal  interface{}
		expectedErr  error // Expected error, nil if no error is expected
		// Use a string for expectedErrString when comparing exact error messages
		expectedErrString string
	}{
		// --- int tests ---
		{
			name:         "Int_EnvExists_ValidValue",
			envKey:       "TEST_INT",
			envValue:     "123",
			setEnv:       true,
			defaultValue: 0,
			expectedVal:  123,
			expectedErr:  nil,
		},
		{
			name:         "Int_EnvNotExists",
			envKey:       "NON_EXISTENT_INT",
			setEnv:       false,
			defaultValue: 456,
			expectedVal:  456,
			expectedErr:  nil,
		},
		{
			name:         "Int_EnvExists_EmptyValue",
			envKey:       "TEST_EMPTY",
			envValue:     "",
			setEnv:       true,
			defaultValue: 789,
			expectedVal:  789,
			expectedErr:  nil,
		},
		{
			name:              "Int_EnvExists_InvalidValue",
			envKey:            "TEST_INVALID_INT",
			envValue:          "abc",
			setEnv:            true,
			defaultValue:      10,
			expectedVal:       10,                                                                            // Returns default on conversion error
			expectedErr:       strconv.ErrSyntax,                                                             // Expect the sentinel error
			expectedErrString: `failed to convert "abc" to int: strconv.Atoi: parsing "abc": invalid syntax`, // Exact string match for clarity/debugging
		},
		{
			name:         "Int_NegativeValue",
			envKey:       "TEST_INT",
			envValue:     "-50",
			setEnv:       true,
			defaultValue: 0,
			expectedVal:  -50,
			expectedErr:  nil,
		},
		// --- int64 tests ---
		{
			name:         "Int64_EnvExists_ValidValue",
			envKey:       "TEST_INT64",
			envValue:     "9876543210",
			setEnv:       true,
			defaultValue: int64(0),
			expectedVal:  int64(9876543210),
			expectedErr:  nil,
		},
		{
			name:         "Int64_EnvNotExists",
			envKey:       "NON_EXISTENT_INT64",
			setEnv:       false,
			defaultValue: int64(12345),
			expectedVal:  int64(12345),
			expectedErr:  nil,
		},
		{
			name:         "Int64_EnvExists_EmptyValue",
			envKey:       "TEST_EMPTY",
			envValue:     "",
			setEnv:       true,
			defaultValue: int64(54321),
			expectedVal:  int64(54321),
			expectedErr:  nil,
		},
		{
			name:              "Int64_EnvExists_InvalidValue",
			envKey:            "TEST_INVALID_INT64",
			envValue:          "xyz",
			setEnv:            true,
			defaultValue:      int64(99),
			expectedVal:       int64(99),         // Returns default on conversion error
			expectedErr:       strconv.ErrSyntax, // Expect sentinel error
			expectedErrString: `failed to convert "xyz" to int64: strconv.ParseInt: parsing "xyz": invalid syntax`,
		},
		{
			name:         "Int64_MaxPositiveValue",
			envKey:       "TEST_INT64",
			envValue:     "9223372036854775807", // Max int64
			setEnv:       true,
			defaultValue: int64(0),
			expectedVal:  int64(9223372036854775807),
			expectedErr:  nil,
		},
		// --- String tests ---
		{
			name:         "String_EnvExists_ValidValue",
			envKey:       "TEST_STRING",
			envValue:     "hello_world",
			setEnv:       true,
			defaultValue: "default",
			expectedVal:  "hello_world",
			expectedErr:  nil,
		},
		{
			name:         "String_EnvNotExists",
			envKey:       "NON_EXISTENT_STRING",
			setEnv:       false,
			defaultValue: "fallback",
			expectedVal:  "fallback",
			expectedErr:  nil,
		},
		{
			name:         "String_EnvExists_EmptyValue",
			envKey:       "TEST_EMPTY",
			envValue:     "",
			setEnv:       true,
			defaultValue: "empty_default",
			expectedVal:  "empty_default",
			expectedErr:  nil,
		},

		// --- Boolean tests ---
		{
			name:         "Bool_EnvExists_True",
			envKey:       "TEST_BOOL",
			envValue:     "true",
			setEnv:       true,
			defaultValue: false,
			expectedVal:  true,
			expectedErr:  nil,
		},
		{
			name:         "Bool_EnvExists_False",
			envKey:       "TEST_BOOL",
			envValue:     "false",
			setEnv:       true,
			defaultValue: true,
			expectedVal:  false,
			expectedErr:  nil,
		},
		{
			name:         "Bool_EnvNotExists",
			envKey:       "NON_EXISTENT_BOOL",
			setEnv:       false,
			defaultValue: true,
			expectedVal:  true,
			expectedErr:  nil,
		},
		{
			name:         "Bool_EnvExists_EmptyValue",
			envKey:       "TEST_EMPTY",
			envValue:     "",
			setEnv:       true,
			defaultValue: true,
			expectedVal:  true,
			expectedErr:  nil,
		},
		{
			name:              "Bool_EnvExists_InvalidValue",
			envKey:            "TEST_INVALID_BOOL",
			envValue:          "not_a_bool",
			setEnv:            true,
			defaultValue:      true,
			expectedVal:       true,              // Returns default on error
			expectedErr:       strconv.ErrSyntax, // Expect sentinel error
			expectedErrString: `failed to convert "not_a_bool" to bool: strconv.ParseBool: parsing "not_a_bool": invalid syntax`,
		},

		// --- Float64 tests ---
		{
			name:         "Float64_EnvExists_ValidValue",
			envKey:       "TEST_FLOAT",
			envValue:     "3.14",
			setEnv:       true,
			defaultValue: 0.0,
			expectedVal:  3.14,
			expectedErr:  nil,
		},
		{
			name:         "Float64_EnvNotExists",
			envKey:       "NON_EXISTENT_FLOAT",
			setEnv:       false,
			defaultValue: 1.23,
			expectedVal:  1.23,
			expectedErr:  nil,
		},
		{
			name:         "Float64_EnvExists_EmptyValue",
			envKey:       "TEST_EMPTY",
			envValue:     "",
			setEnv:       true,
			defaultValue: 9.99,
			expectedVal:  9.99,
			expectedErr:  nil,
		},
		{
			name:              "Float64_EnvExists_InvalidValue",
			envKey:            "TEST_INVALID_FLOAT",
			envValue:          "not_a_float",
			setEnv:            true,
			defaultValue:      5.5,
			expectedVal:       5.5,               // Returns default on error
			expectedErr:       strconv.ErrSyntax, // Expect sentinel error
			expectedErrString: `failed to convert "not_a_float" to float64: strconv.ParseFloat: parsing "not_a_float": invalid syntax`,
		},

		// --- Unsupported type test ---
		{
			name:              "UnsupportedType_Struct",
			envKey:            "TEST_STRUCT",
			envValue:          "{}",
			setEnv:            true,
			defaultValue:      struct{ Name string }{Name: "default"},
			expectedVal:       struct{ Name string }{Name: "default"},                                                          // Returns default
			expectedErr:       fmt.Errorf("unsupported type for environment variable conversion: %T", struct{ Name string }{}), // Exact type formatting
			expectedErrString: `unsupported type for environment variable conversion: struct { Name string }`,
		},
		{
			name:              "UnsupportedType_Slice",
			envKey:            "TEST_SLICE",
			envValue:          "1,2,3",
			setEnv:            true,
			defaultValue:      []int{10, 20},
			expectedVal:       []int{10, 20},                                                                   // Returns default
			expectedErr:       fmt.Errorf("unsupported type for environment variable conversion: %T", []int{}), // Exact type formatting
			expectedErrString: `unsupported type for environment variable conversion: []int`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set the environment variable if required for the test case
			if tt.setEnv {
				os.Setenv(tt.envKey, tt.envValue)
			} else {
				// Ensure it's not set if explicitly not needed
				os.Unsetenv(tt.envKey)
			}

			// Use a type switch to call ReadEnv with the correct generic type
			var actualVal interface{}
			var actualErr error

			switch def := tt.defaultValue.(type) {
			case int:
				actualVal, actualErr = ReadEnv[int](tt.envKey, def)
			case int64:
				actualVal, actualErr = ReadEnv[int64](tt.envKey, def)
			case string:
				actualVal, actualErr = ReadEnv[string](tt.envKey, def)
			case bool:
				actualVal, actualErr = ReadEnv[bool](tt.envKey, def)
			case float64:
				actualVal, actualErr = ReadEnv[float64](tt.envKey, def)
			case struct{ Name string }: // Specific case for our test struct
				actualVal, actualErr = ReadEnv[struct{ Name string }](tt.envKey, def)
			case []int: // Specific case for our test slice
				actualVal, actualErr = ReadEnv[[]int](tt.envKey, def)
			default:
				t.Fatalf("Test setup error: Unsupported defaultValue type for test case: %T", tt.defaultValue)
			}

			// --- Compare the actual result value with the expected result value ---
			// Use reflect.DeepEqual for uncomparable types like slices and structs
			// Otherwise, use direct comparison for comparable types
			if !reflect.DeepEqual(actualVal, tt.expectedVal) {
				t.Errorf("ReadEnv[%T](%q, %v) returned value %v; want %v", tt.defaultValue, tt.envKey, tt.defaultValue, actualVal, tt.expectedVal)
			}

			// --- Compare the actual error with the expected error ---
			if tt.expectedErr != nil {
				if actualErr == nil {
					t.Errorf("ReadEnv[%T](%q, %v) expected an error, but got nil", tt.defaultValue, tt.envKey, tt.defaultValue)
				} else {
					// Use errors.Is for wrapped errors (like strconv.ErrSyntax, strconv.ErrRange)
					if errors.Is(tt.expectedErr, strconv.ErrSyntax) || errors.Is(tt.expectedErr, strconv.ErrRange) {
						if !errors.Is(actualErr, tt.expectedErr) {
							t.Errorf("ReadEnv[%T](%q, %v) returned error %q; want it to contain underlying error %q", tt.defaultValue, tt.envKey, tt.defaultValue, actualErr, tt.expectedErr)
						}
					} else if tt.expectedErrString != "" { // For other specific errors, compare their string representation
						if actualErr.Error() != tt.expectedErrString {
							t.Errorf("ReadEnv[%T](%q, %v) returned error %q; want %q", tt.defaultValue, tt.envKey, tt.defaultValue, actualErr, tt.expectedErrString)
						}
					} else { // Fallback for unexpected non-nil expectedErr
						if actualErr.Error() != tt.expectedErr.Error() {
							t.Errorf("ReadEnv[%T](%q, %v) returned error %q; want %q", tt.defaultValue, tt.envKey, tt.defaultValue, actualErr, tt.expectedErr)
						}
					}
				}
			} else { // tt.expectedErr is nil
				if actualErr != nil {
					t.Errorf("ReadEnv[%T](%q, %v) returned unexpected error: %q", tt.defaultValue, tt.envKey, tt.defaultValue, actualErr)
				}
			}
		})
	}
}

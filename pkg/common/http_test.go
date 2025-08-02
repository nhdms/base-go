package common

import (
	"strings"
	"testing"
)

func TestExtractQueryParams(t *testing.T) {
	testCases := []struct {
		name           string
		queryParams    string
		target         string
		expectedResult string
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:           "Case 1: Target is the only parameter",
			queryParams:    "source_id=abc-123",
			target:         "source_id",
			expectedResult: "abc-123",
			expectError:    false,
		},
		{
			name:           "Case 2: Target is one of multiple parameters",
			queryParams:    "source_id=def-456&user_id=john_doe",
			target:         "source_id",
			expectedResult: "def-456",
			expectError:    false,
		},
		{
			name:           "Case 3: Target parameter not found",
			queryParams:    "user_id=jane_doe",
			target:         "source_id",
			expectedResult: "",
			expectError:    true,
			expectedErrMsg: "target source_id not found in query params",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actualResult, err := ExtractQueryParams(tc.queryParams, tc.target)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected an error but got none")
				}
				if tc.expectedErrMsg != "" && !strings.Contains(err.Error(), tc.expectedErrMsg) {
					t.Errorf("expected error message to contain '%s', but got '%s'", tc.expectedErrMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("did not expect an error but got: %v", err)
				}
			}

			if actualResult != tc.expectedResult {
				t.Errorf("expected result '%s', but got '%s'", tc.expectedResult, actualResult)
			}
		})
	}
}

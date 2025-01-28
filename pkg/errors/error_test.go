// Copyright 2025 SGNL.ai, Inc.

// nolint: lll
package customerror_test

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	framework "github.com/sgnl-ai/adapter-framework"
	api_adapter_v1 "github.com/sgnl-ai/adapter-framework/api/adapter/v1"
	customerror "github.com/sgnl-ai/adapter-sgnl/pkg/errors"
)

func TestUpdateError(t *testing.T) {
	errContextDeadlineExceeded := fmt.Errorf("timed out: %w", context.DeadlineExceeded)

	tests := map[string]struct {
		inputError     *framework.Error
		inputModifiers []customerror.ErrorModifier
		wantError      *framework.Error
	}{
		"success_request_exceeded_timeout": {
			inputError: &framework.Error{
				Message: fmt.Sprintf("Failed to execute PagerDuty request: %v.", errContextDeadlineExceeded),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
			inputModifiers: []customerror.ErrorModifier{customerror.WithRequestTimeoutMessage(errContextDeadlineExceeded, 30)},
			wantError: &framework.Error{
				Message: fmt.Sprintf("Failed to execute PagerDuty request: %v. Request exceeded configured timeout of 30 seconds. Please increase the request timeout.", errContextDeadlineExceeded),
				Code:    api_adapter_v1.ErrorCode_ERROR_CODE_INTERNAL,
			},
		},
		"success_nil_input": {
			inputError:     nil,
			inputModifiers: []customerror.ErrorModifier{customerror.WithRequestTimeoutMessage(errContextDeadlineExceeded, 30)},
			wantError:      nil,
		},
		"success_empty_struct": {
			inputError:     &framework.Error{},
			inputModifiers: []customerror.ErrorModifier{customerror.WithRequestTimeoutMessage(errContextDeadlineExceeded, 30)},
			wantError: &framework.Error{
				Message: "Request exceeded configured timeout of 30 seconds. Please increase the request timeout.",
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := customerror.UpdateError(tt.inputError, tt.inputModifiers...); !reflect.DeepEqual(got, tt.wantError) {
				t.Errorf("UpdateError() = %v, want %v", got, tt.wantError)
			}
		})
	}
}

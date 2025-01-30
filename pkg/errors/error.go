// Copyright 2025 SGNL.ai, Inc.
package customerror

import (
	"context"
	"errors"
	"fmt"

	framework "github.com/sgnl-ai/adapter-framework"
)

type ErrorModifier func(*framework.Error)

func WithRequestTimeoutMessage(reqErr error, timeout int) ErrorModifier {
	return func(frameworkErr *framework.Error) {
		if frameworkErr == nil {
			return
		}

		if errors.Is(reqErr, context.DeadlineExceeded) {
			msgToAppend := fmt.Sprintf(
				"Request exceeded configured timeout of %d seconds. Please increase the request timeout.",
				timeout,
			)

			if frameworkErr.Message == "" {
				frameworkErr.Message = msgToAppend

				return
			}

			frameworkErr.Message += " " + msgToAppend
		}
	}
}

func UpdateError(err *framework.Error, modifiers ...ErrorModifier) *framework.Error {
	for _, modifier := range modifiers {
		modifier(err)
	}

	return err
}

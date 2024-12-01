package apierror

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Kaese72/riskie-lib/logging"
)

type APIError struct {
	// Code indicates semantics based on HTTP status codes
	Code         int   `json:"code"`
	WrappedError error `json:"error"`
}

func (apierror APIError) MarshalJSON() ([]byte, error) {
	intermediary := struct {
		Code  int    `json:"code"`
		Error string `json:"error"`
	}{
		Code:  apierror.Code,
		Error: apierror.WrappedError.Error(),
	}
	bytes, err := json.Marshal(intermediary)
	if err != nil {
		return nil, err
	}

	return bytes, nil
}

func (apierror APIError) UnWrap() error {
	return apierror.WrappedError
}

func (apierror APIError) Error() string {
	return fmt.Sprintf("APIError: %s", apierror.WrappedError.Error())
}

func TerminalHTTPError(ctx context.Context, w http.ResponseWriter, err error) {
	var apiError APIError
	if errors.As(err, &apiError) {
		if apiError.Code == 500 {
			// When an unknown error occurs, do not send the error to the client
			http.Error(w, "Internal Server Error", apiError.Code)
			logging.Error(ctx, err.Error())
			return

		} else {
			bytes, intErr := json.MarshalIndent(apiError, "", "   ")
			if intErr != nil {
				// Must send a normal Error an not APIError just in case of eternal loop
				TerminalHTTPError(ctx, w, fmt.Errorf("error encoding response: %s", intErr.Error()))
				return
			}
			http.Error(w, string(bytes), apiError.Code)
			return
		}
	} else {
		TerminalHTTPError(ctx, w, APIError{Code: http.StatusInternalServerError, WrappedError: err})
		return
	}
}

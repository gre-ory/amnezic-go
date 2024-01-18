package api

import (
	"net/http"
)

func onPanic(resp http.ResponseWriter) func() {
	return func() {
		r := recover()
		if r != nil {
			if err, ok := r.(error); ok {
				encodeError(resp, http.StatusInternalServerError, err.Error())
			} else {
				encodeError(resp, http.StatusInternalServerError, "unknown error")
			}
		}
	}
}

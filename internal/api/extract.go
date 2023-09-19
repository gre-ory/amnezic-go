package api

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gre-ory/amnezic-go/internal/util"
	"github.com/julienschmidt/httprouter"
)

// //////////////////////////////////////////////////
// decode

func extractParameter(req *http.Request, name string) string {
	return strings.Trim(req.FormValue(name), " ")
}

func extractPathParameter(req *http.Request, name string) string {
	params := httprouter.ParamsFromContext(req.Context())
	return strings.Trim(params.ByName(name), " ")
}

// func toBool(value string) bool {
// 	value = strings.ToLower(value)
// 	return value == "true" || value == "1"
// }

func toStrings(values string) []string {
	return util.Convert(
		strings.Split(values, ","),
		func(value string) string { return strings.Trim(value, " ") },
	)
}

func toInt(value string) int {
	if result, err := strconv.Atoi(value); err == nil {
		return result
	}
	return 0
}

func toInt64(value string) int64 {
	if result, err := strconv.ParseInt(value, 10, 64); err == nil {
		return result
	}
	return 0
}

package api

import (
	"fmt"
	"net/http"
	"os"
)

// //////////////////////////////////////////////////
// game api

func NewRootHandler() http.Handler {
	return &rootHandler{}
}

type rootHandler struct {
}

func (h *rootHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	name := os.Getenv("NAME")
	if name == "" {
		name = "World"
	}
	fmt.Fprintf(resp, "Hello %s!\n", name)
}

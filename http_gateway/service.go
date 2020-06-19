package http_gateway

import (
	"io"
	"net/http"
	"strings"

	"github.com/TheSmallBoat/fiesta"
	"github.com/julienschmidt/httprouter"
)

func Handle(node *fiesta.Node, services []string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := make(map[string]string)
		for key := range r.Header {
			headers[strings.ToLower(key)] = r.Header.Get(key)
		}

		for key := range r.URL.Query() {
			headers["query."+strings.ToLower(key)] = r.URL.Query().Get(key)
		}

		params := httprouter.ParamsFromContext(r.Context())
		for _, param := range params {
			headers["params."+strings.ToLower(param.Key)] = param.Value
		}

		stream, err := node.StreamNode.Push(services, headers, r.Body)
		if err != nil {
			_, _ = w.Write([]byte(err.Error()))
			return
		}

		for name, val := range stream.Header.Headers {
			w.Header().Set(name, val)
		}

		_, _ = io.Copy(w, stream.Reader)
	})
}

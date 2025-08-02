package rocketentrypoint

import (
	"net/http"
)

func HandleFindRocketV1HTTP() http.HandlerFunc {
	return func(_ http.ResponseWriter, _ *http.Request) {}
}

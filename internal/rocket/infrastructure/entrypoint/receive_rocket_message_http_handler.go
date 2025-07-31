package rocketentrypoint

import (
	"net/http"
)

func HandleReceiveRocketMessageV1HTTP() http.HandlerFunc {
	return func(_ http.ResponseWriter, _ *http.Request) {}
}

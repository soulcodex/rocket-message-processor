package httpserver

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/soulcodex/rockets-message-processor/pkg/errutil"
)

const reverseProxyForwardedByHeader = "X-Forwarded-For"

var (
	ErrUnableToReadRequestBody      = errutil.NewError("unable to read request body")
	ErrUnableToUnmarshalRequestBody = errutil.NewError("unable to unmarshal request body")
)

func ClientIP(req *http.Request) string {
	ipAddress := req.RemoteAddr
	fwdAddress := req.Header.Get(reverseProxyForwardedByHeader)
	if fwdAddress != "" {
		ipAddress = fwdAddress

		ips := strings.Split(fwdAddress, ", ")
		if len(ips) > 1 {
			ipAddress = ips[0]
		}
	}

	return ipAddress
}

func CloneRequest(r *http.Request) *http.Request {
	var bodyBytes []byte
	newRequest := *r.WithContext(r.Context())

	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
	}

	r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	newRequest.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	return &newRequest
}

func AllParamsRequest(r *http.Request) (map[string]interface{}, error) {
	newRequest := CloneRequest(r)

	allParams, err := ConvertRequestToBodyMap(newRequest)
	if err != nil {
		return allParams, err
	}

	queryParams := QueryParams(newRequest)

	for k, v := range queryParams {
		allParams[k] = v
	}

	return allParams, nil
}

func ConvertRequestToBodyMap(r *http.Request) (map[string]interface{}, error) {
	requestBody := make(map[string]interface{})

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return requestBody, ErrUnableToReadRequestBody.Wrap(err)
	}

	if unmarshallErr := json.Unmarshal(b, &requestBody); unmarshallErr != nil {
		return requestBody, ErrUnableToUnmarshalRequestBody.Wrap(unmarshallErr)
	}

	return requestBody, nil
}

func QueryParams(r *http.Request) map[string]interface{} {
	query := make(map[string]interface{})

	for k, v := range r.URL.Query() {
		query[k] = v[0]
	}

	return query
}

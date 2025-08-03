package rocketentrypoint

import (
	"net/http"
	"strings"

	rocketqueries "github.com/soulcodex/rockets-message-processor/internal/rocket/application/queries"
	"github.com/soulcodex/rockets-message-processor/pkg/bus"
	querybus "github.com/soulcodex/rockets-message-processor/pkg/bus/query"
	httpserver "github.com/soulcodex/rockets-message-processor/pkg/http-server"
)

const (
	sortQueryParam = "sort"
	defaultSort    = "-created_at"
)

type SortParams struct {
	Sort string
	Asc  bool
}

func newSortParams(r *http.Request) SortParams {
	sort, asc := httpserver.FetchStringQueryParamValue(r.URL.Query(), sortQueryParam, defaultSort), true

	if strings.HasPrefix(sort, "-") {
		asc = false
		sort = strings.TrimPrefix(sort, "-")
	}

	return SortParams{
		Sort: sort,
		Asc:  asc,
	}
}

func HandleSearchRocketsV1HTTP(
	queryBus querybus.Bus,
	responseWriter *httpserver.JSONResponseWriter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		sortParams := newSortParams(r)

		searchQuery := &rocketqueries.SearchRocketsQuery{
			Sort: sortParams.Sort,
			Asc:  sortParams.Asc,
		}

		resp, err := bus.DispatchWithResponse[*rocketqueries.SearchRocketsQuery, rocketqueries.RocketsResponse](
			queryBus,
		)(r.Context(), searchQuery)

		switch err {
		case nil:
			response := newRocketsResponseV1(resp)
			responseWriter.WriteResponse(r.Context(), w, response, http.StatusOK)
		default:
			responseWriter.WriteErrorResponse(r.Context(), w, []string{err.Error()}, http.StatusInternalServerError)
		}
	}
}

package rocketentrypoint

import (
	"net/http"
	"time"

	"github.com/gorilla/mux"

	rocketqueries "github.com/soulcodex/rockets-message-processor/internal/rocket/application/queries"
	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
	"github.com/soulcodex/rockets-message-processor/pkg/bus"
	querybus "github.com/soulcodex/rockets-message-processor/pkg/bus/query"
	httpserver "github.com/soulcodex/rockets-message-processor/pkg/http-server"
	"github.com/soulcodex/rockets-message-processor/pkg/utils"
)

type FindRocketByIDResponse struct {
	ID          string    `json:"id"`
	RocketType  string    `json:"rocket_type"`
	LaunchSpeed int64     `json:"launch_speed"`
	Mission     string    `json:"mission"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func newFindRocketByIDResponseFromPrimitives(r rocketqueries.FindRocketByIDResponse) FindRocketByIDResponse {
	return FindRocketByIDResponse{
		ID:          r.ID,
		RocketType:  r.RocketType,
		LaunchSpeed: r.LaunchSpeed,
		Mission:     r.Mission,
		CreatedAt:   r.CreatedAt,
		UpdatedAt:   r.UpdatedAt,
	}
}

func HandleFindRocketV1HTTP(
	queryBus querybus.Bus,
	responseWriter *httpserver.JSONResponseWriter,
) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rocketID := mux.Vars(r)["rocket_id"]
		if rocketID == "" {
			responseWriter.WriteErrorResponse(r.Context(), w, []string{"rocket_id is required"}, http.StatusBadRequest)
			return
		}

		if err := utils.GuardUUID(rocketID); err != nil {
			responseWriter.WriteErrorResponse(r.Context(), w, []string{"invalid rocket_id format"}, http.StatusBadRequest)
			return
		}

		findQuery := &rocketqueries.FindRocketByIDQuery{RocketID: rocketID}

		resp, err := bus.DispatchWithResponse[*rocketqueries.FindRocketByIDQuery, rocketqueries.FindRocketByIDResponse](
			queryBus,
		)(r.Context(), findQuery)

		switch {
		case err == nil:
			response := newFindRocketByIDResponseFromPrimitives(resp)
			responseWriter.WriteResponse(r.Context(), w, response, http.StatusOK)
		case rocketdomain.IsRocketNotFoundError(err):
			responseWriter.WriteErrorResponse(r.Context(), w, []string{"rocket not found"}, http.StatusNotFound)
		default:
			responseWriter.WriteErrorResponse(r.Context(), w, []string{err.Error()}, http.StatusInternalServerError)
		}
	}
}

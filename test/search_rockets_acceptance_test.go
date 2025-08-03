package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/soulcodex/rockets-message-processor/cmd/di"
	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
	rocketentrypoint "github.com/soulcodex/rockets-message-processor/internal/rocket/infrastructure/entrypoint"
	rockettest "github.com/soulcodex/rockets-message-processor/test/rocket"
	testutils "github.com/soulcodex/rockets-message-processor/test/utils"
)

type SearchRocketsAcceptanceTestSuite struct {
	suite.Suite

	common       *di.CommonServices
	rocketModule *di.RocketModule

	rocketsByTimestamps rocketdomain.RocketCollection
	rocketsBySpeed      rocketdomain.RocketCollection
}

func TestSearchRockets(t *testing.T) {
	suite.Run(t, new(SearchRocketsAcceptanceTestSuite))
}

func (suite *SearchRocketsAcceptanceTestSuite) SetupSuite() {
	suite.common = di.MustInitCommonServicesWithEnvFiles(
		suite.T().Context(),
		"../.env",
		".test.env",
	)
	suite.rocketModule = di.NewRocketModule(suite.T().Context(), suite.common)
	suite.common.RedisClient.FlushAll(suite.T().Context())

	rocketOne := rockettest.NewRocketMother(
		rockettest.WithRocketID(suite.common.UUIDProvider.New().String()),
		rockettest.WithLaunchSpeed(3000),
		rockettest.WithUpdateDate(time.Now().Add(3*time.Hour)),
	).Build(suite.T())
	err := suite.rocketModule.Repository.Save(suite.T().Context(), rocketOne)
	suite.Require().NoError(err, "failed to save rocket one for suite setup")

	rocketTwo := rockettest.NewRocketMother(
		rockettest.WithRocketID(suite.common.UUIDProvider.New().String()),
		rockettest.WithLaunchSpeed(5000),
	).Build(suite.T())
	err = suite.rocketModule.Repository.Save(suite.T().Context(), rocketTwo)
	suite.Require().NoError(err, "failed to save rocket two for suite setup")

	rocketThree := rockettest.NewRocketMother(
		rockettest.WithRocketID(suite.common.UUIDProvider.New().String()),
		rockettest.WithLaunchSpeed(7000),
		rockettest.WithSoftDeletion(),
	).Build(suite.T())
	err = suite.rocketModule.Repository.Save(suite.T().Context(), rocketThree)
	suite.Require().NoError(err, "failed to save rocket three for suite setup")

	suite.rocketsByTimestamps = rocketdomain.NewRocketCollection(rocketOne, rocketTwo)
	suite.rocketsBySpeed = rocketdomain.NewRocketCollection(rocketTwo, rocketOne)
}

func (suite *SearchRocketsAcceptanceTestSuite) TestSearchRockets_SuccessByCreationDate() {
	const path = "/rockets?sort=-created_at"
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodGet, path, nil)
	suite.Equal(http.StatusOK, response.Code, "Expected status code 200 OK")
	rocketsResponse := suite.rocketResponse(response)
	suite.Len(rocketsResponse, 2, "Expected number of rockets to match")

	for idx, rocketResponse := range suite.rocketsByTimestamps.All() {
		suite.NotEmpty(rocketResponse.ID, "Expected rocket ID to be present")
		rocket := rocketsResponse[idx]
		suite.Equal(rocket.ID, rocketResponse.Primitives().ID, "Expected rocket ID to match")
		suite.Equal(rocket.RocketType, rocketResponse.Primitives().RocketType, "Expected rocket name to match")
		suite.Equal(rocket.Mission, rocketResponse.Primitives().Mission, "Expected rocket mission to match")
		suite.Equal(rocket.LaunchSpeed, rocketResponse.Primitives().LaunchSpeed, "Expected rocket launch speed to match")
	}
}

func (suite *SearchRocketsAcceptanceTestSuite) TestSearchRockets_SuccessByLaunchSpeed() {
	const path = "/rockets?sort=-launch_speed"
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodGet, path, nil)
	suite.Equal(http.StatusOK, response.Code, "Expected status code 200 OK")
	rocketsResponse := suite.rocketResponse(response)
	suite.Len(rocketsResponse, 2, "Expected number of rockets to match")

	for idx, rocketResponse := range suite.rocketsBySpeed.All() {
		suite.NotEmpty(rocketResponse.ID, "Expected rocket ID to be present")
		rocket := rocketsResponse[idx]
		suite.Equal(rocket.ID, rocketResponse.Primitives().ID, "Expected rocket ID to match")
		suite.Equal(rocket.RocketType, rocketResponse.Primitives().RocketType, "Expected rocket name to match")
		suite.Equal(rocket.Mission, rocketResponse.Primitives().Mission, "Expected rocket mission to match")
		suite.Equal(rocket.LaunchSpeed, rocketResponse.Primitives().LaunchSpeed, "Expected rocket launch speed to match")
	}
}

func (suite *SearchRocketsAcceptanceTestSuite) TestSearchRockets_SuccessUpdateDate() {
	const path = "/rockets?sort=-updated_at"
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodGet, path, nil)
	suite.Equal(http.StatusOK, response.Code, "Expected status code 200 OK")
	rocketsResponse := suite.rocketResponse(response)
	suite.Len(rocketsResponse, 2, "Expected number of rockets to match")

	for idx, rocketResponse := range suite.rocketsByTimestamps.All() {
		suite.NotEmpty(rocketResponse.ID, "Expected rocket ID to be present")
		rocket := rocketsResponse[idx]
		suite.Equal(rocket.ID, rocketResponse.Primitives().ID, "Expected rocket ID to match")
		suite.Equal(rocket.RocketType, rocketResponse.Primitives().RocketType, "Expected rocket name to match")
		suite.Equal(rocket.Mission, rocketResponse.Primitives().Mission, "Expected rocket mission to match")
		suite.Equal(rocket.LaunchSpeed, rocketResponse.Primitives().LaunchSpeed, "Expected rocket launch speed to match")
	}
}

func (suite *SearchRocketsAcceptanceTestSuite) rocketResponse(res *httptest.ResponseRecorder) rocketentrypoint.RocketsResponseV1 {
	suite.T().Helper()

	var rocketResponse rocketentrypoint.RocketsResponseV1
	err := json.Unmarshal(res.Body.Bytes(), &rocketResponse)
	suite.NoError(err, "failed to unmarshal rocket response")

	return rocketResponse
}

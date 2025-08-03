package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/soulcodex/rockets-message-processor/cmd/di"
	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
	rocketentrypoint "github.com/soulcodex/rockets-message-processor/internal/rocket/infrastructure/entrypoint"
	rockettest "github.com/soulcodex/rockets-message-processor/test/rocket"
	testutils "github.com/soulcodex/rockets-message-processor/test/utils"
)

type FindRocketByIDAcceptanceTestSuite struct {
	suite.Suite

	common       *di.CommonServices
	rocketModule *di.RocketModule

	rocketID rocketdomain.RocketID
}

func TestFindRocketByID(t *testing.T) {
	suite.Run(t, new(FindRocketByIDAcceptanceTestSuite))
}

func (suite *FindRocketByIDAcceptanceTestSuite) SetupSuite() {
	suite.common = di.MustInitCommonServicesWithEnvFiles(
		suite.T().Context(),
		"../.env",
		".test.env",
	)
	suite.rocketModule = di.NewRocketModule(suite.T().Context(), suite.common)
	suite.common.RedisClient.FlushAll(suite.T().Context())
	suite.rocketID = rocketdomain.RocketID(suite.common.UUIDProvider.New().String())

	rocket := rockettest.NewRocketMother(rockettest.WithRocketID(suite.rocketID.String())).Build(suite.T())
	err := suite.rocketModule.Repository.Save(suite.T().Context(), rocket)
	suite.Require().NoError(err, "failed to save rocket for suite setup")
}

func (suite *FindRocketByIDAcceptanceTestSuite) TestFindRocketByID_RocketNotFound() {
	rocketID := suite.common.UUIDProvider.New().String()
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodGet, "/rockets/"+rocketID, nil)
	suite.Equal(http.StatusNotFound, response.Code, "Expected status code 404 Not Found")
}

func (suite *FindRocketByIDAcceptanceTestSuite) TestFindRocketByID_InvalidRocketID() {
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodGet, "/rockets/1", nil)
	suite.Equal(http.StatusBadRequest, response.Code, "Expected status code 400 Bad Request")
}

func (suite *FindRocketByIDAcceptanceTestSuite) TestFindRocketByID_Success() {
	rocketID := suite.rocketID.String()
	response := testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, http.MethodGet, "/rockets/"+rocketID, nil)
	suite.Equal(http.StatusOK, response.Code, "Expected status code 200 OK")
	suite.assertRocket(response.Body.Bytes())
}

func (suite *FindRocketByIDAcceptanceTestSuite) assertRocket(body []byte) {
	suite.T().Helper()

	var rocketResponse rocketentrypoint.RocketResponseV1
	err := json.Unmarshal(body, &rocketResponse)
	suite.NoError(err, "failed to unmarshal rocket response")

	suite.Equal(suite.rocketID.String(), rocketResponse.ID, "Expected rocket ID to match")
	suite.NotEmpty(rocketResponse.RocketType, "Expected rocket type to be non-empty")
	suite.Greater(rocketResponse.LaunchSpeed, int64(0), "Expected launch speed to be greater than zero")
	suite.NotEmpty(rocketResponse.Mission, "Expected mission to be non-empty")
}

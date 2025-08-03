package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	"github.com/soulcodex/rockets-message-processor/cmd/di"
	rocketdomain "github.com/soulcodex/rockets-message-processor/internal/rocket/domain"
	rockettest "github.com/soulcodex/rockets-message-processor/test/rocket"
	testutils "github.com/soulcodex/rockets-message-processor/test/utils"
)

type RocketMessageReceiveSmokeTestSuite struct {
	suite.Suite

	common       *di.CommonServices
	rocketModule *di.RocketModule

	rocketID rocketdomain.RocketID
}

func TestRocketMessageReceive(t *testing.T) {
	suite.Run(t, new(RocketMessageReceiveSmokeTestSuite))
}

func (suite *RocketMessageReceiveSmokeTestSuite) SetupSuite() {
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

func (suite *RocketMessageReceiveSmokeTestSuite) TestReceiveRocketEvent_FailUnexpectedFormat() {
	eventBody := []byte(`
		{
			"message": {
				"type": "Falcon-9",
				"launchSpeed": 500,
				"mission": "ARTEMIS"
			}
		}
	`)
	response := suite.executeJSONRequest(http.MethodPost, "/messages", eventBody)
	suite.Equal(http.StatusBadRequest, response.Code, "Expected status code 200 OK")
}

func (suite *RocketMessageReceiveSmokeTestSuite) TestReceiveRocketLaunched_Success() {
	rocketID := suite.common.UUIDProvider.New().String()
	eventBody := suite.rocketLaunchedEventBody(
		rocketID,
		"Falcon-9",
		500,
		"ARTEMIS",
		time.Now(),
	)
	response := suite.executeJSONRequest(http.MethodPost, "/messages", eventBody)
	suite.Equal(http.StatusNoContent, response.Code, "Expected status code 200 OK")
}

func (suite *RocketMessageReceiveSmokeTestSuite) TestReceiveRocketExploded_Success() {
	eventBody := suite.rocketExplodedEventBody(suite.rocketID.String(), time.Now())
	response := suite.executeJSONRequest(http.MethodPost, "/messages", eventBody)
	suite.Equal(http.StatusNoContent, response.Code, "Expected status code 200 OK")
}

func (suite *RocketMessageReceiveSmokeTestSuite) TestReceiveRocketMissionChanged_Success() {
	eventBody := suite.rocketMissionChangedEventBody(suite.rocketID.String(), "LUNAR", time.Now())
	response := suite.executeJSONRequest(http.MethodPost, "/messages", eventBody)
	suite.Equal(http.StatusNoContent, response.Code, "Expected status code 200 OK")
}

func (suite *RocketMessageReceiveSmokeTestSuite) TestReceiveRocketSpeedIncreased_Success() {
	eventBody := suite.rocketSpeedIncreasedEventBody(suite.rocketID.String(), 500, time.Now())
	response := suite.executeJSONRequest(http.MethodPost, "/messages", eventBody)
	suite.Equal(http.StatusNoContent, response.Code, "Expected status code 200 OK")
}

func (suite *RocketMessageReceiveSmokeTestSuite) TestReceiveRocketSpeedDecreased_Success() {
	eventBody := suite.rocketSpeedDecreasedEventBody(suite.rocketID.String(), 500, time.Now())
	response := suite.executeJSONRequest(http.MethodPost, "/messages", eventBody)
	suite.Equal(http.StatusNoContent, response.Code, "Expected status code 200 OK")
}

func (suite *RocketMessageReceiveSmokeTestSuite) executeJSONRequest(
	verb,
	path string,
	body []byte,
) *httptest.ResponseRecorder {
	return testutils.ExecuteJSONRequest(suite.T(), suite.common.Router, verb, path, body)
}

func (suite *RocketMessageReceiveSmokeTestSuite) rocketLaunchedEventBody(
	rocketID string,
	rocketType string,
	launchSpeed int,
	mission string,
	at time.Time,
) []byte {
	body := fmt.Sprintf(`{"type": "%s","launchSpeed": %d,"mission": "%s"}`, rocketType, launchSpeed, mission)
	return suite.rocketEventBody(rocketID, 1, []byte(body), "RocketLaunched", at)
}

func (suite *RocketMessageReceiveSmokeTestSuite) rocketExplodedEventBody(
	rocketID string,
	at time.Time,
) []byte {
	body := `{"reason": "PRESSURE_VESSEL_FAILURE"}`
	return suite.rocketEventBody(rocketID, 6, []byte(body), "RocketExploded", at)
}

func (suite *RocketMessageReceiveSmokeTestSuite) rocketMissionChangedEventBody(
	rocketID string,
	newMission string,
	at time.Time,
) []byte {
	body := fmt.Sprintf(`{"newMission": "%s"}`, newMission)
	return suite.rocketEventBody(rocketID, 2, []byte(body), "RocketMissionChanged", at)
}

func (suite *RocketMessageReceiveSmokeTestSuite) rocketSpeedIncreasedEventBody(
	rocketID string,
	speed int,
	at time.Time,
) []byte {
	body := fmt.Sprintf(`{"by": %d}`, speed)
	return suite.rocketEventBody(rocketID, 3, []byte(body), "RocketSpeedIncreased", at)
}

func (suite *RocketMessageReceiveSmokeTestSuite) rocketSpeedDecreasedEventBody(
	rocketID string,
	speed int,
	at time.Time,
) []byte {
	body := fmt.Sprintf(`{"by": %d}`, speed)
	return suite.rocketEventBody(rocketID, 4, []byte(body), "RocketSpeedDecreased", at)
}

func (suite *RocketMessageReceiveSmokeTestSuite) rocketEventBody(
	messageID string,
	messageCount int,
	messageContent []byte,
	messageType string,
	at time.Time,
) []byte {
	body := fmt.Sprintf(`
		{
			"metadata": {
				"channel": "%s",
				"messageNumber": %d,
				"messageTime": "%s",
				"messageType": "%s"
			},
			"message": %s
		}
	`, messageID, messageCount, at.Format(time.RFC3339), messageType, messageContent)

	return []byte(body)
}

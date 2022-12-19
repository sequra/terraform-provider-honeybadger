package cli

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"net/http"
	"testing"
)

func TestGetTeamsNotProperlyAnswering(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/teams"

	expectedHoneybadgerResponse := HoneybadgerTeams{}
	expectedErrorResponse := errors.New(`status: 500, body: {"results":null,"links":{"self":"","prev":"","next":""}}`)
	expectedBody, _ := json.Marshal(expectedHoneybadgerResponse)
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusInternalServerError).
		JSON(expectedBody)

	actualHoneybadgerResponse, errResponse := honeybadgerCli.GetTeams()

	assert.Equal(expectedHoneybadgerResponse.Teams, actualHoneybadgerResponse, "Actual response is different from expected response")
	assert.Equal(errResponse, expectedErrorResponse, "Reponse error must be 500")
}

func TestGetTeams(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/teams"

	expectedHoneybadgerResponse := HoneybadgerTeams{
		Teams: []HoneybadgerTeam{
			{
				ID:   1234,
				Name: "Test Sequra Team",
				Owner: HoneybadgerTeamOwner{
					ID:    945,
					Email: "test.sequra@sequra.es",
				},
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedHoneybadgerResponse)
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusOK).
		JSON(expectedBody)

	actualHoneybadgerResponse, errResponse := honeybadgerCli.GetTeams()

	assert.Equal(expectedHoneybadgerResponse.Teams, actualHoneybadgerResponse, "Actual response is different from expected response")
	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TestFindTeamByName(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/teams"

	expectedHoneybadgerResponse := HoneybadgerTeams{
		Teams: []HoneybadgerTeam{
			{
				ID:   1234,
				Name: "Test Sequra Team",
				Owner: HoneybadgerTeamOwner{
					ID:    945,
					Email: "test.sequra@sequra.es",
				},
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedHoneybadgerResponse)
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusOK).
		JSON(expectedBody)

	actualHoneybadgerResponse, errResponse := honeybadgerCli.FindTeamByName("Test Sequra Team")

	assert.Equal(expectedHoneybadgerResponse.Teams[0], actualHoneybadgerResponse, "Actual response is different from expected response")
	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TestFindTeamByNameNotFound(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/teams"

	expectedHoneybadgerResponse := HoneybadgerTeams{
		Teams: []HoneybadgerTeam{
			{
				ID:   1234,
				Name: "Test Sequra Team",
				Owner: HoneybadgerTeamOwner{
					ID:    945,
					Email: "test.sequra@sequra.es",
				},
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedHoneybadgerResponse)
	expectedErrorResponse := errors.New("Team not found")
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusOK).
		JSON(expectedBody)

	actualHoneybadgerResponse, errResponse := honeybadgerCli.FindTeamByName("Test Team Not Found")

	assert.Equal(HoneybadgerTeam{}, actualHoneybadgerResponse, "Actual response is different from expected response")
	assert.Equal(errResponse, expectedErrorResponse, "Reponse error is different from expected")
}

func TestFindTeamByID(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/teams"

	expectedHoneybadgerResponse := HoneybadgerTeams{
		Teams: []HoneybadgerTeam{
			{
				ID:   1234,
				Name: "Test Sequra Team",
				Owner: HoneybadgerTeamOwner{
					ID:    945,
					Email: "test.sequra@sequra.es",
				},
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedHoneybadgerResponse)
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusOK).
		JSON(expectedBody)

	actualHoneybadgerResponse, errResponse := honeybadgerCli.FindTeamByID(1234)

	assert.Equal(expectedHoneybadgerResponse.Teams[0], actualHoneybadgerResponse, "Actual response is different from expected response")
	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TestFindTeamByIDNotFound(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/teams"

	expectedHoneybadgerResponse := HoneybadgerTeams{
		Teams: []HoneybadgerTeam{
			{
				ID:   1234,
				Name: "Test Sequra Team",
				Owner: HoneybadgerTeamOwner{
					ID:    945,
					Email: "test.sequra@sequra.es",
				},
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedHoneybadgerResponse)
	expectedErrorResponse := errors.New("Team not found")
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusOK).
		JSON(expectedBody)

	actualHoneybadgerResponse, errResponse := honeybadgerCli.FindTeamByID(0)

	assert.Equal(HoneybadgerTeam{}, actualHoneybadgerResponse, "Actual response is different from expected response")
	assert.Equal(errResponse, expectedErrorResponse, "Reponse error must be nil")
}

func TestCreateTeam(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/teams"

	expectedBody, _ := json.Marshal(nil)
	gock.New(honeybadgerAPIHost).
		Post(urlPath).
		Reply(http.StatusCreated).
		JSON(expectedBody)

	_, errResponse := honeybadgerCli.CreateTeam("New Team")

	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TestUpdateTeam(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	teamID := 999
	teamName := "Updated test team"
	urlPath := fmt.Sprintf("/v2/teams/%d", teamID)

	expectedBody, _ := json.Marshal(nil)
	gock.New(honeybadgerAPIHost).
		Put(urlPath).
		Reply(http.StatusNoContent).
		JSON(expectedBody)

	errResponse := honeybadgerCli.UpdateTeam(teamName, teamID)

	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TestDeleteTeam(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	teamID := 999
	urlPath := fmt.Sprintf("/v2/teams/%d", teamID)

	expectedBody, _ := json.Marshal(nil)
	gock.New(honeybadgerAPIHost).
		Delete(urlPath).
		Reply(http.StatusNoContent).
		JSON(expectedBody)

	errResponse := honeybadgerCli.DeleteTeam(teamID)

	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

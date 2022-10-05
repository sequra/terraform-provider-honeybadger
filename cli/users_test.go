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

var honeybadgerAPIHost = "http://localhost"
var honeybadgerAPIKey = "213123"
var honeybadgerTeamID = 23434
var honeybadgerCli = NewClient(&honeybadgerAPIHost, &honeybadgerAPIKey)

func TestGetUsersIsNotProperlyAnswering(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := fmt.Sprintf("/v2/teams/%d/team_members", honeybadgerTeamID)

	expectedHoneybadgerResponse := HoneybadgerUsers{}
	expectedErrorResponse := errors.New(`status: 500, body: {"results":null,"links":{"self":"","prev":"","next":""}}`)

	expectedBody, _ := json.Marshal(expectedHoneybadgerResponse)
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusInternalServerError).
		JSON(expectedBody)

	actualHoneybadgerResponse, errResponse := honeybadgerCli.GetUsers(honeybadgerTeamID)

	assert.Equal(expectedHoneybadgerResponse.Users, actualHoneybadgerResponse, "Actual response is different from expected response")
	assert.Equal(errResponse, expectedErrorResponse, "Reponse error must be 500")
}

func TestGetUsersWithoutPagination(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := fmt.Sprintf("/v2/teams/%d/team_members", honeybadgerTeamID)

	expectedHoneybadgerResponse := HoneybadgerUsers{
		Users: []HoneybadgerUser{
			{
				ID:      10,
				Name:    "Test Sequra Page2",
				Email:   "test.sequra.page2@sequra.es",
				IsAdmin: false,
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedHoneybadgerResponse)
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusOK).
		JSON(expectedBody)

	actualHoneybadgerResponse, errResponse := honeybadgerCli.GetUsers(honeybadgerTeamID)

	assert.Equal(expectedHoneybadgerResponse.Users, actualHoneybadgerResponse, "Actual response is different from expected response")
	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TestGetUsersWithPagination(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)

	var hbExpectedUserList []HoneybadgerUser

	expectedPaginatedResponse := []struct {
		hbUsers          HoneybadgerUsers
		urlPath          string
		httpResponseCode int
	}{
		{
			httpResponseCode: http.StatusOK,
			urlPath:          fmt.Sprintf("/v2/teams/%d/team_members", honeybadgerTeamID),
			hbUsers: HoneybadgerUsers{
				Users: []HoneybadgerUser{
					{
						ID:      10,
						Name:    "Test Sequra Page1",
						Email:   "test.sequra.page1@sequra.es",
						IsAdmin: false,
					},
				},
				Links: HoneybadgerLink{
					Self:     "http://localhost",
					NextPage: "/page2",
				},
			},
		},
		{
			httpResponseCode: http.StatusOK,
			urlPath:          "/page2",
			hbUsers: HoneybadgerUsers{
				Users: []HoneybadgerUser{
					{
						ID:      10,
						Name:    "Test Sequra Page2",
						Email:   "test.sequra.page2@sequra.es",
						IsAdmin: false,
					},
				},
				Links: HoneybadgerLink{
					Self:         "http://localhost",
					PreviousPage: "/v2/teams/23434/team_members",
				},
			},
		},
	}

	for _, expectedResponse := range expectedPaginatedResponse {
		expectedBodyPage, _ := json.Marshal(expectedResponse.hbUsers)
		gock.New(honeybadgerAPIHost).
			Get(expectedResponse.urlPath).
			Reply(expectedResponse.httpResponseCode).
			JSON(expectedBodyPage)

		hbExpectedUserList = append(hbExpectedUserList, expectedResponse.hbUsers.Users...)
	}

	actualHoneybadgerResponse, errResponse := honeybadgerCli.GetUsers(honeybadgerTeamID)

	assert.Equal(hbExpectedUserList, actualHoneybadgerResponse, "Actual response is different from expected response")
	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TestCreateUser(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := fmt.Sprintf("/v2/teams/%d/team_invitations", honeybadgerTeamID)

	expectedBody, _ := json.Marshal(nil)
	gock.New(honeybadgerAPIHost).
		Post(urlPath).
		Reply(http.StatusCreated).
		JSON(expectedBody)

	errResponse := honeybadgerCli.CreateUser("new.user@sequra.es", false, honeybadgerTeamID)

	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TestUpdateUser(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	userID := 999
	isAdmin := false
	urlPath := fmt.Sprintf("/v2/teams/%d/team_members/%d", honeybadgerTeamID, userID)

	expectedBody, _ := json.Marshal(nil)
	gock.New(honeybadgerAPIHost).
		Put(urlPath).
		Reply(http.StatusNoContent).
		JSON(expectedBody)

	errResponse := honeybadgerCli.UpdateUser(userID, isAdmin, honeybadgerTeamID)

	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TestDeleteUser(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	userID := 999
	urlPath := fmt.Sprintf("/v2/teams/%d/team_members/%d", honeybadgerTeamID, userID)

	expectedBody, _ := json.Marshal(nil)
	gock.New(honeybadgerAPIHost).
		Delete(urlPath).
		Reply(http.StatusNoContent).
		JSON(expectedBody)

	errResponse := honeybadgerCli.DeleteUser(userID, honeybadgerTeamID)

	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TesGetUserFromTeams(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	userID := 999
	teamID := 999

	expectedHoneybadgerTeamResponse := HoneybadgerTeams{
		Teams: []HoneybadgerTeam{
			{
				ID:   teamID,
				Name: "Test Sequra Team",
				Owner: HoneybadgerTeamOwner{
					ID:    945,
					Email: "test.sequra@sequra.es",
				},
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedHoneybadgerTeamResponse)
	gock.New(honeybadgerAPIHost).
		Get("/v2/teams").
		Reply(http.StatusOK).
		JSON(expectedBody)

	expectedHoneybadgerTeamUserResponse := HoneybadgerUsers{
		Users: []HoneybadgerUser{
			{
				ID:      userID,
				Name:    "Test Sequra Page2",
				Email:   "test.sequra.page2@sequra.es",
				IsAdmin: false,
			},
		},
	}
	expectedBody, _ = json.Marshal(expectedHoneybadgerTeamUserResponse)
	gock.New(honeybadgerAPIHost).
		Get(fmt.Sprintf("/v2/teams/%d/team_members", teamID)).
		Reply(http.StatusNoContent).
		JSON(expectedBody)

	actualHoneybadgerResponse, actualErrResponse := honeybadgerCli.GetUserFromTeams("test.sequra.page2@sequra.es")

	assert.Equal(expectedHoneybadgerTeamUserResponse, actualHoneybadgerResponse, "Actual response is different from expected response")
	assert.Equal(actualErrResponse, nil, "Reponse error does not match")
}

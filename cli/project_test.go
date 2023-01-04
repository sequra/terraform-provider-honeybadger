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

func TestGetProjectsNotProperlyAnswering(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/projects"

	expectedResponse := HoneybadgerProjects{}
	expectedErrorResponse := errors.New(`status: 500, body: {"results":null,"links":{"self":"","prev":"","next":""}}`)
	expectedBody, _ := json.Marshal(expectedResponse)
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusInternalServerError).
		JSON(expectedBody)

	actualResponse, actualErrResponse := honeybadgerCli.GetProjects()

	assert.Equal(expectedResponse.Projects, actualResponse, "Actual response is different from expected response")
	assert.Equal(actualErrResponse, expectedErrorResponse, "Reponse error must be 500")
}

func TestGetProjects(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/projects"

	expectedResponse := HoneybadgerProjects{
		Projects: []HoneybadgerProject{
			{
				ID:           1234,
				Name:         "Test Sequra Project",
				Environments: []string{"development", "production"},
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedResponse)
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusOK).
		JSON(expectedBody)

	actualResponse, actualErrResponse := honeybadgerCli.GetProjects()

	assert.Equal(expectedResponse.Projects, actualResponse, "Actual response is different from expected response")
	assert.Equal(actualErrResponse, nil, "Reponse error must be nil")
}

func TestFindProjectByNameProjects(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/projects"

	expectedResponse := HoneybadgerProjects{
		Projects: []HoneybadgerProject{
			{
				ID:           1234,
				Name:         "Test Sequra Project",
				Environments: []string{"development", "production"},
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedResponse)
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusOK).
		JSON(expectedBody)

	actualResponse, actualErrResponse := honeybadgerCli.FindProjectByName("Test Sequra Project")

	assert.Equal(expectedResponse.Projects[0], actualResponse, "Actual response is different from expected response")
	assert.Equal(actualErrResponse, nil, "Reponse error must be nil")
}

func TestFindProjectByNameProjectNotFound(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/projects"

	expectedResponse := HoneybadgerProjects{
		Projects: []HoneybadgerProject{
			{
				ID:           1234,
				Name:         "Test Sequra Project",
				Environments: []string{"development", "production"},
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedResponse)
	expectedErrorResponse := errors.New("Project not found")
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusOK).
		JSON(expectedBody)

	actualResponse, actualErrResponse := honeybadgerCli.FindProjectByName("Test Sequra Project1")

	assert.Equal(HoneybadgerProject{}, actualResponse, "Actual response is different from expected response")
	assert.Equal(actualErrResponse, expectedErrorResponse, "Reponse error must be nil")
}

func TestFindProjectByIDProjects(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/projects"

	expectedResponse := HoneybadgerProjects{
		Projects: []HoneybadgerProject{
			{
				ID:           1234,
				Name:         "Test Sequra Project",
				Environments: []string{"development", "production"},
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedResponse)
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusOK).
		JSON(expectedBody)

	actualResponse, actualErrResponse := honeybadgerCli.FindProjectByID(1234)

	assert.Equal(expectedResponse.Projects[0], actualResponse, "Actual response is different from expected response")
	assert.Equal(actualErrResponse, nil, "Reponse error must be nil")
}

func TestFindProjectByIDProjectsNotFound(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/projects"

	expectedResponse := HoneybadgerProjects{
		Projects: []HoneybadgerProject{
			{
				ID:           1234,
				Name:         "Test Sequra Project",
				Environments: []string{"development", "production"},
			},
		},
	}
	expectedBody, _ := json.Marshal(expectedResponse)
	expectedErrorResponse := errors.New("Project not found")
	gock.New(honeybadgerAPIHost).
		Get(urlPath).
		Reply(http.StatusOK).
		JSON(expectedBody)

	actualResponse, actualErrResponse := honeybadgerCli.FindProjectByID(0)

	assert.Equal(HoneybadgerProject{}, actualResponse, "Actual response is different from expected response")
	assert.Equal(actualErrResponse, expectedErrorResponse, "Expected error is different from actual error")
}

func TestCreateProject(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	urlPath := "/v2/projects"

	expectedBody, _ := json.Marshal(nil)
	gock.New(honeybadgerAPIHost).
		Post(urlPath).
		Reply(http.StatusCreated).
		JSON(expectedBody)

	_, errResponse := honeybadgerCli.CreateProject("New Project", "ruby")

	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TestUpdateProject(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	projectID := 1234
	urlPath := fmt.Sprintf("/v2/projects/%d", projectID)

	expectedBody, _ := json.Marshal(nil)
	gock.New(honeybadgerAPIHost).
		Put(urlPath).
		Reply(http.StatusCreated).
		JSON(expectedBody)

	errResponse := honeybadgerCli.UpdateProject("New Project", projectID, "ruby")

	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

func TestDeleteProject(t *testing.T) {
	defer gock.Off()
	assert := assert.New(t)
	projectID := 1234
	urlPath := fmt.Sprintf("/v2/projects/%d", projectID)

	expectedBody, _ := json.Marshal(nil)
	gock.New(honeybadgerAPIHost).
		Delete(urlPath).
		Reply(http.StatusCreated).
		JSON(expectedBody)

	errResponse := honeybadgerCli.DeleteProject(projectID)

	assert.Equal(errResponse, nil, "Reponse error must be nil")
}

package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// GetTeams - Get Honeybadger Teams
func (hbc *HoneybadgerClient) GetProjects() ([]HoneybadgerProject, error) {
	var hbProjects HoneybadgerProjects

	url := fmt.Sprintf("%s/v2/projects", hbc.HostURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return hbProjects.Projects, err
	}

	body, err := hbc.DoRequest(req)
	if err != nil {
		return hbProjects.Projects, err
	}

	err = json.Unmarshal(body, &hbProjects)
	if err != nil {
		return hbProjects.Projects, err
	}

	return hbProjects.Projects, nil
}

// FindProjectByName - Find Project by name
func (hbc *HoneybadgerClient) FindProjectByName(projectName string) (HoneybadgerProject, error) {
	hbProjects, err := hbc.GetProjects()
	if err != nil {
		return HoneybadgerProject{}, err
	}

	for _, project := range hbProjects {
		if project.Name == projectName {
			return project, nil
		}
	}
	return HoneybadgerProject{}, errors.New("Project not found")
}

// FindProjectByID - Find Project by ID
func (hbc *HoneybadgerClient) FindProjectByID(projectID int) (HoneybadgerProject, error) {
	hbProjects, err := hbc.GetProjects()
	if err != nil {
		return HoneybadgerProject{}, err
	}

	for _, project := range hbProjects {
		if project.ID == projectID {
			return project, nil
		}
	}
	return HoneybadgerProject{}, errors.New("Project not found")
}

// CreateProject - Create Project
func (hbc *HoneybadgerClient) CreateProject(projectName string, language string) (HoneybadgerProject, error) {
	var hbProject HoneybadgerProject
	var jsonPayload = []byte(`{"project":{"name":"` + projectName + `"}}`)

	if len(language) > 0 {
		jsonPayload = []byte(`{"project":{"name":"` + projectName + `", "language": "` + language + `"}}`)
	}

	url := fmt.Sprintf("%s/v2/projects", hbc.HostURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return HoneybadgerProject{}, err
	}

	body, err := hbc.DoRequest(req)
	if err != nil {
		return HoneybadgerProject{}, err
	}

	err = json.Unmarshal(body, &hbProject)
	if err != nil {
		return HoneybadgerProject{}, err
	}

	return hbProject, nil
}

// UpdateProject - Update Project
func (hbc *HoneybadgerClient) UpdateProject(projectName string, projectID int, language string) error {
	var jsonPayload = []byte(`{"project":{"name":"` + projectName + `"}}`)

	if len(language) > 0 {
		jsonPayload = []byte(`{"project":{"name":"` + projectName + `", "language": "` + language + `"}}`)
	}

	url := fmt.Sprintf("%s/v2/projects/%d", hbc.HostURL, projectID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	_, err = hbc.DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

// DeleteProject - Delete Project
func (hbc *HoneybadgerClient) DeleteProject(projectID int) error {
	url := fmt.Sprintf("%s/v2/projects/%d", hbc.HostURL, projectID)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	_, err = hbc.DoRequest(req)
	if err != nil {
		return err
	}

	return nil
}

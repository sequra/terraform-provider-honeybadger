package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// GetTeams - Get Honeybadger Teams
func (hbc *HoneybadgerClient) GetTeams() ([]HoneybadgerTeam, error) {
	var hbTeams HoneybadgerTeams

	url := fmt.Sprintf("%s/v2/teams", hbc.HostURL)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return hbTeams.Teams, err
	}

	body, err := hbc.DoRequest(req)
	if err != nil {
		return hbTeams.Teams, err
	}

	err = json.Unmarshal(body, &hbTeams)
	if err != nil {
		return hbTeams.Teams, err
	}

	return hbTeams.Teams, nil
}

// FindTeamByName - Find Team by name
func (hbc *HoneybadgerClient) FindTeamByName(teamName string) (HoneybadgerTeam, error) {
	hbTeams, err := hbc.GetTeams()
	if err != nil {
		return HoneybadgerTeam{}, err
	}

	for _, team := range hbTeams {
		if team.Name == teamName {
			return team, nil
		}
	}
	return HoneybadgerTeam{}, errors.New("Team not found")
}

// FindTeamByID - Find Team by ID
func (hbc *HoneybadgerClient) FindTeamByID(teamID int) (HoneybadgerTeam, error) {
	hbTeams, err := hbc.GetTeams()
	if err != nil {
		return HoneybadgerTeam{}, err
	}

	for _, team := range hbTeams {
		if team.ID == teamID {
			return team, nil
		}
	}
	return HoneybadgerTeam{}, errors.New("Team not found")
}

// CreateTeam - Create Team
func (hbc *HoneybadgerClient) CreateTeam(teamName string) (HoneybadgerTeam, error) {
	var hbTeam HoneybadgerTeam
	var jsonPayload = []byte(`{"team":{"name":"` + teamName + `"}}`)

	url := fmt.Sprintf("%s/v2/teams", hbc.HostURL)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return HoneybadgerTeam{}, err
	}

	body, err := hbc.DoRequest(req)
	if err != nil {
		return HoneybadgerTeam{}, err
	}

	err = json.Unmarshal(body, &hbTeam)
	if err != nil {
		return HoneybadgerTeam{}, err
	}

	return hbTeam, nil
}

// UpdateTeam - Update Team
func (hbc *HoneybadgerClient) UpdateTeam(teamName string, teamID int) error {
	var jsonPayload = []byte(`{"team":{"name":"` + teamName + `"}}`)

	url := fmt.Sprintf("%s/v2/teams/%d", hbc.HostURL, teamID)
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

// DeleteTeam - Delete Team
func (hbc *HoneybadgerClient) DeleteTeam(teamID int) error {
	url := fmt.Sprintf("%s/v2/teams/%d", hbc.HostURL, teamID)
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

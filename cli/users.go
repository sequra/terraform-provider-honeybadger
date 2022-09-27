package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// GetUsersPaginated - Returns all registered users in Honeybadger using pagination
func (hbc *HoneybadgerClient) GetUsersPaginated(pagePath string, hbUserList []HoneybadgerUser) ([]HoneybadgerUser, error) {
	var hbUsers HoneybadgerUsers
	urlPath := fmt.Sprintf("%s/%s", hbc.HostURL, pagePath)

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return hbUsers.Users, err
	}

	body, err := hbc.DoRequest(req)
	if err != nil {
		return hbUsers.Users, err
	}

	err = json.Unmarshal(body, &hbUsers)
	if err != nil {
		return hbUsers.Users, err
	}

	if hbUsers.Links.NextPage != "" {
		return hbc.GetUsersPaginated(hbUsers.Links.NextPage, hbUsers.Users)
	}

	hbUserList = append(hbUserList, hbUsers.Users...)

	return hbUserList, nil
}

// GetUsers - Returns all registered users in Honeybadger
func (hbc *HoneybadgerClient) GetUsers(teamID int) ([]HoneybadgerUser, error) {
	var hbUsers HoneybadgerUsers
	urlPath := fmt.Sprintf("/v2/teams/%d/team_members", teamID)

	return hbc.GetUsersPaginated(urlPath, hbUsers.Users)
}

// FindUserByID - Returns a user by ID
func (hbc *HoneybadgerClient) FindUserByID(userID int, teamID int) (HoneybadgerUser, error) {
	hbUsers, err := hbc.GetUsers(teamID)
	if err != nil {
		return HoneybadgerUser{}, err
	}

	for _, user := range hbUsers {
		if user.ID == userID {
			return user, nil
		}
	}
	return HoneybadgerUser{}, errors.New("User not found")
}

// FindUserByEmail - Returns a user by Email
func (hbc *HoneybadgerClient) FindUserByEmail(userEmail string, teamID int) (HoneybadgerUser, error) {
	hbUsers, err := hbc.GetUsers(teamID)
	if err != nil {
		return HoneybadgerUser{}, err
	}

	for _, user := range hbUsers {
		if user.Email == userEmail {
			return user, nil
		}
	}
	return HoneybadgerUser{}, errors.New("User not found")
}

// CreateUser - Create Honeybadger User
func (hbc *HoneybadgerClient) CreateUser(userEmail string, teamID int) error {
	var hbUser HoneybadgerUser
	var jsonPayload = []byte(`{"team_invitation":{"email":"` + userEmail + `"}}`)

	url := fmt.Sprintf("%s/v2/teams/%d/team_invitations", hbc.HostURL, teamID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	body, err := hbc.DoRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &hbUser)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser - Update Honeybadger User Information
func (hbc *HoneybadgerClient) UpdateUser(userID int, isAdmin bool, teamID int) error {
	var jsonPayload = []byte(`{"team_member":{"admin":` + strconv.FormatBool(isAdmin) + `}}`)

	url := fmt.Sprintf("%s/v2/teams/%d/team_members/%d", hbc.HostURL, teamID, userID)
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

// DeleteUser - Delete Honeybadger User
func (hbc *HoneybadgerClient) DeleteUser(userID int, teamID int) error {
	url := fmt.Sprintf("%s/v2/teams/%d/team_members/%d", hbc.HostURL, teamID, userID)
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

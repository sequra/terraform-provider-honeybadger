package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

// GetUsersPaginated - Returns all registered users in HoneyBadger using pagination
func (hbc *HoneyBadgerClient) GetUsersPaginated(pagePath string, hbUserList []HoneyBadgerUser) ([]HoneyBadgerUser, error) {
	var hbUsers HoneyBadgerUsers
	urlPath := fmt.Sprintf("%s/%s", hbc.HostURL, pagePath)

	req, err := http.NewRequest("GET", urlPath, nil)
	if err != nil {
		return hbUsers.Users, err
	}

	body, err := hbc.doRequest(req)
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

// GetUsers - Returns all registered users in HoneyBadger
func (hbc *HoneyBadgerClient) GetUsers() ([]HoneyBadgerUser, error) {
	var hbUsers HoneyBadgerUsers
	urlPath := fmt.Sprintf("/v2/teams/%d/team_members", hbc.TeamID)

	return hbc.GetUsersPaginated(urlPath, hbUsers.Users)
}

// FindUserByID - Returns a user by ID
func (hbc *HoneyBadgerClient) FindUserByID(userID int) (HoneyBadgerUser, error) {
	hbUsers, err := hbc.GetUsers()
	if err != nil {
		return HoneyBadgerUser{}, err
	}

	for _, user := range hbUsers {
		if user.ID == userID {
			return user, nil
		}
	}
	return HoneyBadgerUser{}, errors.New("User not found")
}

// FindUserByEmail - Returns a user by Email
func (hbc *HoneyBadgerClient) FindUserByEmail(userEmail string) (HoneyBadgerUser, error) {
	hbUsers, err := hbc.GetUsers()
	if err != nil {
		return HoneyBadgerUser{}, err
	}

	for _, user := range hbUsers {
		if user.Email == userEmail {
			return user, nil
		}
	}
	return HoneyBadgerUser{}, errors.New("User not found")
}

// CreateUser - Crea a HoneyBadger User
func (hbc *HoneyBadgerClient) CreateUser(userEmail string) error {
	var hbUser HoneyBadgerUser
	var jsonPayload = []byte(`{"team_invitation":{"email":"` + userEmail + `"}}`)

	url := fmt.Sprintf("%s/v2/teams/%d/team_invitations", hbc.HostURL, hbc.TeamID)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	body, err := hbc.doRequest(req)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &hbUser)
	if err != nil {
		return err
	}

	return nil
}

// UpdateUser - Update HoneyBadger User Information
func (hbc *HoneyBadgerClient) UpdateUser(userID int, isAdmin bool) error {
	var jsonPayload = []byte(`{"team_member":{"admin":` + strconv.FormatBool(isAdmin) + `}}`)

	url := fmt.Sprintf("%s/v2/teams/%d/team_members/%d", hbc.HostURL, hbc.TeamID, userID)
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return err
	}

	_, err = hbc.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

// DeleteUser - Delete HoneyBadger User
func (hbc *HoneyBadgerClient) DeleteUser(userID int) error {
	url := fmt.Sprintf("%s/v2/teams/%d/team_members/%d", hbc.HostURL, hbc.TeamID, userID)
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return err
	}

	_, err = hbc.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

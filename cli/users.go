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

// CreateUser - Create Honeybadger User
func (hbc *HoneybadgerClient) CreateUser(userEmail string, isAdmin bool, teamID int) error {
	var hbUser HoneybadgerUser
	var jsonPayload = []byte(`{"team_invitation":{"email":"` + userEmail + `", "admin":"` + strconv.FormatBool(isAdmin) + `"}}`)

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

// GetUserFromTeams - Get User information from Teams
func (hbc *HoneybadgerClient) GetUserFromTeams(userEmail string) (userTeams []HoneybadgerUser, err error) {
	insertedUser := make(map[string]bool)

	teams, err := hbc.GetTeams()
	if err != nil {
		return userTeams, err
	}

	for _, team := range teams {
		for _, user := range team.Users {
			if user.Email != userEmail {
				continue
			}
			user.TeamID = team.ID
			userTeams = append(userTeams, user)
			insertedUser[userEmail] = true
		}

		// Check if the user is invited but its not already in member list
		if _, ok := insertedUser[userEmail]; !ok {
			for _, userInvitation := range team.Invitations {
				if userInvitation.Email != userEmail {
					continue
				}
				userTeams = append(
					userTeams,
					HoneybadgerUser{
						ID:        userInvitation.ID,
						Email:     userInvitation.Email,
						IsAdmin:   userInvitation.IsAdmin,
						TeamID:    team.ID,
						CreatedAt: userInvitation.CreatedAt,
					},
				)

			}
		}

	}

	return userTeams, nil
}

// GetUserFromTeams - Get User information from Teams
func (hbc *HoneybadgerClient) GetUserForTeam(userEmail string, teamID int) (HoneybadgerUser, error) {
	userTeams, err := hbc.GetUserFromTeams(userEmail)
	if err != nil {
		return HoneybadgerUser{}, err
	}

	for _, userTeam := range userTeams {
		if userTeam.TeamID == teamID {
			return userTeam, nil
		}
	}

	return HoneybadgerUser{}, errors.New("User " + userEmail + "not found in team " + strconv.Itoa(teamID))
}

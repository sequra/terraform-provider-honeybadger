package cli

type HoneybadgerUsers struct {
	Users []HoneybadgerUser `json:"results"`
	Links HoneybadgerLink   `json:"links"`
}

type HoneybadgerUser struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	IsAdmin   bool   `json:"admin"`
}

type HoneybadgerLink struct {
	Self         string `json:"self"`
	PreviousPage string `json:"prev"`
	NextPage     string `json:"next"`
}

type HoneybadgerTeams struct {
	Teams []HoneybadgerTeam `json:"results"`
	Links HoneybadgerLink   `json:"links"`
}

type HoneybadgerTeam struct {
	ID        int                  `json:"id"`
	Name      string               `json:"name"`
	CreatedAt string               `json:"created_at"`
	Users     []HoneybadgerUser    `json:"members"`
	Projects  []HoneybadgerProject `json:"projects"`
	Owner     HoneybadgerTeamOwner `json:"owner"`
}

type HoneybadgerProject struct {
	ID                 int      `json:"id"`
	Name               string   `json:"name"`
	CreatedAt          string   `json:"created_at"`
	DisablePublicLinks bool     `json:"disable_public_links"`
	Token              string   `json:"token"`
	Environments       []string `json:"environments"`
}

type HoneybadgerTeamOwner struct {
	ID    int    `json:"id`
	Email string `json:"email"`
	Name  string `json:"name"`
}

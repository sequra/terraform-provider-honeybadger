package cli

type HoneyBadgerUsers struct {
	Users []HoneyBadgerUser `json:"results"`
	Links HoneyBadgerLink   `json:"links"`
}

type HoneyBadgerUser struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	IsAdmin   bool   `json:"admin"`
}

type HoneyBadgerLink struct {
	Self         string `json:"self"`
	PreviousPage string `json:"prev"`
	NextPage     string `json:"next"`
}

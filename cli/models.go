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

package model

// Organization represents a LinkedIn company or organization page.
type Organization struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	VanityName    string `json:"vanityName"`
	Description   string `json:"description"`
	LogoURL       string `json:"logoUrl"`
	Website       string `json:"website"`
	FollowerCount int    `json:"followerCount"`
}

// OrgFollowerStats contains follower statistics for an organization.
type OrgFollowerStats struct {
	OrganicCount int            `json:"organicCount"`
	PaidCount    int            `json:"paidCount"`
	TotalCount   int            `json:"totalCount"`
	ByFunction   map[string]int `json:"byFunction"`
	BySeniority  map[string]int `json:"bySeniority"`
}

// OrgPageStats contains page view statistics for an organization.
type OrgPageStats struct {
	Views          int    `json:"views"`
	UniqueVisitors int    `json:"uniqueVisitors"`
	Clicks         int    `json:"clicks"`
	Period         string `json:"period"`
}

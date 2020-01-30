package data

type Data struct {
	Service *Service
}

type Link struct {
	URL      string `json:"url" validate:"required"`
	Name     string `json:"name"`
	ImageURL string `json:"image_url"`
}

type Team struct {
	Name string `json:"name" validate:"required"`
}

type Tag struct {
	Key   string `json:"key" validate:"required"`
	Value string `json:"value" validate:"required"`
}

type Service struct {
	Name           string     `json:"name" validate:"required"`
	Description    string     `json:"description"`
	ImageUrl       string     `json:"image_url"`
	Hashtags       []*string  `json:"hashtags"`
	Tags           []*Tag     `json:"tags"`
	TeamOwner      *Team      `json:"team_owner"`
	Dependencies   []*Service `json:"dependencies"`
	Chat           *Link      `json:"chat"`
	Dashboards     []*Link    `json:"dashboards"`
	Documentation  *Link      `json:"documentation"`
	Email          *Link      `json:"email"`
	Runbook        *Link      `json:"runbook"`
	VersionControl *Link      `json:"version_control"`
}

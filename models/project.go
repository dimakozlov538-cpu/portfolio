package models

import "github.com/lib/pq"

type Profile struct {
	ID       int
	Name     string
	Title    string
	PhotoURL string
}

type About struct {
	ID          int
	Title       string
	Description string
}

type Skill struct {
	ID       int
	Name     string
	Category string
}

type Project struct {
	ID          int
	Title       string
	Description string
	ImageURL    string
	RepoURL     string
	DemoURL     string
	TechStack   pq.StringArray
}

type Hackathon struct {
	ID          int
	Title       string
	Description string
	Date        string
	Result      string
}

type Award struct {
	ID          int
	Title       string
	Description string
	Date        string
}

type Education struct {
	ID          int
	Degree      string
	Institution string
	Description string
	StartDate   string
	EndDate     string
}

type Contact struct {
	ID    int
	Type  string
	Value string
	URL   string
}

type PageData struct {
	Profile    Profile
	About      About
	Skills     []Skill
	Projects   []Project
	Hackathons []Hackathon
	Awards     []Award
	Education  []Education
	Contacts   []Contact
}

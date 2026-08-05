package domain

import "time"

// BadResponse Standard Format
type APIError struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}

// Multi-language JSON Field structure
type LocalizedString map[string]string

// Projects Domain Model
type Project struct {
	ID          string          `json:"id"`
	Slug        string          `json:"slug"`
	Title       string          `json:"title,omitempty"`
	Name        string          `json:"name,omitempty"` // Alias for GET compatibility
	Icon        string          `json:"icon"`
	Category    string          `json:"category"`
	Status      string          `json:"status"`
	Tagline     LocalizedString `json:"tagline"`
	Description LocalizedString `json:"description"`
	TechStack   []string        `json:"techStack"`
	DemoURL     string          `json:"demoUrl"`
	GithubURL   string          `json:"githubUrl"`
	CreatedAt   time.Time       `json:"createdAt,omitempty"`
	UpdatedAt   time.Time       `json:"updatedAt,omitempty"`
}

// Research Paper Domain Model
type ResearchPaper struct {
	ID              string    `json:"id"`
	Slug            string    `json:"slug"`
	Title           string    `json:"title"`
	Category        string    `json:"category"`
	Summary         string    `json:"summary"`
	PublicationDate string    `json:"date"`
	ReadTime        string    `json:"readTime"`
	Tags            []string  `json:"tags"`
	CreatedAt       time.Time `json:"createdAt,omitempty"`
}

// Tech Stack Item Model
type TechItem struct {
	Slug                 string          `json:"slug"`
	Name                 string          `json:"name"`
	IconSVGURL           string          `json:"iconSvgUrl"`
	Category             string          `json:"category"`
	YearsSinceFirstTouch int             `json:"yearsSinceFirstTouch"`
	Level                string          `json:"level"`
	ShortDesc            LocalizedString `json:"shortDesc"`
	ExperienceDetails    LocalizedString `json:"experienceDetails"`
}

// Telemetry Metrics Model
type Telemetry struct {
	ExperienceYears string `json:"experienceYears"`
	DailyEvents     string `json:"dailyEvents"`
	AvgLatencyMs    string `json:"avgLatencyMs"`
	ActiveSessions  string `json:"activeSessions"`
	SystemStatus    string `json:"systemStatus"`
}

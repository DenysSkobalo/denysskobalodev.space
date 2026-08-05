package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
)

type ProjectsHandler struct {
	DB *sql.DB
}

func (h *ProjectsHandler) List(w http.ResponseWriter, r *http.Request) {
	lang := r.URL.Query().Get("lang")
	statusFilter := r.URL.Query().Get("status")

	if h.DB == nil {
		mockProjects := []map[string]interface{}{
			{
				"slug":        "cloudflare-edge-engine",
				"name":        "Cloudflare Edge Engine",
				"icon":        "zap",
				"status":      "published",
				"tagline":     filterLanguage(map[string]string{"en": "WASM Go Backend Service"}, lang),
				"description": filterLanguage(map[string]string{"en": "High performance edge backend on Cloudflare Workers"}, lang),
				"techStack":   []string{"Go", "WebAssembly", "Cloudflare Workers"},
				"demoUrl":     "https://api.denysskobalodev.space",
				"githubUrl":   "https://github.com/DenysSkobalo/denysskobalodev.space",
			},
		}
		respondJSON(w, http.StatusOK, mockProjects)
		return
	}

	query := "SELECT slug, title, icon, status, tagline, description, tech_stack, demo_url, github_url FROM projects"
	var args []interface{}

	if statusFilter != "" {
		query += " WHERE status = ?"
		args = append(args, statusFilter)
	}

	rows, err := h.DB.Query(query, args...)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "INTERNAL_ERROR", "Database fetch error")
		return
	}
	defer rows.Close()

	projects := make([]map[string]interface{}, 0)
	for rows.Next() {
		var slug, title, icon, status, taglineJSON, descJSON, techStackJSON, demoURL, githubURL string
		_ = rows.Scan(&slug, &title, &icon, &status, &taglineJSON, &descJSON, &techStackJSON, &demoURL, &githubURL)

		var tagline, desc map[string]string
		var techStack []string
		_ = json.Unmarshal([]byte(taglineJSON), &tagline)
		_ = json.Unmarshal([]byte(descJSON), &desc)
		_ = json.Unmarshal([]byte(techStackJSON), &techStack)

		item := map[string]interface{}{
			"slug":      slug,
			"name":      title,
			"icon":      icon,
			"status":    status,
			"tagline":   filterLanguage(tagline, lang),
			"description": filterLanguage(desc, lang),
			"techStack": techStack,
			"demoUrl":   demoURL,
			"githubUrl": githubURL,
		}
		projects = append(projects, item)
	}

	respondJSON(w, http.StatusOK, projects)
}

func (h *ProjectsHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string   `json:"title"`
		Icon        string   `json:"icon"`
		Status      string   `json:"status"`
		Category    string   `json:"category"`
		Tagline     string   `json:"tagline"`
		Description string   `json:"description"`
		Stack       []string `json:"stack"`
		DemoURL     string   `json:"demoUrl"`
		GithubURL   string   `json:"githubUrl"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil || input.Title == "" || input.Status == "" {
		respondError(w, http.StatusBadRequest, "VALIDATION_FAILED", "Title and status are required fields")
		return
	}

	slug := slugify(input.Title)
	taglineJSON, _ := json.Marshal(map[string]string{"en": input.Tagline})
	descJSON, _ := json.Marshal(map[string]string{"en": input.Description})
	stackJSON, _ := json.Marshal(input.Stack)

	query := `INSERT INTO projects (id, slug, title, icon, category, status, tagline, description, tech_stack, demo_url, github_url) 
	          VALUES (lower(hex(randomblob(16))), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := h.DB.Exec(query, slug, input.Title, input.Icon, input.Category, input.Status, string(taglineJSON), string(descJSON), string(stackJSON), input.DemoURL, input.GithubURL)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			respondError(w, http.StatusConflict, "DUPLICATE_SLUG", "Project with slug '"+slug+"' already exists")
			return
		}
		respondError(w, http.StatusInternalServerError, "DB_WRITE_ERROR", err.Error())
		return
	}

	// Immediate response to guarantee Read-After-Write consistency on Client
	respondJSON(w, http.StatusCreated, map[string]interface{}{
		"slug":      slug,
		"title":     input.Title,
		"icon":      input.Icon,
		"status":    input.Status,
		"category":  input.Category,
		"tagline":   map[string]string{"en": input.Tagline},
		"description": map[string]string{"en": input.Description},
		"techStack": input.Stack,
		"demoUrl":   input.DemoURL,
		"githubUrl": input.GithubURL,
		"createdAt": timeNowISO(),
	})
}

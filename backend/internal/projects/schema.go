package projects

import "time"

type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type UpdateProjectRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
}

type ProjectResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	OwnerID     string    `json:"owner_id"`
	CreatedAt   time.Time `json:"created_at"`
}

type PaginationMeta struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
	Total int `json:"total"`
}

type ListProjectsResponse struct {
	Projects   []ProjectResponse `json:"projects"`
	Pagination PaginationMeta    `json:"pagination"`
}

type ProjectDetailsResponse struct {
	ProjectResponse
	Tasks []Task `json:"tasks"`
}

type StatusCount struct {
	Status string `json:"status"`
	Count  int    `json:"count"`
}

type AssigneeCount struct {
	AssigneeID   *string `json:"assignee_id"`
	AssigneeName string  `json:"assignee_name"`
	Count        int     `json:"count"`
}

type ProjectStatsResponse struct {
	StatusCounts   []StatusCount   `json:"status_counts"`
	AssigneeCounts []AssigneeCount `json:"assignee_counts"`
}

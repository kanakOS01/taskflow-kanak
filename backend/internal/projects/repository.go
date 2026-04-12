package projects

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("project not found")

type Repository interface {
	Create(ctx context.Context, project *Project) error
	List(ctx context.Context, userID string, page, limit int) ([]Project, int, error)
	GetByID(ctx context.Context, id string) (*Project, []Task, error)
	GetStats(ctx context.Context, id string) ([]StatusCount, []AssigneeCount, error)
	Update(ctx context.Context, project *Project) error
	Delete(ctx context.Context, id string) error
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, p *Project) error {
	query := `
		INSERT INTO projects (id, name, description, owner_id, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(ctx, query, p.ID, p.Name, p.Description, p.OwnerID, p.CreatedAt, p.UpdatedAt)
	return err
}

func (r *postgresRepository) List(ctx context.Context, userID string, page, limit int) ([]Project, int, error) {
	countQuery := `
		SELECT COUNT(DISTINCT p.id)
		FROM projects p
		LEFT JOIN tasks t ON t.project_id = p.id
		WHERE p.owner_id = $1 OR t.assignee_id = $1
	`
	var total int
	if err := r.db.QueryRow(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.owner_id, p.created_at, p.updated_at
		FROM projects p
		LEFT JOIN tasks t ON t.project_id = p.id
		WHERE p.owner_id = $1 OR t.assignee_id = $1
		ORDER BY p.created_at DESC
		LIMIT $2 OFFSET $3
	`
	rows, err := r.db.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, 0, err
		}
		projects = append(projects, p)
	}
	return projects, total, rows.Err()
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*Project, []Task, error) {
	query := `SELECT id, name, description, owner_id, created_at, updated_at FROM projects WHERE id = $1`
	var p Project
	err := r.db.QueryRow(ctx, query, id).Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil, ErrNotFound
		}
		return nil, nil, err
	}

	taskQuery := `SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at FROM tasks WHERE project_id = $1`
	rows, err := r.db.Query(ctx, taskQuery, id)
	if err != nil {
		return &p, nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return &p, nil, err
		}
		tasks = append(tasks, t)
	}
	return &p, tasks, nil
}

func (r *postgresRepository) GetStats(ctx context.Context, id string) ([]StatusCount, []AssigneeCount, error) {
	statusQuery := `
		SELECT status, count(*) 
		FROM tasks 
		WHERE project_id = $1 
		GROUP BY status
	`
	sRows, err := r.db.Query(ctx, statusQuery, id)
	if err != nil {
		return nil, nil, err
	}
	defer sRows.Close()

	var statusCounts []StatusCount
	for sRows.Next() {
		var sc StatusCount
		if err := sRows.Scan(&sc.Status, &sc.Count); err != nil {
			return nil, nil, err
		}
		statusCounts = append(statusCounts, sc)
	}

	assigneeQuery := `
		SELECT 
			t.assignee_id, 
			COALESCE(u.name, 'Unassigned') as assignee_name, 
			count(*) 
		FROM tasks t 
		LEFT JOIN users u ON t.assignee_id = u.id 
		WHERE t.project_id = $1 
		GROUP BY t.assignee_id, u.name
	`
	aRows, err := r.db.Query(ctx, assigneeQuery, id)
	if err != nil {
		return statusCounts, nil, err
	}
	defer aRows.Close()

	var assigneeCounts []AssigneeCount
	for aRows.Next() {
		var ac AssigneeCount
		if err := aRows.Scan(&ac.AssigneeID, &ac.AssigneeName, &ac.Count); err != nil {
			return statusCounts, nil, err
		}
		assigneeCounts = append(assigneeCounts, ac)
	}

	return statusCounts, assigneeCounts, nil
}

func (r *postgresRepository) Update(ctx context.Context, p *Project) error {
	query := `UPDATE projects SET name = $1, description = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(ctx, query, p.Name, p.Description, p.UpdatedAt, p.ID)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM projects WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

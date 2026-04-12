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
	List(ctx context.Context, userID string) ([]Project, error)
	GetByID(ctx context.Context, id string) (*Project, []Task, error)
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

func (r *postgresRepository) List(ctx context.Context, userID string) ([]Project, error) {
	query := `
		SELECT DISTINCT p.id, p.name, p.description, p.owner_id, p.created_at, p.updated_at 
		FROM projects p 
		LEFT JOIN tasks t ON t.project_id = p.id
		WHERE p.owner_id = $1 OR t.assignee_id = $1
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var projects []Project
	for rows.Next() {
		var p Project
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.OwnerID, &p.CreatedAt, &p.UpdatedAt); err != nil {
			return nil, err
		}
		projects = append(projects, p)
	}
	return projects, rows.Err()
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

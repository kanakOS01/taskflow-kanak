package tasks

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("task not found")

type Repository interface {
	Create(ctx context.Context, task *Task) error
	GetByID(ctx context.Context, id string) (*Task, error)
	ListByProject(ctx context.Context, projectID string, status string, assignee string) ([]Task, error)
	Update(ctx context.Context, task *Task) error
	Delete(ctx context.Context, id string) error
	GetProjectOwner(ctx context.Context, projectID string) (string, error)
}

type postgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) Repository {
	return &postgresRepository{db: db}
}

func (r *postgresRepository) Create(ctx context.Context, t *Task) error {
	query := `
		INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	_, err := r.db.Exec(ctx, query, t.ID, t.Title, t.Description, t.Status, t.Priority, t.ProjectID, t.AssigneeID, t.DueDate, t.CreatedAt, t.UpdatedAt)
	return err
}

func (r *postgresRepository) GetByID(ctx context.Context, id string) (*Task, error) {
	query := `SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at FROM tasks WHERE id = $1`
	var t Task
	err := r.db.QueryRow(ctx, query, id).Scan(
		&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, 
		&t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	return &t, nil
}

func (r *postgresRepository) ListByProject(ctx context.Context, projectID string, status string, assignee string) ([]Task, error) {
	query := `SELECT id, title, description, status, priority, project_id, assignee_id, due_date, created_at, updated_at FROM tasks WHERE project_id = $1`
	
	args := []interface{}{projectID}
	argIdx := 2

	if status != "" {
		query += fmt.Sprintf(" AND status = $%d", argIdx)
		args = append(args, status)
		argIdx++
	}

	if assignee != "" {
		query += fmt.Sprintf(" AND assignee_id = $%d", argIdx)
		args = append(args, assignee)
	}

	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var t Task
		if err := rows.Scan(&t.ID, &t.Title, &t.Description, &t.Status, &t.Priority, &t.ProjectID, &t.AssigneeID, &t.DueDate, &t.CreatedAt, &t.UpdatedAt); err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}
	return tasks, rows.Err()
}

func (r *postgresRepository) Update(ctx context.Context, t *Task) error {
	query := `UPDATE tasks SET title = $1, description = $2, status = $3, priority = $4, assignee_id = $5, due_date = $6, updated_at = $7 WHERE id = $8`
	_, err := r.db.Exec(ctx, query, t.Title, t.Description, t.Status, t.Priority, t.AssigneeID, t.DueDate, t.UpdatedAt, t.ID)
	return err
}

func (r *postgresRepository) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM tasks WHERE id = $1`
	_, err := r.db.Exec(ctx, query, id)
	return err
}

func (r *postgresRepository) GetProjectOwner(ctx context.Context, projectID string) (string, error) {
	var ownerID string
	query := `SELECT owner_id FROM projects WHERE id = $1`
	err := r.db.QueryRow(ctx, query, projectID).Scan(&ownerID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrNotFound
		}
		return "", err
	}
	return ownerID, nil
}

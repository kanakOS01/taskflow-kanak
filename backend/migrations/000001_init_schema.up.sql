-- USERS
CREATE TABLE users (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- PROJECTS
CREATE TABLE projects (
    id UUID PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- LOOKUP: STATUS
CREATE TABLE task_statuses (
    value TEXT PRIMARY KEY
);

-- LOOKUP: PRIORITY (with rank for ordering)
CREATE TABLE task_priorities (
    value TEXT PRIMARY KEY,
    rank INT NOT NULL UNIQUE
);

-- SEED DATA
INSERT INTO task_statuses (value) VALUES
('todo'),
('in_progress'),
('done');

INSERT INTO task_priorities (value, rank) VALUES
('low', 1),
('medium', 10),
('high', 100);

-- TASKS
CREATE TABLE tasks (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT,

    status TEXT NOT NULL DEFAULT 'todo'
        REFERENCES task_statuses(value),

    priority TEXT NOT NULL DEFAULT 'medium'
        REFERENCES task_priorities(value),

    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    assignee_id UUID REFERENCES users(id) ON DELETE SET NULL,
    due_date DATE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- INDEXES
CREATE INDEX idx_tasks_project_id ON tasks(project_id);
CREATE INDEX idx_tasks_assignee_id ON tasks(assignee_id);
CREATE INDEX idx_projects_owner_id ON projects(owner_id);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_task_priorities_rank ON task_priorities(rank);

INSERT INTO users (id, name, email, password)
VALUES ('7b1897c5-559d-43dd-bbf7-268a983b6f00', 'Test User', 'test@example.com', '$2a$12$R9h/cIPz0gi.URNNX3kh2OPST9/PgBkqquzi.Ss7KIUgO2t0jWMUW')
ON CONFLICT (id) DO NOTHING;

INSERT INTO projects (id, name, description, owner_id)
VALUES ('8b1897c5-559d-43dd-bbf7-268a983b6f01', 'Test Project', 'A seeded project for testing', '7b1897c5-559d-43dd-bbf7-268a983b6f00')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id)
VALUES 
    ('9b1897c5-559d-43dd-bbf7-268a983b6f02', 'Task 1', 'Test task todo', 'todo', 'high', '8b1897c5-559d-43dd-bbf7-268a983b6f01', '7b1897c5-559d-43dd-bbf7-268a983b6f00'),
    ('9b1897c5-559d-43dd-bbf7-268a983b6f03', 'Task 2', 'Test task in progress', 'in_progress', 'medium', '8b1897c5-559d-43dd-bbf7-268a983b6f01', NULL),
    ('9b1897c5-559d-43dd-bbf7-268a983b6f04', 'Task 3', 'Test task done', 'done', 'low', '8b1897c5-559d-43dd-bbf7-268a983b6f01', '7b1897c5-559d-43dd-bbf7-268a983b6f00')
ON CONFLICT (id) DO NOTHING;

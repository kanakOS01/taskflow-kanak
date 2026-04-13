INSERT INTO users (id, name, email, password)
VALUES ('0dfba3db-ffb9-49df-887c-2ebba8c24088', 'kanak', 'kanak@email.com', '$2a$12$wDbsQl/B2GvIyt..wQuwUeymD439EIVYKV1P2vMmdinR9o.Mkhxs.')
ON CONFLICT (id) DO NOTHING;

INSERT INTO users (id, name, email, password)
VALUES ('156d0a79-b4a5-4c21-9f24-555b99675f45', 'test', 'test@email.com', '$2a$12$wDbsQl/B2GvIyt..wQuwUeymD439EIVYKV1P2vMmdinR9o.Mkhxs.')
ON CONFLICT (id) DO NOTHING;

INSERT INTO projects (id, name, description, owner_id)
VALUES ('8b1897c5-559d-43dd-bbf7-268a983b6f01', 'Test Project', 'A seeded project for testing', '0dfba3db-ffb9-49df-887c-2ebba8c24088')
ON CONFLICT (id) DO NOTHING;

INSERT INTO tasks (id, title, description, status, priority, project_id, assignee_id, created_by)
VALUES 
    ('9b1897c5-559d-43dd-bbf7-268a983b6f02', 'Task 1', 'Test task todo', 'todo', 'high', '8b1897c5-559d-43dd-bbf7-268a983b6f01', '0dfba3db-ffb9-49df-887c-2ebba8c24088', '0dfba3db-ffb9-49df-887c-2ebba8c24088'),
    ('9b1897c5-559d-43dd-bbf7-268a983b6f03', 'Task 2', 'Test task in progress', 'in_progress', 'medium', '8b1897c5-559d-43dd-bbf7-268a983b6f01', NULL, '0dfba3db-ffb9-49df-887c-2ebba8c24088'),
    ('9b1897c5-559d-43dd-bbf7-268a983b6f04', 'Task 3', 'Test task done', 'done', 'low', '8b1897c5-559d-43dd-bbf7-268a983b6f01', '0dfba3db-ffb9-49df-887c-2ebba8c24088', '0dfba3db-ffb9-49df-887c-2ebba8c24088')
ON CONFLICT (id) DO NOTHING;

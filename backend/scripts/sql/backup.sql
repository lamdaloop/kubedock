CREATE TABLE IF NOT EXISTS backup_history (
    id SERIAL PRIMARY KEY,
    cluster_id TEXT REFERENCES clusters(id),
    status TEXT NOT NULL,
    path TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

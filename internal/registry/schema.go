package registry

import (
	"fmt"
)

func (r *Registry) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS projects (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_name TEXT NOT NULL,
		instance_name TEXT NOT NULL,
		path TEXT NOT NULL,
		UNIQUE(project_name, instance_name)
	);

	CREATE TABLE IF NOT EXISTS allocations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id INTEGER NOT NULL,
		service_name TEXT NOT NULL,
		port INTEGER NOT NULL UNIQUE,
		pid INTEGER,
		status TEXT NOT NULL, -- 'reserved', 'running', 'stale'
		updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(project_id) REFERENCES projects(id) ON DELETE CASCADE
	);
	`
	_, err := r.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("schema migration failed: %w", err)
	}
	return nil
}

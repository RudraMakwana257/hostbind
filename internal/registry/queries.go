package registry

import (
	"database/sql"
	"fmt"
)

type Allocation struct {
	ID          int
	ProjectName string
	Instance    string
	ServiceName string
	Port        int
	PID         int
	Status      string
}

// EnsureProject returns the ID of a project, creating it if it doesn't exist.
func (r *Registry) EnsureProject(projectName, instanceName, path string) (int, error) {
	var id int
	err := r.db.QueryRow(`SELECT id FROM projects WHERE project_name = ? AND instance_name = ?`, projectName, instanceName).Scan(&id)
	
	if err == sql.ErrNoRows {
		res, err := r.db.Exec(`INSERT INTO projects (project_name, instance_name, path) VALUES (?, ?, ?)`, projectName, instanceName, path)
		if err != nil {
			return 0, fmt.Errorf("failed to insert project: %w", err)
		}
		lastID, _ := res.LastInsertId()
		return int(lastID), nil
	} else if err != nil {
		return 0, fmt.Errorf("failed to query project: %w", err)
	}
	
	return id, nil
}

// GetActiveAllocations returns all currently active ports across all projects.
func (r *Registry) GetActiveAllocations() ([]Allocation, error) {
	rows, err := r.db.Query(`
		SELECT p.project_name, p.instance_name, a.service_name, a.port, a.pid, a.status 
		FROM allocations a
		JOIN projects p ON p.id = a.project_id
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allocs []Allocation
	for rows.Next() {
		var a Allocation
		if err := rows.Scan(&a.ProjectName, &a.Instance, &a.ServiceName, &a.Port, &a.PID, &a.Status); err != nil {
			return nil, err
		}
		allocs = append(allocs, a)
	}
	return allocs, nil
}

// ReservePort safely reserves a port for a service to prevent race conditions.
func (r *Registry) ReservePort(projectID int, serviceName string, port int) error {
	_, err := r.db.Exec(`
		INSERT INTO allocations (project_id, service_name, port, status) 
		VALUES (?, ?, ?, 'reserved')
		ON CONFLICT(port) DO UPDATE SET 
			project_id=excluded.project_id, 
			service_name=excluded.service_name, 
			status='reserved', 
			updated_at=CURRENT_TIMESTAMP
	`, projectID, serviceName, port)
	return err
}

// UpdatePID updates a reserved port to 'running' with the actual process ID.
func (r *Registry) UpdatePID(port int, pid int) error {
	_, err := r.db.Exec(`UPDATE allocations SET pid = ?, status = 'running', updated_at = CURRENT_TIMESTAMP WHERE port = ?`, pid, port)
	return err
}

// ReleasePort releases a port when a service stops.
func (r *Registry) ReleasePort(port int) error {
	_, err := r.db.Exec(`DELETE FROM allocations WHERE port = ?`, port)
	return err
}

// GetProjectAllocations returns allocations for a specific project instance.
func (r *Registry) GetProjectAllocations(projectName, instanceName string) ([]Allocation, error) {
	rows, err := r.db.Query(`
		SELECT p.project_name, p.instance_name, a.service_name, a.port, a.pid, a.status 
		FROM allocations a
		JOIN projects p ON p.id = a.project_id
		WHERE p.project_name = ? AND p.instance_name = ?
	`, projectName, instanceName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var allocs []Allocation
	for rows.Next() {
		var a Allocation
		if err := rows.Scan(&a.ProjectName, &a.Instance, &a.ServiceName, &a.Port, &a.PID, &a.Status); err != nil {
			return nil, err
		}
		allocs = append(allocs, a)
	}
	return allocs, nil
}

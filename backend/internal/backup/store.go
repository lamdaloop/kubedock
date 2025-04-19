package backup

import (
    "database/sql"
)

type HistoryEntry struct {
    ClusterID string
    Status    string
    Path      string
}

type BackupStore interface {
    SaveHistory(entry HistoryEntry) error
}

type PostgresStore struct {
    DB *sql.DB
}

func (s *PostgresStore) SaveHistory(entry HistoryEntry) error {
    _, err := s.DB.Exec(`
        INSERT INTO backup_history (cluster_id, status, path)
        VALUES ($1, $2, $3)
    `, entry.ClusterID, entry.Status, entry.Path)
    return err
}

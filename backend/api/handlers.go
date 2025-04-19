package api

import (
    "encoding/json"
    "fmt"
    "net/http"
    "log"

    "github.com/gorilla/mux"
    "github.com/lamdaloop/kubedock/backend/internal/backup"
)

type Cluster struct {
    ID    string `json:"id" db:"id"`
    Name  string `json:"name" db:"name"`
    URL   string `json:"url" db:"url"`
    Token string `json:"token" db:"token"`
}

type BackupRecord struct {
    ID        int       `db:"id" json:"id"`
    Status    string    `db:"status" json:"status"`
    Path      string    `db:"path" json:"path"`
    CreatedAt string    `db:"created_at" json:"created_at"`
}

// CreateClusterHandler inserts or updates a cluster in the DB
func CreateClusterHandler(w http.ResponseWriter, r *http.Request) {
    var cluster Cluster
    if err := json.NewDecoder(r.Body).Decode(&cluster); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    _, err := DB.Exec(`
        INSERT INTO clusters (id, name, url, token)
        VALUES ($1, $2, $3, $4)
        ON CONFLICT (id) DO UPDATE
        SET name = EXCLUDED.name, url = EXCLUDED.url, token = EXCLUDED.token
    `, cluster.ID, cluster.Name, cluster.URL, cluster.Token)

    if err != nil {
        http.Error(w, "Failed to store cluster", http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(cluster)
}

// ListClustersHandler returns all clusters from the DB
func ListClustersHandler(w http.ResponseWriter, r *http.Request) {
    var clusters []Cluster
    err := DB.Select(&clusters, `SELECT * FROM clusters`)
    if err != nil {
        http.Error(w, "Failed to list clusters", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(clusters)
}


// TriggerBackupHandler performs a backup for a specific cluster
func TriggerBackupHandler(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]

    var cluster Cluster
    err := DB.Get(&cluster, `SELECT * FROM clusters WHERE id=$1`, id)
    if err != nil {
        http.Error(w, "Cluster not found", http.StatusNotFound)
        return
    }

    fmt.Printf("📦 Triggering backup for cluster: %s (%s)\n", cluster.Name, cluster.URL)

    status, path, backupErr := backup.RunBackup(cluster.URL, cluster.Token, cluster.ID)

    _, dbErr := DB.Exec(`
        INSERT INTO backup_history (cluster_id, status, path)
        VALUES ($1, $2, $3)
    `, cluster.ID, status, path)

    if dbErr != nil {
        log.Println("⚠️ Failed to save backup history:", dbErr)
    }

    if backupErr != nil {
        http.Error(w, fmt.Sprintf("Backup failed: %v", backupErr), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    fmt.Fprintf(w, "Backup completed successfully. Saved to %s", path)
}



func GetBackupHistoryHandler(w http.ResponseWriter, r *http.Request) {
    id := mux.Vars(r)["id"]

    var history []BackupRecord
    query := `
        SELECT id, status, path, created_at
        FROM backup_history
        WHERE cluster_id = $1
        ORDER BY created_at DESC
    `

    err := DB.Select(&history, query, id)
    if err != nil {
        fmt.Printf("❌ Failed to get history for cluster '%s': %v\n", id, err)
        http.Error(w, "Failed to get history: "+err.Error(), http.StatusInternalServerError)
        return
    }

    fmt.Printf("📜 Returning %d history entries for cluster '%s'\n", len(history), id)
    json.NewEncoder(w).Encode(history)
}


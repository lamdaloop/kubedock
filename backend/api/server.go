package api

import (
    "github.com/gorilla/mux"
)

func NewRouter() *mux.Router {
    r := mux.NewRouter()

    r.HandleFunc("/clusters", CreateClusterHandler).Methods("POST")
    r.HandleFunc("/clusters", ListClustersHandler).Methods("GET")
    r.HandleFunc("/clusters/{id}/backup", TriggerBackupHandler).Methods("POST")
    r.HandleFunc("/clusters/{id}/history", GetBackupHistoryHandler).Methods("GET")


    return r
}

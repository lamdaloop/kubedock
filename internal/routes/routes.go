package routes

import (
	"fmt"
	_ "fmt"
	_ "kubedock.org/kubedock-package/kubedock/k8s-client"
	client_config "kubedock.org/kubedock-package/kubedock/k8s-client"
	"net/http"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", healthCheck)
	mux.HandleFunc("/fetch-pod-manifest", FetchPodManifest)

	return mux
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, "helath-check reached!")
}

func FetchPodManifest(w http.ResponseWriter, r *http.Request) {
	//To Do
	var clientset, err = client_config.GetKubeConfig()
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create Kubernetes client: %s", err), http.StatusInternalServerError)
		return
	}
	namespace := r.URL.Query().Get("namespace")
	podName := r.URL.Query().Get("pod")

	if namespace == "" || podName == "" {
		http.Error(w, fmt.Sprintf("Failed to get namespace and pod name from URL: %s", r.URL), http.StatusInternalServerError)
		return
	}

	manifestOutput, err := client_config.FetchPodManifest(clientset, namespace, podName) // Corrected this line
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching pod manifest: %s", err), http.StatusInternalServerError)
		return
	}

	// Write the YAML to the response
	w.Header().Set("Content-Type", "application/x-yaml")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(manifestOutput)) // Corrected the typo here
}

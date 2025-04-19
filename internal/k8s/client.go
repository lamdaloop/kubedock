package k8s

import (
	"os"
    "log"

    "k8s.io/client-go/dynamic"
    "k8s.io/client-go/discovery"
    "k8s.io/client-go/rest"
)

func CreateClient() (*rest.Config, *discovery.DiscoveryClient, dynamic.Interface) {
    server := os.Getenv("KUBEDOCK_SERVER_URL")
    token := os.Getenv("KUBEDOCK_BEARER_TOKEN")

    if server == "" || token == "" {
        log.Fatalf("❌ Missing KUBEDOCK_SERVER_URL or KUBEDOCK_BEARER_TOKEN")
    }

    config := &rest.Config{
        Host:        server,
        BearerToken: token,
        TLSClientConfig: rest.TLSClientConfig{
            Insecure: true, // set to false if you're providing a CA cert
        },
    }

    discoClient, err := discovery.NewDiscoveryClientForConfig(config)
    if err != nil {
        log.Fatalf("❌ Failed to create discovery client: %v", err)
    }

    dynClient, err := dynamic.NewForConfig(config)
    if err != nil {
        log.Fatalf("❌ Failed to create dynamic client: %v", err)
    }

    return config, discoClient, dynClient
}

func homeDir() string {
    return os.Getenv("HOME")
}

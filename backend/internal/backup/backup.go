package backup

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"
    "time"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/client-go/discovery"
    "k8s.io/client-go/dynamic"
    "k8s.io/client-go/rest"
)

func RunBackup(serverURL, bearerToken, clusterID string) (status string, path string, err error) {
    config := &rest.Config{
        Host:        serverURL,
        BearerToken: bearerToken,
        TLSClientConfig: rest.TLSClientConfig{
            Insecure: true,
        },
    }

    discoClient, err := discovery.NewDiscoveryClientForConfig(config)
    if err != nil {
        return "failed", "", fmt.Errorf("discovery client error: %w", err)
    }

    dynClient, err := dynamic.NewForConfig(config)
    if err != nil {
        return "failed", "", fmt.Errorf("dynamic client error: %w", err)
    }

    gvrs, err := DiscoverResources(discoClient)
    if err != nil {
        return "failed", "", fmt.Errorf("discover resources error: %w", err)
    }

    timestamp := time.Now().Format("2006-01-02_15-04-05")
    backupDir := filepath.Join("backups", clusterID, timestamp)

    err = FetchAndDumpResources(dynClient, gvrs, backupDir)
    if err != nil {
        return "failed", backupDir, err
    }

    return "success", backupDir, nil
}

func DiscoverResources(disco *discovery.DiscoveryClient) ([]schema.GroupVersionResource, error) {
    resourceList, err := disco.ServerPreferredResources()
    if err != nil {
        return nil, err
    }

    var gvrs []schema.GroupVersionResource
    for _, list := range resourceList {
        for _, res := range list.APIResources {
            if !containsVerb(res.Verbs, "list") || res.Name == "events" {
                continue
            }
            gv, err := schema.ParseGroupVersion(list.GroupVersion)
            if err != nil {
                continue
            }
            gvrs = append(gvrs, schema.GroupVersionResource{
                Group:    gv.Group,
                Version:  gv.Version,
                Resource: res.Name,
            })
        }
    }
    return gvrs, nil
}

func FetchAndDumpResources(client dynamic.Interface, gvrs []schema.GroupVersionResource, backupPath string) error {
    for _, gvr := range gvrs {
        list, err := client.Resource(gvr).Namespace("").List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            continue
        }

        for _, item := range list.Items {
            ns := item.GetNamespace()
            name := item.GetName()

            var path string
            if ns == "" {
                path = filepath.Join(backupPath, "cluster", gvr.Resource)
            } else {
                path = filepath.Join(backupPath, ns, gvr.Resource)
            }

            os.MkdirAll(path, 0755)
            outFile := filepath.Join(path, fmt.Sprintf("%s.yaml", name))
            data, err := json.MarshalIndent(item.Object, "", "  ")
            if err != nil {
                continue
            }

            os.WriteFile(outFile, data, 0644)
        }
    }

    return nil
}

func containsVerb(verbs []string, verb string) bool {
    for _, v := range verbs {
        if v == verb {
            return true
        }
    }
    return false
}

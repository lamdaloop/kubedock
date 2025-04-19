package k8s

import (
    "context"
    "encoding/json"
    "fmt"
    "os"
    "path/filepath"

    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
    "k8s.io/apimachinery/pkg/runtime/schema"
    "k8s.io/client-go/discovery"
    "k8s.io/client-go/dynamic"
)

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
    total := 0
    for _, gvr := range gvrs {
        resourceList, err := client.Resource(gvr).Namespace("").List(context.TODO(), metav1.ListOptions{})
        if err != nil {
            fmt.Printf("⚠️  Skipping %s.%s (%s): %v\n", gvr.Resource, gvr.Group, gvr.Version, err)
            continue
        }

        for _, item := range resourceList.Items {
            ns := item.GetNamespace()
            name := item.GetName()

            var path string
            if ns == "" {
                path = filepath.Join(backupPath, "cluster", gvr.Resource)
            } else {
                path = filepath.Join(backupPath, ns, gvr.Resource)
            }

            if err := os.MkdirAll(path, 0755); err != nil {
                fmt.Printf("❌ Failed to create path: %s | %v\n", path, err)
                continue
            }

            outFile := filepath.Join(path, fmt.Sprintf("%s.yaml", name))
            data, err := json.MarshalIndent(item.Object, "", "  ")
            if err != nil {
                fmt.Printf("❌ Failed to marshal %s/%s: %v\n", ns, name, err)
                continue
            }

            err = os.WriteFile(outFile, data, 0644)
            if err != nil {
                fmt.Printf("❌ Failed to write %s: %v\n", outFile, err)
                continue
            }

            if ns == "" {
                fmt.Printf("📦 Backed up: [%s] %s\n", gvr.Resource, name)
            } else {
                fmt.Printf("📦 Backed up: [%s] %s/%s\n", gvr.Resource, ns, name)
            }

            total++
        }
    }

    fmt.Printf("\n✅ Backup complete! Total resources backed up: %d\n", total)
    return nil
}

func containsVerb(verbs []string, target string) bool {
    for _, v := range verbs {
        if v == target {
            return true
        }
    }
    return false
}

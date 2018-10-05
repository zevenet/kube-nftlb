package main

import (
    "fmt"

    v1 "k8s.io/api/core/v1"
    fields "k8s.io/apimachinery/pkg/fields"
    kubernetes "k8s.io/client-go/kubernetes"
    rest "k8s.io/client-go/rest"
    cache "k8s.io/client-go/tools/cache"
)

func main() {
    // Get in-cluster config
    config, err := rest.InClusterConfig()
    if err != nil {
        panic(err.Error())
    }
    // Get ClientSet for auth
    clientset, err := kubernetes.NewForConfig(config)
    if err != nil {
        panic(err.Error())
    }
    // Make list of resources to be observed
    watchlist := cache.NewListWatchFromClient(
        clientset.CoreV1().RESTClient(),    // REST interface
        string(v1.ResourceServices),        // Resource to watch for: "Service"
        v1.NamespaceAll,                    // Resource can be found in ALL namespaces
        fields.Everything(),                // Get ALL fields from requested resource
    )
    // Make channel before writing messages
    receiveCh := make(chan string)
    // Notify every change based on watchlist
    _, controller := cache.NewInformer(
        watchlist,          // Resources to watch for
        &v1.Service{},      // Resource struct
        0,
        // Event handler: new, deleted or updated resource
        cache.ResourceEventHandlerFuncs {
            AddFunc: func(obj interface{}) {
                receiveCh <- fmt.Sprintf("New Service: %s\n\n", obj)
            },
            DeleteFunc: func(obj interface{}) {
                receiveCh <- fmt.Sprintf("Deleted Service: %s\n\n", obj)
            },
            UpdateFunc: func(oldObj, newObj interface{}) {
                receiveCh <- fmt.Sprintf("Updated Service:\nBEFORE: %s\nNOW: %s\n\n", oldObj, newObj)
            },
        },
    )
    // Make stop channel and defer its close()
    stopCh := make(chan struct{})
    defer close(stopCh)
    // Event loop start, run it as background process
    go controller.Run(stopCh)
    // Print every message received from the channel
    for message := range receiveCh {
        fmt.Println(message)
    }
    // This line is unreachable: working as intended
}

K8s Utils
=========

Tools for a nicer Kubernetes development experience.

## Tools

### Unmarshal a Kubernetes ConfigMap to a Struct

```go
package main

import (
	"fmt"

	"github.com/nickschuch/k8sutils/configmap"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Backend struct {
	Host string `k8s_configmap:"backend.host"`
	Port string `k8s_configmap:"backend.port"`
}

func main() {
	cfg := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "example",
			Name:      "conf",
		},
		Data: map[string]string{
			"backend.host": "1.1.1.1",
			"backend.port": "443",
		},
	}

	var b Backend

	err := configmap.Unmarshal(cfg, &b)
	if err != nil {
		panic(err)
	}

	fmt.Println("Host:", b.Host)
	fmt.Println("Port:", b.Port)
}
```

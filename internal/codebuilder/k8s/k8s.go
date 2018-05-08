package k8s

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	skprcronjob "github.com/previousnext/skpr/utils/k8s/cronjob"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// SyncParams for function Sync.
type SyncParams struct {
	Namespace string
	Frequency string
	Image     string
	Role      string
	Bucket    string
	Keys      Keys
	Resources Resources
	Docker    Docker
	AWS       AWSCredentials
}

// Sync ConfigMaps to CronJobs to build Mysql containers.
func Sync(w io.Writer, params SyncParams) error {
	fmt.Println("Loading Kubernets config")

	k8sconfig, err := rest.InClusterConfig()
	if err != nil {
		panic(err)
	}

	fmt.Println("Connecting to Kubernetes")

	k8sclient, err := kubernetes.NewForConfig(k8sconfig)
	if err != nil {
		panic(err)
	}

	fmt.Fprintln(w, "Querying Kubernetes for ConfigMaps in namespace:", params.Namespace)

	configmaps, err := k8sclient.CoreV1().ConfigMaps(params.Namespace).List(metav1.ListOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to lookup ConfigMaps")
	}

	for _, configmap := range configmaps.Items {
		fmt.Println("Syncing CronJob:", configmap.ObjectMeta.Namespace, "|", configmap.ObjectMeta.Name)

		cronjob, err := generateCronJob(configmap, params)
		if err != nil {
			fmt.Fprintln(w, "Skipping CronJob:", configmap.ObjectMeta.Namespace, "|", configmap.ObjectMeta.Name, ":", err)
			continue
		}

		err = skprcronjob.Deploy(k8sclient, cronjob)
		if err != nil {
			return errors.Wrap(err, "failed to deploy CronJob")
		}
	}

	return nil
}

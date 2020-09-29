package kube

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clientset "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure" // auth for AKS clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/exec"  // auth for OIDC
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"   // auth for GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"  // auth for OIDC
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {
	Client *clientset.Clientset

	Config *restclient.Config
}

func NewClient(context string) (*KubeClient, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// if you want to change the loading rules (which files in which order), you can do so here

	configOverrides := &clientcmd.ConfigOverrides{
		CurrentContext: context,
	}
	// if you want to change override values or bind them to flags, there are methods to help you

	var config *restclient.Config
	var err error

	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err = kubeConfig.ClientConfig()
	if err != nil {
		return nil, err
	}

	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &KubeClient{
		Client: client,
		Config: config,
	}, nil
}

func (kubeClient *KubeClient) ListPods(namespace string) ([]v1.Pod, error) {
	objs, err := kubeClient.Client.CoreV1().Pods(namespace).List(metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return objs.Items, nil
}

func (kubeClient *KubeClient) GetClusterUID() (string, error) {
	ns, err := kubeClient.Client.CoreV1().Namespaces().Get("kube-system", metav1.GetOptions{})

	if err != nil {
		return "", err
	}

	return string(ns.UID), nil
}

package framework

import (
	"fmt"
	"github.com/alcideio/iskan/pkg/scan"
	"github.com/alcideio/iskan/pkg/types"
	"github.com/alcideio/iskan/pkg/vulnprovider/api"
	"time"

	"github.com/onsi/ginkgo"
	"github.com/onsi/gomega"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/apimachinery/pkg/util/wait"
	clientset "k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/azure" // auth for AKS clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/exec"  // auth for OIDC
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"   // auth for GKE clusters
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"  // auth for OIDC
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
)

// Initial pod start can be delayed O(minutes) by slow docker pulls
const PodStartTimeout = 45 * time.Second

// How often to Poll pods, nodes and claims.
const Poll = 2 * time.Second

func log(level string, format string, args ...interface{}) {
	fmt.Fprintf(ginkgo.GinkgoWriter, time.Now().Format(time.StampMilli)+": "+level+": "+format+"\n", args...)
}

func Logf(format string, args ...interface{}) {
	log("INFO", format, args...)
}

type Framework struct {
	Namespace string

	basename string

	Client clientset.Interface
}

func NewDefaultFramework(basename string) (*Framework, error) {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	// if you want to change the loading rules (which files in which order), you can do so here

	configOverrides := &clientcmd.ConfigOverrides{
		CurrentContext: "",
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

	f := &Framework{
		Client:   client,
		basename: basename,
	}

	ginkgo.BeforeEach(f.BeforeEach)
	ginkgo.AfterEach(f.AfterEach)

	return f, nil
}

func (f *Framework) BeforeEach() {
	f.Namespace = fmt.Sprint("iskan-e2e-", f.basename, "-", rand.String(5))
	_, err := f.Client.CoreV1().Namespaces().Create(&v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: f.Namespace,
		}})

	ExpectNoError(err)
	Logf("Created test namespace '%v'", f.Namespace)
}

func (f *Framework) AfterEach() {
	f.Client.CoreV1().Namespaces().Delete(f.Namespace, &metav1.DeleteOptions{})
	Logf("Deleted test namespace '%v'", f.Namespace)
}

func (f *Framework) NewImageScannerWithConfig(policy *types.Policy, config *api.VulnProvidersConfig) *scan.ImageScanner {
	var err error

	if policy == nil {
		policy = types.NewDefaultPolicy()
	}

	policy.Init()

	scanner, err := scan.NewImageScanner(policy, config)
	ExpectNoError(err)

	return scanner
}

func (f *Framework) NewImageScanner(policy *types.Policy) *scan.ImageScanner {
	var err error

	config, err := api.LoadVulnProvidersConfigFromBuffer([]byte(GlobalConfig.ApiConfigFile))
	ExpectNoError(err)

	return f.NewImageScannerWithConfig(policy, config)
}

func (f *Framework) NewClusterScanner(policy *types.Policy) *scan.ClusterScanner {
	var err error

	config, err := api.LoadVulnProvidersConfigFromBuffer([]byte(GlobalConfig.ApiConfigFile))
	ExpectNoError(err)

	return f.NewClusterScannerWithConfig(policy, config)
}

func (f *Framework) NewClusterScannerWithConfig(policy *types.Policy, config *api.VulnProvidersConfig) *scan.ClusterScanner {
	var err error

	if policy == nil {
		policy = types.NewDefaultPolicy()
	}

	policy.Init()

	scanner, err := scan.NewClusterScanner("", policy, config)
	ExpectNoError(err)

	return scanner
}

func (f *Framework) CreateImagePullSecret(name string, secret string) *v1.Secret {
	s := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: f.Namespace,
		},
		Type: v1.SecretTypeDockerConfigJson,
		Data: map[string][]byte{
			".dockerconfigjson": []byte(secret),
		},
	}

	createdSecret, err := f.Client.CoreV1().Secrets(f.Namespace).Create(s)
	ExpectNoError(err)

	return createdSecret
}

func (f *Framework) CreatePodWithContainerImage(name string, image string, imagePullSecretName string) *v1.Pod {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: f.Namespace,
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "mycontainer",
					Image: image,
				},
			},
		},
	}

	if imagePullSecretName != "" {
		pod.Spec.ImagePullSecrets = []v1.LocalObjectReference{
			{
				Name: imagePullSecretName,
			},
		}
	}

	createdPod, err := f.Client.CoreV1().Pods(f.Namespace).Create(pod)
	ExpectNoError(err)

	ExpectNoError(WaitForPodNameRunningInNamespace(f.Client, createdPod.Name, createdPod.Namespace))
	// Get the newest pod after it becomes running, some status may change after pod created, such as pod ip.
	p, err := f.Client.CoreV1().Pods(f.Namespace).Get(pod.Name, metav1.GetOptions{})
	ExpectNoError(err)
	return p
}

func (f *Framework) DeployTestImage(info *TestImageInfo) (*v1.Secret, *v1.Pod) {
	var secret *v1.Secret
	var pod *v1.Pod

	Logf("Deploying '%v' (%v)", info.Description, info.Image)

	secretName := fmt.Sprintf("secret-%v-%v", info.PullSecret, rand.String(4))
	ginkgo.By(fmt.Sprintf("creating image pull secret '%v'", secretName), func() {
		pullSecret, exist := GlobalConfig.PullSecrets[info.PullSecret]
		gomega.Expect(exist).NotTo(gomega.BeFalse())
		gomega.Expect(pullSecret).NotTo(gomega.BeNil())

		if *pullSecret != "" {
			secret = f.CreateImagePullSecret(secretName, *pullSecret)
			gomega.Expect(secret).NotTo(gomega.BeNil())
		}
	})

	podName := fmt.Sprintf("pod-%v", rand.String(4))
	ginkgo.By(fmt.Sprintf("creating pod '%v' that use the image '%v' ", podName, info.Image), func() {
		secretName := ""
		if secret != nil {
			secretName = secret.Name
		}
		pod = f.CreatePodWithContainerImage(podName, info.Image, secretName)
		gomega.Expect(pod).NotTo(gomega.BeNil())
	})

	return secret, pod
}

// Waits default amount of time (PodStartTimeout) for the specified pod to become running.
// Returns an error if timeout occurs first, or pod goes in to failed state.
func WaitForPodNameRunningInNamespace(c clientset.Interface, podName, namespace string) error {
	return waitTimeoutForPodRunningInNamespace(c, podName, namespace, PodStartTimeout)
}

func waitTimeoutForPodRunningInNamespace(c clientset.Interface, podName, namespace string, timeout time.Duration) error {
	return wait.PollImmediate(Poll, timeout, podRunning(c, podName, namespace))
}

func podRunning(c clientset.Interface, podName, namespace string) wait.ConditionFunc {
	return func() (bool, error) {
		pod, err := c.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
		if err != nil {
			return false, err
		}

		if pod.Status.PodIP == "" {
			return false, nil
		}

		switch pod.Status.Phase {
		case v1.PodRunning:
			return true, nil
		case v1.PodFailed, v1.PodSucceeded:
			return false, fmt.Errorf("pod ran to completion")
		}
		return false, nil
	}
}

func ExpectNoError(err error, explain ...interface{}) {
	expectNoErrorWithOffset(1, err, explain...)
}

// ExpectNoErrorWithOffset checks if "err" is set, and if so, fails assertion while logging the error at "offset" levels above its caller
// (for example, for call chain f -> g -> ExpectNoErrorWithOffset(1, ...) error would be logged for "f").
func expectNoErrorWithOffset(offset int, err error, explain ...interface{}) {
	if err != nil {
		klog.Errorf("Unexpected error occurred: %v", err)
	}
	gomega.ExpectWithOffset(1+offset, err).NotTo(gomega.HaveOccurred(), explain...)
}

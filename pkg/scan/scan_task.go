package scan

import (
	"fmt"
	"strings"
	"sync"

	"github.com/alcideio/iskan/pkg/registry"
	"github.com/alcideio/iskan/pkg/util"
	"github.com/alcideio/iskan/types"
	dockerref "github.com/docker/distribution/reference"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog"
)

func RegistryConfigForImage(image string, registriesConfig map[string]*types.RegistryConfig) *types.RegistryConfig {
	repo, _, _, _ := util.ParseImageName(image)

	if config, exist := registriesConfig[repo]; exist {
		klog.V(7).Infof("found exact config for %v", repo)
		return config
	}

	for regsitry, config := range registriesConfig {
		if strings.Contains(repo, regsitry) {
			klog.V(7).Infof("found config for %v", repo)
			return config
		}
	}

	kind, err := registry.DetectRegistryKind(repo)
	if err != nil {
		klog.V(5).Infof("Failed to detect %v - %v", repo, err)
	}

	if kind == registry.RegistryKind_UNKNOWN {
		klog.V(5).Infof("Failed to detect registry kind from image name - %v", repo)
	}

	return &types.RegistryConfig{
		Kind: kind,
	}
}

type ScanTaskResult struct {
	Findings map[string]*types.ImageScanResult

	ScannedPods []*v1.Pod
	SkippedPods []*v1.Pod
}

func ScanTask(pods []v1.Pod, policy *types.Policy, registriesConfig map[string]*types.RegistryConfig) (*ScanTaskResult, error) {
	scanTaskREsult := &ScanTaskResult{
		Findings:    nil,
		ScannedPods: []*v1.Pod{},
		SkippedPods: []*v1.Pod{},
	}
	errs := []error{}
	containers := sets.NewString()

	klog.V(5).Infof("[pods=%v]", len(pods))
	for i, pod := range pods {
		//Filter
		if !RuntimeScanFilter.ShouldScan(policy, &pod, "") {
			scanTaskREsult.SkippedPods = append(scanTaskREsult.SkippedPods, &pods[i])
			klog.V(5).Infof("[%v/%v] skipping", pod.Namespace, pod.Name)
			continue
		}

		podContainers := [][]v1.ContainerStatus{
			pod.Status.InitContainerStatuses,
			pod.Status.EphemeralContainerStatuses,
			pod.Status.ContainerStatuses,
		}

		var analyze bool

		analyze = false
		for _, l := range podContainers {
			for _, c := range l {
				repo, _, _, err := util.ParseImageName(c.Image)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				if !RuntimeScanFilter.ShouldScan(policy, &pod, repo) {
					klog.V(5).Infof("[%v/%v] skipping", pod.Namespace, pod.Name)
					continue
				}
				analyze = true
				containers.Insert(getImageId(c.Image, c.ImageID))
			}
		}

		if analyze {
			scanTaskREsult.ScannedPods = append(scanTaskREsult.ScannedPods, &pods[i])
		} else {
			scanTaskREsult.SkippedPods = append(scanTaskREsult.SkippedPods, &pods[i])
		}
	}

	images := containers.List()
	wg := sync.WaitGroup{}
	results := map[string]*types.ImageScanResult{}
	resLock := sync.Mutex{}

	for _, image := range images {
		wg.Add(1)
		go func(image string) {
			defer wg.Done()

			res, err := ScanImage(image, policy, RegistryConfigForImage(image, registriesConfig))
			if err != nil {
				resLock.Lock()
				errs = append(errs, err)
				resLock.Unlock()
				return
			}

			resLock.Lock()
			results[image] = res
			resLock.Unlock()

		}(image)
	}

	wg.Wait() //Wait for all tasks to complete

	scanTaskREsult.Findings = results

	return scanTaskREsult, errors.NewAggregate(errs)
}

func getImageId(image string, imageId string) string {
	if strings.HasPrefix(imageId, "docker-pullable://") {
		img := strings.TrimPrefix(imageId, "docker-pullable://")
		named, _ := dockerref.ParseNormalizedNamed(img)

		return named.String()

	} else {
		return fmt.Sprintf("%v@%v", image, imageId)
	}
}

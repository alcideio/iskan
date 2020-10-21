package scan

import (
	"fmt"
	"github.com/fatih/color"
	"k8s.io/client-go/util/flowcontrol"
	"strings"
	"sync"

	"github.com/alcideio/iskan/pkg/types"
	"github.com/alcideio/iskan/pkg/util"
	"github.com/alcideio/iskan/pkg/vulnprovider"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/klog"
)

func RegistryConfigForImage(image string, registriesConfig map[string]*types.VulnProviderConfig) *types.VulnProviderConfig {
	repo, _, _, _ := util.ParseImageName(image)

	if config, exist := registriesConfig[repo]; exist {
		klog.V(7).Infof("found exact config for %v", repo)
		return config
	}

	var c *types.VulnProviderConfig = nil
	for registry, config := range registriesConfig {
		if strings.Contains(repo, registry) || registry == "*" {
			if c == nil || c.Repository == "*" && registry != "*" {
				klog.V(5).Infof("found config for [image=%v][repo=%v][registry=%v]", image, repo, registry)
				c = config
			}
		}
	}

	if c != nil {
		return c
	}

	kind, err := vulnprovider.DetectProviderKind(repo)
	if err != nil {
		klog.V(5).Infof("Failed to detect %v - %v", repo, err)
	}

	if kind == vulnprovider.ProviderKind_UNKNOWN {
		klog.V(5).Infof("Failed to detect registry kind from image name - %v", repo)
	}

	return &types.VulnProviderConfig{
		Kind: kind,
	}
}

func ScanTask(pods []v1.Pod, policy *types.Policy, registriesConfig map[string]*types.VulnProviderConfig, flowControl flowcontrol.RateLimiter) (*types.ScanTaskResult, error) {
	scanTaskREsult := &types.ScanTaskResult{
		Findings:    nil,
		ScannedPods: []*v1.Pod{},
		SkippedPods: []*v1.Pod{},
	}
	errs := []error{}
	containers := sets.NewString()

	klog.V(5).Infof("[pods=%v]", len(pods))
	for i, _ := range pods {
		var pod *v1.Pod
		var podContainers [][]v1.ContainerStatus

		pod = &pods[i]
		//NAMESPACE Filter
		if !RuntimeScanFilter.ShouldScan(policy, pod, "") {
			scanTaskREsult.SkippedPods = append(scanTaskREsult.SkippedPods, &pods[i])
			klog.V(5).Infof("[%v/%v] skipping", pod.Namespace, pod.Name)
			continue
		}

		klog.V(5).Infof("[%v/%v] processing", pod.Namespace, pod.Name)

		podContainers = [][]v1.ContainerStatus{
			pod.Status.InitContainerStatuses,
			pod.Status.EphemeralContainerStatuses,
			pod.Status.ContainerStatuses,
		}

		var analyze bool

		analyze = false
		for i, _ := range podContainers {
			for j, _ := range podContainers[i] {
				c := podContainers[i][j]
				_, _, _, err := util.ParseImageName(c.Image)
				if err != nil {
					errs = append(errs, err)
					continue
				}

				if !RuntimeScanFilter.ShouldScan(policy, pod, c.Image) {
					klog.V(5).Infof("[%v/%v] skipping (%v)", pod.Namespace, pod.Name, c.Image)
					continue
				}

				analyze = true
				containers.Insert(util.GetImageId(c.Image, c.ImageID))
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
		regConfig := RegistryConfigForImage(image, registriesConfig)
		util.ConsolePrinter(fmt.Sprintf("Get vulnerability info for '%v' using '%v'", color.HiBlueString(image), color.HiGreenString(regConfig.Kind)))
		wg.Add(1)
		go func(image string, regConfig *types.VulnProviderConfig) {
			defer wg.Done()

			res, err := ScanImage(image, policy, regConfig, flowControl)
			if err != nil {
				resLock.Lock()
				errs = append(errs, err)
				resLock.Unlock()
				return
			}

			resLock.Lock()
			results[image] = res
			resLock.Unlock()

		}(image, regConfig)
	}

	wg.Wait() //Wait for all tasks to complete

	scanTaskREsult.Findings = results

	return scanTaskREsult, errors.NewAggregate(errs)
}

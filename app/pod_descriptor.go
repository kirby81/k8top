package app

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodDescriptor struct {
	Name             string
	NbContainerReady int
	Totalcontainer   int
	Restart          int32
	Status           string
	Age              time.Duration
}

func (pd *PodDescriptor) GetPodReadiness() string {
	return fmt.Sprintf("%d/%d", pd.NbContainerReady, pd.Totalcontainer)
}

func GetPods(namespace string) []PodDescriptor {
	k8Client := Connect()
	result := make([]PodDescriptor, 0)

	pods, err := k8Client.CoreV1().Pods("default").List(metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}

	for _, pod := range pods.Items {
		delta := time.Since(pod.GetCreationTimestamp().Time)
		nbContainerReady := 0
		for _, c := range pod.Status.ContainerStatuses {
			if c.Ready {
				nbContainerReady++
			}
		}
		result = append(result, PodDescriptor{
			Name:             pod.GetName(),
			NbContainerReady: nbContainerReady,
			Totalcontainer:   len(pod.Status.ContainerStatuses),
			Restart:          pod.Status.ContainerStatuses[0].RestartCount,
			Status:           string(pod.Status.Phase),
			Age:              delta.Truncate(time.Second),
		})
	}

	return result
}

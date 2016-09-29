package kuberang

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

type KubeOutput struct {
	Success     bool
	CombinedOut string
	RawOut      []byte
}

func RunKubectl(args ...string) KubeOutput {
	kubeCmd := exec.Command("kubectl", args...)
	bytes, err := kubeCmd.CombinedOutput()
	if err != nil {
		return KubeOutput{
			Success:     false,
			CombinedOut: string(bytes),
			RawOut:      bytes,
		}
	}
	return KubeOutput{
		Success:     true,
		CombinedOut: string(bytes),
		RawOut:      bytes,
	}
}

func RunGetService(svcName string) KubeOutput {
	return RunKubectl("get", "service", svcName, "-o", "json")
}

func RunGetPodByImage(name string) KubeOutput {
	return RunKubectl("get", "deployment", name, "-o", "json")
}

func RunGetDeployment(name string) KubeOutput {
	return RunKubectl("get", "deployment", name, "-o", "json")
}

func RunPod(name string, image string, count int64) KubeOutput {
	return RunKubectl("run", name, "--image="+image, "--replicas="+strconv.FormatInt(count, 10), "-o", "json")
}

func (ko KubeOutput) ObservedReplicaCount() int64 {
	resp := DeploymentResponse{}
	json.Unmarshal(ko.RawOut, &resp)
	return resp.Status.AvaiableReplicas
}

type DeploymentResponse struct {
	Status struct {
		AvaiableReplicas int64 `json:"availableReplicas"`
	} `json:"status"`
}

func (ko KubeOutput) ServiceCluserIP() string {
	resp := ServiceResponse{}
	json.Unmarshal(ko.RawOut, &resp)
	return resp.Spec.ClusterIP
}

type ServiceResponse struct {
	Spec struct {
		ClusterIP string `json:"clusterIP"`
	} `json:"spec"`
}

func (ko KubeOutput) PodIPs() []string {
	//In Scala, this code would be gorgeous. In Golang, it's a blood blister
	resp := PodsResponse{}
	if err := json.Unmarshal(ko.RawOut, &resp); err != nil {
		fmt.Println(err)
	}
	podIPs := make([]string, len(resp.Items))
	for i, item := range resp.Items {
		podIPs[i] = item.Status.PodIP
	}
	return podIPs
}

func (ko KubeOutput) FirstPodName() string {
	resp := PodsResponse{}
	if err := json.Unmarshal(ko.RawOut, &resp); err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(ko.RawOut, &resp)
	if len(resp.Items) < 1 {
		return ""
	}
	return resp.Items[0].Metadata.Name
}

type PodsResponse struct {
	Items []struct {
		Metadata struct {
			Name string `json:"name"`
		} `json:"metadata"`
		Status struct {
			PodIP string `json:"podIP"`
		} `json:"status"`
	} `json:"items"`
}

package diverts

import (
	"encoding/json"
	"fmt"
	"strings"

	okLabels "github.com/okteto/okteto/pkg/k8s/labels"
	"github.com/okteto/okteto/pkg/model"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type divertServiceModification struct {
	ProxyProtocol      string `json:"proxy_port"`
	OriginalPort       string `json:"original_port"`
	OriginalTargetPort string `json:"original_target_port"`
}

// DivertName returns the name of the diverted version of a given resource
func DivertName(username, name string) string {
	return fmt.Sprintf("%s-%s", username, name)
}

func translateDeployment(username string, dev *model.Dev, d *appsv1.Deployment) *appsv1.Deployment {
	result := d.DeepCopy()
	result.UID = ""
	result.Name = DivertName(username, dev.Name)
	result.Labels = map[string]string{model.OktetoDivertLabel: username}
	result.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: map[string]string{
			model.OktetoDivertLabel:      username,
			okLabels.InteractiveDevLabel: result.Name,
		},
	}
	result.Spec.Template.Labels = map[string]string{
		model.OktetoDivertLabel:      username,
		okLabels.InteractiveDevLabel: result.Name,
	}
	if result.Annotations == nil {
		result.Annotations = map[string]string{}
	}
	result.Annotations[model.OktetoAutoCreateAnnotation] = model.OktetoUpCmd
	result.ResourceVersion = ""
	return result
}

func translateService(username string, dev *model.Dev, d *appsv1.Deployment, s *apiv1.Service) (*apiv1.Service, error) {
	result := s.DeepCopy()
	result.UID = ""
	result.Name = DivertName(username, dev.Name)
	result.Labels = map[string]string{model.OktetoDivertLabel: username}
	if s.Annotations != nil {
		modification := s.Annotations[model.OktetoDivertServiceModificationAnnotation]
		if modification != "" {
			var dsm divertServiceModification
			if err := json.Unmarshal([]byte(modification), &dsm); err != nil {
				return nil, fmt.Errorf("bad divert service modification: %s", modification)
			}
			for i := range result.Spec.Ports {
				if strings.Compare(fmt.Sprint(result.Spec.Ports[i].Port), dsm.OriginalPort) == 0 {
					result.Spec.Ports[i].TargetPort = intstr.FromString(dsm.OriginalTargetPort)
				}
			}
		}
	}
	delete(result.Annotations, okLabels.OktetoAutoIngressAnnotation)
	delete(result.Annotations, model.OktetoDivertServiceModificationAnnotation)
	result.Spec.Selector = map[string]string{
		model.OktetoDivertLabel:      username,
		okLabels.InteractiveDevLabel: d.Name,
	}
	result.ResourceVersion = ""
	return result, nil
}

func translateIngress(username string, dev *model.Dev, i *networkingv1.Ingress) *networkingv1.Ingress {
	result := i.DeepCopy()
	result.UID = ""
	result.Name = DivertName(username, dev.Name)
	result.Labels = map[string]string{model.OktetoDivertLabel: username}
	if result.Annotations == nil {
		result.Annotations = map[string]string{}
	}
	if host := result.Annotations[okLabels.OktetoIngressAutoGenerateHost]; host != "" {
		if host != "true" {
			result.Annotations[okLabels.OktetoIngressAutoGenerateHost] = fmt.Sprintf("%s-%s", username, host)
		}
	} else {
		result.Annotations[okLabels.OktetoIngressAutoGenerateHost] = "true"
	}
	result.ResourceVersion = ""
	return result
}

func translateDivertCRD(username string, dev *model.Dev, s *apiv1.Service, i *networkingv1.Ingress) *Divert {
	return &Divert{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Divert",
			APIVersion: "weaver.okteto.com/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: dev.Namespace,
		},
		Spec: DivertSpec{
			Ingress: IngressDivertSpec{
				Name:      i.Name,
				Namespace: dev.Namespace,
				Value:     username,
			},
			FromService: ServiceDivertSpec{
				Name:      dev.Divert.Service,
				Namespace: dev.Namespace,
				Port:      dev.Divert.Port,
			},
			ToService: ServiceDivertSpec{
				Name:      s.Name,
				Namespace: dev.Namespace,
				Port:      dev.Divert.Port,
			},
			Deployment: DeploymentDivertSpec{
				Name:      dev.Name,
				Namespace: dev.Namespace,
			},
		},
	}
}

func translateDev(username string, dev *model.Dev, d *appsv1.Deployment) {
	dev.Name = d.Name
	dev.Labels = nil
}

package k8s

import (
	"fmt"
	"path/filepath"

	"github.com/okteto/okteto/pkg/k8s/secrets"
	"github.com/okteto/okteto/pkg/model"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	resource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	oktetoSyncLabel    = "syncthing.okteto.com"
	syncImageTag       = "okteto/syncthing:1.2.2"
	syncTCPPort        = 22000
	syncGUIPort        = 8384
	oktetoContainer    = "okteto"
	oktetoSecretVolume = "okteto-secret"
	oktetoMount        = "/var/okteto"
)

var (
	devReplicas int32 = 1
)

func translate(dev *model.Dev, d *appsv1.Deployment, c *apiv1.Container) *appsv1.StatefulSet {
	initContainer := translateInitContainer(dev)

	reqMem, _ := resource.ParseQuantity("64Mi")
	reqCPU, _ := resource.ParseQuantity("50m")
	limMem, _ := resource.ParseQuantity("256Mi")
	limCPU, _ := resource.ParseQuantity("500m")
	ss := &appsv1.StatefulSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      dev.GetStatefulSetName(),
			Namespace: dev.Namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			ServiceName:         dev.GetStatefulSetName(),
			PodManagementPolicy: appsv1.ParallelPodManagement,
			Replicas:            &devReplicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					oktetoSyncLabel: dev.Name,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						oktetoSyncLabel: dev.Name,
					},
				},
				Spec: apiv1.PodSpec{
					TerminationGracePeriodSeconds: &devTerminationGracePeriodSeconds,
					SecurityContext:               d.Spec.Template.Spec.SecurityContext,
					ImagePullSecrets:              d.Spec.Template.Spec.ImagePullSecrets,
					InitContainers:                []apiv1.Container{*initContainer},
					Containers: []apiv1.Container{
						apiv1.Container{
							Name:            oktetoContainer,
							Image:           syncImageTag,
							ImagePullPolicy: c.ImagePullPolicy,
							Resources: apiv1.ResourceRequirements{
								Requests: apiv1.ResourceList{
									apiv1.ResourceMemory: reqMem,
									apiv1.ResourceCPU:    reqCPU,
								},
								Limits: apiv1.ResourceList{
									apiv1.ResourceMemory: limMem,
									apiv1.ResourceCPU:    limCPU,
								},
							},
							VolumeMounts: []apiv1.VolumeMount{
								apiv1.VolumeMount{
									Name:      oktetoSecretVolume,
									MountPath: "/var/syncthing/secret/",
								},
								apiv1.VolumeMount{
									Name:      dev.GetVolumeTemplateName(0),
									MountPath: oktetoMount,
								},
							},
							Ports: []apiv1.ContainerPort{
								apiv1.ContainerPort{
									ContainerPort: syncGUIPort,
								},
								apiv1.ContainerPort{
									ContainerPort: syncTCPPort,
								},
							},
						},
					},
					Volumes: []apiv1.Volume{
						apiv1.Volume{
							Name: oktetoSecretVolume,
							VolumeSource: apiv1.VolumeSource{
								Secret: &apiv1.SecretVolumeSource{
									SecretName: secrets.GetSecretName(dev),
								},
							},
						},
					},
					NodeSelector: d.Spec.Template.Spec.NodeSelector,
					Affinity:     d.Spec.Template.Spec.Affinity,
				},
			},
		},
	}

	if dev.SecurityContext != nil {
		ss.Spec.Template.Spec.SecurityContext = &apiv1.PodSecurityContext{
			RunAsUser:  dev.SecurityContext.RunAsUser,
			RunAsGroup: dev.SecurityContext.RunAsGroup,
			FSGroup:    dev.SecurityContext.FSGroup,
		}
	}

	ss.Spec.VolumeClaimTemplates = translateVolumeClaimTemplates(dev)
	return ss
}

func translateInitContainer(dev *model.Dev) *apiv1.Container {
	reqMem, _ := resource.ParseQuantity("16Mi")
	reqCPU, _ := resource.ParseQuantity("50m")
	limMem, _ := resource.ParseQuantity("16Mi")
	limCPU, _ := resource.ParseQuantity("50m")
	source := filepath.Join(dev.MountPath, "*")

	target := "/okteto/init"
	if dev.SubPath != "" {
		target = filepath.Join(target, dev.SubPath)
	}
	c := &apiv1.Container{
		Name:    model.OktetoInitContainer,
		Image:   dev.Image,
		Command: []string{"sh", "-c", fmt.Sprintf("(ls -A /okteto/init | grep -v lost+found || mkdir -p %s && cp -Rf %s %s); touch /okteto/init/%s", target, source, target, dev.DevPath)},
		Resources: apiv1.ResourceRequirements{
			Requests: apiv1.ResourceList{
				apiv1.ResourceMemory: reqMem,
				apiv1.ResourceCPU:    reqCPU,
			},
			Limits: apiv1.ResourceList{
				apiv1.ResourceMemory: limMem,
				apiv1.ResourceCPU:    limCPU,
			},
		},
		VolumeMounts: []apiv1.VolumeMount{
			apiv1.VolumeMount{
				Name:      dev.GetVolumeTemplateName(0),
				MountPath: "/okteto/init",
			},
		},
	}

	for i, v := range dev.Volumes {
		c.VolumeMounts = append(
			c.VolumeMounts,
			apiv1.VolumeMount{
				Name:      dev.GetVolumeTemplateName(i + 1),
				MountPath: v,
			},
		)
	}

	return c
}

func translateVolumeClaimTemplates(dev *model.Dev) []apiv1.PersistentVolumeClaim {
	quantDisk, _ := resource.ParseQuantity("10Gi")
	result := []apiv1.PersistentVolumeClaim{}
	for i := 0; i <= len(dev.Volumes); i++ {
		result = append(
			result,
			apiv1.PersistentVolumeClaim{
				ObjectMeta: metav1.ObjectMeta{
					Name: dev.GetVolumeTemplateName(i),
				},
				Spec: apiv1.PersistentVolumeClaimSpec{
					AccessModes: []apiv1.PersistentVolumeAccessMode{apiv1.ReadWriteOnce},
					Resources: apiv1.ResourceRequirements{
						Requests: apiv1.ResourceList{
							"storage": quantDisk,
						},
					},
				},
			},
		)
	}
	return result
}

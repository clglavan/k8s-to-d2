package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func writeObjectInfo(file *os.File, objectPath string, title string, fields map[string]string) {
	file.WriteString(fmt.Sprintf(objectPath + ".extradata: |md\n"))
	file.WriteString(fmt.Sprintf(title + "\n"))

	for fieldName, fieldValue := range fields {
		file.WriteString(fmt.Sprintf("- %s = %s\n", fieldName, fieldValue))
	}

	file.WriteString(fmt.Sprintf("|\n"))
}

func main() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}

	// move this to parameter
	allowedNamespaces := []string{"default"}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	file, err := os.Create("output.in")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("namespace: { grid-columns: 1 }\n"))
	// file.WriteString(fmt.Sprintf("namespace.grid-columns: 1\n"))

	file.WriteString(fmt.Sprintf(`classes: {
		nginx: {label: Nginx; icon: https://www.vectorlogo.zone/logos/nginx/nginx-icon.svg}
		kibana: {label: Kibana; icon: https://www.vectorlogo.zone/logos/elasticco_kibana/elasticco_kibana-icon.svg}
		elasticsearch: {label: ElasticSearch; icon: https://www.vectorlogo.zone/logos/elastic/elastic-icon.svg}
		oauth2-proxy: {label: OAuth2 Proxy; icon: https://github.com/oauth2-proxy/oauth2-proxy/raw/master/docs/static/img/logos/OAuth2_Proxy_icon.svg}
		kyverno: {label: Kyverno; icon: https://github.com/cncf/artwork/raw/main/projects/kyverno/icon/color/kyverno-icon-color.svg}
		flux: {label: Flux CD; icon: https://github.com/cncf/artwork/raw/main/projects/flux/icon/color/flux-icon-color.svg}
		fluent-bit: {label: Fluent Bit; icon: https://asset.brandfetch.io/idxVhszl6V/id-eYdnm_0.svg}
		prometheus: {label: Prometheus; icon: https://github.com/cncf/artwork/raw/main/projects/prometheus/icon/color/prometheus-icon-color.svg}
		alertmanager: {label: Alertmanager; icon: https://github.com/cncf/artwork/raw/main/projects/prometheus/icon/color/prometheus-icon-color.svg}
		grafana: {label: Grafana; icon: https://www.vectorlogo.zone/logos/grafana/grafana-icon.svg}
		thanos: {label: Thanos; icon: https://github.com/cncf/artwork/raw/main/projects/thanos/icon/color/thanos-icon-color.svg}
		boundary: {label: Bondary; icon: https://www.svgrepo.com/download/448275/boundary.svg}
		falco: {label: Falco; icon: https://www.vectorlogo.zone/logos/falco/falco-icon.svg}
		cert-manager: {label: Cert Manager; icon: https://github.com/cncf/artwork/raw/main/projects/cert-manager/icon/color/cert-manager-icon-color.svg}
		kubecost: {label: Kubecost; icon: https://github.com/cncf/landscape/raw/master/hosted_logos/kubecost.svg}
		velero: {label: Velero; icon: https://github.com/cncf/landscape/raw/master/hosted_logos/project-velero.svg}
		percona: {label: Percona; icon: https://docs.percona.com/percona-software-repositories/_images/percona-logo.svg}
		symfony: {label: Symfony; icon: https://www.vectorlogo.zone/logos/symfony/symfony-icon.svg}
		dotnet: {label: .NET; icon: https://github.com/dotnet/brand/raw/main/logo/dotnet-logo.svg}
		wordpress: {label: WordPress; icon: https://www.vectorlogo.zone/logos/wordpress/wordpress-icon.svg}
		rabbitmq: {label: RabbitMQ; icon: https://www.vectorlogo.zone/logos/rabbitmq/rabbitmq-icon.svg}
		mongodb: {label: MongoDB; icon: https://www.vectorlogo.zone/logos/mongodb/mongodb-icon.svg}
		k8s: {label: Kubernetes; icon: https://github.com/kubernetes/kubernetes/raw/master/logo/logo.svg}
		k8s-ns: {label: Namespace; icon: https://raw.githubusercontent.com/kubernetes/community/master/icons/svg/resources/labeled/ns.svg}
		k8s-svc: {label: Service; icon: https://github.com/kubernetes/community/raw/master/icons/svg/resources/labeled/svc.svg}
		k8s-ingress: {label: Ingress; icon: https://github.com/kubernetes/community/raw/master/icons/svg/resources/labeled/ing.svg}
		k8s-deployment: {label: Deployment; icon: https://raw.githubusercontent.com/kubernetes/community/master/icons/svg/resources/labeled/deploy.svg}
		k8s-configmap: {label: ConfigMap; icon: https://raw.githubusercontent.com/kubernetes/community/master/icons/svg/resources/labeled/cm.svg}
		k8s-secret: {label: Secret; icon: https://raw.githubusercontent.com/kubernetes/community/master/icons/svg/resources/labeled/secret.svg}
		k8s-hpa: {label: HPA; icon: https://raw.githubusercontent.com/kubernetes/community/master/icons/svg/resources/labeled/hpa.svg}
	  }` + "\n"))

	for _, namespace := range allowedNamespaces {
		file.WriteString(fmt.Sprintf("namespace.'%s': { grid-columns: 1 }\n", namespace))
		file.WriteString(fmt.Sprintf("namespace.'%s'.deployment: { grid-columns: 3 }\n", namespace))
		file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `': {
			class: k8s-ns 
			label: Namespace ` + namespace + `
		}` + "\n"))

		ns, err := clientset.CoreV1().Namespaces().Get(context.TODO(), namespace, v1.GetOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}
		annotations := make([]string, 0, len(ns.GetObjectMeta().GetAnnotations()))
		for k, v := range ns.ObjectMeta.Annotations {
			annotations = append(annotations, k+"="+v)
		}

		labels := make([]string, 0, len(ns.GetObjectMeta().GetLabels()))
		for k, v := range ns.ObjectMeta.Labels {
			labels = append(labels, k+"="+v)
		}

		file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.extradata: |md
				# Namespace ` + namespace + `
				 ` + strings.Join(labels, " , ") + `
				 ` + strings.Join(annotations, " , ") + `
				|` + "\n"))

		ingresses, err := clientset.NetworkingV1().Ingresses(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Println(err)
			return
		}

		// Resources below are in the scope of an ingress
		for _, ingress := range ingresses.Items {
			// jsonIngress, err := json.MarshalIndent(ingress, "", "  ")
			if err != nil {
				// handle error
			}

			file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.ingress.'` + ingress.Name + `': {
				class: k8s-ingress
				label: Ingress ` + ingress.Name + `
			}` + "\n"))

			// fmt.Print(ingress.Status.LoadBalancer.Ingress[0].IP)
			lbIngress := make([]string, 0, len(ingress.Status.LoadBalancer.Ingress))
			for _, v := range ingress.Status.LoadBalancer.Ingress {
				// fmt.Println(k, v.IP)
				lbIngress = append(lbIngress, v.IP)
			}
			file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.ingress.'` + ingress.Name + `'.extradata: |md
			# Ingress ` + ingress.Name + `
			 ` + strings.Join(lbIngress, " , ") + `
			|` + "\n"))
			// fmt.Println(string(jsonIngress))

			// Resources below are in the scope of an ingress rule
			for _, rule := range ingress.Spec.Rules {
				file.WriteString(fmt.Sprintf("namespace.'%s'.ingress.'%s' -> namespace.'%s'.ingress.rule.'%s'\n", namespace, ingress.Name, namespace, rule.Host))

				// Resources below are in the scope of an ingress path
				for _, path := range rule.HTTP.Paths {
					service, err := clientset.CoreV1().Services(namespace).Get(context.TODO(), path.Backend.Service.Name, metav1.GetOptions{})

					file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.service.'` + service.Name + `': {
						class: k8s-svc
						label: Service ` + service.Name + `
					}` + "\n"))
					file.WriteString(fmt.Sprintf("namespace.'%s'.ingress.rule.'%s' -> namespace.'%s'.service.'%s': HTTP path '%s'\n", namespace, rule.Host, namespace, service.Name, path.Path))

					// file.WriteString(fmt.Sprintf("Namespace: %s, Ingress: %s, Service: %s, Path: %s, Backend: %s\n", namespace, ingress.Name, service.Name, path.Path, path.Backend.Service.Name))

					deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
					if err != nil {
						fmt.Println(err)
						continue
					}
					// fmt.Printf("Found %d deployments in namespace %s\n", len(deployments.Items), namespace)
					// fmt.Println("Service selector:", service.Spec.Selector)
					for _, deployment := range deployments.Items {
						if containsLabels(deployment.Spec.Selector.MatchLabels, service.Spec.Selector) {
							var ports []string
							for _, port := range service.Spec.Ports {
								ports = append(ports, fmt.Sprintf("%v", port.TargetPort.String()))
							}
							portsString := strings.Join(ports, ", ")

							// fmt.Println("Ports:", portsString)

							file.WriteString(fmt.Sprintf("namespace.'%s'.service.'%s' -> namespace.'%s'.deployment.'%s': Port '%s' \n", namespace, service.Name, namespace, deployment.Name, portsString))
							file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.deployment.'` + deployment.Name + `': {
								class: k8s-deployment

							}` + "\n"))
						}
					}

					statefulSets, err := clientset.AppsV1().StatefulSets(namespace).List(context.TODO(), metav1.ListOptions{})
					if err != nil {
						fmt.Println(err)
						continue
					}

					for _, statefulSet := range statefulSets.Items {
						// file.WriteString(fmt.Sprintf("Service: %s, StatefulSet: %s\n", service.Name, statefulSet.Name))
						if containsLabels(statefulSet.Spec.Selector.MatchLabels, service.Spec.Selector) {

							file.WriteString(fmt.Sprintf("namespace.'%s'.service.'%s' -> namespace.'%s'.statefulSet.'%s' \n", namespace, service.Name, namespace, statefulSet.Name))
						}
					}

					daemonSets, err := clientset.AppsV1().DaemonSets(namespace).List(context.TODO(), metav1.ListOptions{})
					if err != nil {
						fmt.Println(err)
						continue
					}

					for _, daemonSet := range daemonSets.Items {
						if containsLabels(daemonSet.Spec.Selector.MatchLabels, service.Spec.Selector) {

							file.WriteString(fmt.Sprintf("namespace.'%s'.service.'%s' -> namespace.'%s'.daemonSet.'%s' \n", namespace, service.Name, namespace, daemonSet.Name))
						}
					}
				}
			}
		}

		deployments, err := clientset.AppsV1().Deployments(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			fmt.Println(err)
			continue
		}

		// file.WriteString(fmt.Sprintf("namespace.'%s'.deployment.grid-columns: 5\n", namespace))
		// fmt.Printf("Found %d deployments in namespace %s\n", len(deployments.Items), namespace)
		// fmt.Println("Service selector:", service.Spec.Selector)
		for _, deployment := range deployments.Items {
			// ConfigMaps and Secrets
			file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.deployment.'` + deployment.Name + `': {
				class: k8s-deployment
			}` + "\n"))

			replicas := ""
			if deployment.Spec.Replicas != nil {
				replicas = strconv.FormatInt(int64(*deployment.Spec.Replicas), 10)
			}

			writeObjectInfo(file, "namespace.'"+namespace+"'.deployment.'"+deployment.Name+"'",
				"# Deployment "+deployment.Name,
				map[string]string{"Replicas": replicas,
					"Strategy":                deployment.Spec.Strategy.RollingUpdate.String(),
					"MinReadySeconds":         strconv.FormatInt(int64(deployment.Spec.MinReadySeconds), 10),
					"ProgressDeadlineSeconds": strconv.FormatInt(int64(*deployment.Spec.ProgressDeadlineSeconds), 10),
				})

			file.WriteString(fmt.Sprintf("namespace.'%s'.deployment.'%s'.grid-columns: 1\n", namespace, deployment.Name))
			for _, volume := range deployment.Spec.Template.Spec.Volumes {
				// fmt.Println(volume.Name)
				if volume.ConfigMap != nil {
					file.WriteString(fmt.Sprintf("namespace.'%s'.deployment.'%s'.configMap.'%s'\n", namespace, deployment.Name, volume.ConfigMap.Name))
					file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.deployment.'` + deployment.Name + `'.configMap.'` + volume.ConfigMap.Name + `': {
						class: k8s-configmap
					}` + "\n"))
				}
				if volume.Secret != nil {
					file.WriteString(fmt.Sprintf("namespace.'%s'.deployment.'%s'.secret.'%s'\n", namespace, deployment.Name, volume.Secret.SecretName))
					file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.deployment.'` + deployment.Name + `'.secret.'` + volume.Secret.SecretName + `': {
						class: k8s-secret
						label: Secret ` + volume.Secret.SecretName + `
						
					}` + "\n"))
				}
			}
			for _, container := range deployment.Spec.Template.Spec.Containers {
				for _, envFrom := range container.EnvFrom {
					if envFrom.ConfigMapRef != nil {
						file.WriteString(fmt.Sprintf("namespace.'%s'.deployment.'%s'.configMap.'%s'\n", namespace, deployment.Name, envFrom.ConfigMapRef.Name))
						file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.deployment.'` + deployment.Name + `'.configMap.'` + envFrom.ConfigMapRef.Name + `': {
							class: k8s-configmap							
						}` + "\n"))
						configMap, err := clientset.CoreV1().ConfigMaps(namespace).Get(context.TODO(), envFrom.ConfigMapRef.Name, metav1.GetOptions{})
						if err != nil {
							// handle error
						}
						annotations := make([]string, 0, len(configMap.GetObjectMeta().GetAnnotations()))
						for k, v := range configMap.ObjectMeta.Annotations {
							annotations = append(annotations, k+"="+v)
						}
						labels := make([]string, 0, len(configMap.GetObjectMeta().GetLabels()))
						for k, v := range configMap.ObjectMeta.Labels {
							labels = append(labels, k+"="+v)
						}
						configMapKeys := []string{}
						for key := range configMap.Data {
							configMapKeys = append(configMapKeys, key)
						}
						writeObjectInfo(file, "namespace.'"+namespace+"'.deployment.'"+deployment.Name+"'.configMap.'"+envFrom.ConfigMapRef.Name+"'",
							"# ConfigMap "+envFrom.ConfigMapRef.Name,
							map[string]string{
								"Items":       fmt.Sprintf("%v", configMapKeys),
								"Labels":      fmt.Sprintf("%v", labels),
								"Annotations": fmt.Sprintf("%v", annotations)})

					}
					if envFrom.SecretRef != nil {
						file.WriteString(fmt.Sprintf("namespace.'%s'.deployment.'%s'.secret.'%s'\n", namespace, deployment.Name, envFrom.SecretRef.Name))
						file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.deployment.'` + deployment.Name + `'.secret.'` + envFrom.SecretRef.Name + `': {
							class: k8s-secret
							label: Secret ` + envFrom.SecretRef.Name + `
							
						}` + "\n"))
						if envFrom.SecretRef != nil {
							file.WriteString(fmt.Sprintf("namespace.'%s'.deployment.'%s'.secret.'%s'\n", namespace, deployment.Name, envFrom.SecretRef.Name))
							file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.deployment.'` + deployment.Name + `'.secret.'` + envFrom.SecretRef.Name + `': {
								class: k8s-secret
								label: Secret ` + envFrom.SecretRef.Name + `
								
							}` + "\n"))
							secret, err := clientset.CoreV1().Secrets(namespace).Get(context.TODO(), envFrom.SecretRef.Name, metav1.GetOptions{})
							if err != nil {
								// handle error
							}
							annotations := make([]string, 0, len(secret.GetObjectMeta().GetAnnotations()))
							for k, v := range secret.ObjectMeta.Annotations {
								annotations = append(annotations, k+"="+v)
							}
							labels := make([]string, 0, len(secret.GetObjectMeta().GetLabels()))
							for k, v := range secret.ObjectMeta.Labels {
								labels = append(labels, k+"="+v)
							}
							secretKeys := []string{}
							for key := range secret.Data {
								secretKeys = append(secretKeys, key)
							}
							writeObjectInfo(file, "namespace.'"+namespace+"'.deployment.'"+deployment.Name+"'.secret.'"+envFrom.SecretRef.Name+"'",
								"# Secret "+envFrom.SecretRef.Name,
								map[string]string{
									"Items":       fmt.Sprintf("%v", secretKeys),
									"Labels":      fmt.Sprintf("%v", labels),
									"Annotations": fmt.Sprintf("%v", annotations)})
						}
					}
				}
			}

			// HorizontalPodAutoscalers
			hpas, err := clientset.AutoscalingV2().HorizontalPodAutoscalers(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, hpa := range hpas.Items {
				if hpa.Spec.ScaleTargetRef.Name == deployment.Name {
					file.WriteString(fmt.Sprintf("namespace.'%s'.deployment.'%s'.hpa.'%s'\n", namespace, deployment.Name, hpa.Name))
					file.WriteString(fmt.Sprintf(`namespace.'` + namespace + `'.deployment.'` + deployment.Name + `'.hpa.'` + hpa.Name + `': {
						class: k8s-hpa
						label: HPA ` + hpa.Name + `
						
					}` + "\n"))
				}
				annotations := make([]string, 0, len(hpa.GetObjectMeta().GetAnnotations()))
				for k, v := range hpa.ObjectMeta.Annotations {
					annotations = append(annotations, k+"="+v)
				}
				labels := make([]string, 0, len(hpa.GetObjectMeta().GetLabels()))
				for k, v := range hpa.ObjectMeta.Labels {
					labels = append(labels, k+"="+v)
				}
				// secretKeys := []string{}
				// for key := range secret.Data {
				// 	secretKeys = append(secretKeys, key)
				// }
				writeObjectInfo(file, "namespace.'"+namespace+"'.deployment.'"+deployment.Name+"'.hpa.'"+hpa.Name+"'",
					"# HPA "+hpa.Name,
					map[string]string{
						// "Items":       fmt.Sprintf("%v", secretKeys),
						"Labels":      fmt.Sprintf("%v", labels),
						"Annotations": fmt.Sprintf("%v", annotations)})
			}

			// PodDisruptionBudgets
			pdbs, err := clientset.PolicyV1().PodDisruptionBudgets(namespace).List(context.TODO(), metav1.ListOptions{})
			if err != nil {
				fmt.Println(err)
				continue
			}
			for _, pdb := range pdbs.Items {
				if containsLabels(pdb.Spec.Selector.MatchLabels, deployment.Spec.Template.Labels) {
					file.WriteString(fmt.Sprintf("namespace.'%s'.deployment.'%s'.pdb.'%s'\n", namespace, deployment.Name, pdb.Name))
				}
			}

			// VolumeMounts PVCs
			for _, volume := range deployment.Spec.Template.Spec.Volumes {
				if volume.PersistentVolumeClaim != nil {
					file.WriteString(fmt.Sprintf("namespace.'%s'.deployment.'%s'.pvc.'%s'\n", namespace, deployment.Name, volume.PersistentVolumeClaim.ClaimName))
				}
			}
		}

	}
}

func containsLabels(deploymentMatchLabels map[string]string, serviceSelector map[string]string) bool {
	for key, value := range serviceSelector {
		if deploymentMatchLabels[key] != value {
			return false
		}
	}
	return true
}

## k8s Deployment

Structure
- k8s/: This is the main directory where all Kubernetes-related files are stored.
- namespaces/: Contains YAML files for defining any custom namespaces (like cryptkeeper-namespace.yaml).
- deployments/: Stores all your Deployment definitions. Each deployment has its own YAML file, making it easier to manage and update.
- services/: Contains Service definitions corresponding to each Deployment, again each in its own file.
- configmaps/: This folder is for ConfigMap definitions, separating the configuration files from the deployment and service definitions.
- persistent-volume-claims/: Contains Persistent Volume Claim definitions for your stateful applications (like PostgreSQL and Grafana).
- other-resources/: If you have additional resource types (like secrets, roles, etc.), they can go here.

README.md: Both at the root of your repository and in the k8s/ directory, providing documentation on how to deploy your services, usage instructions, and any necessary context.



```
cryptkeeper/
├── k8s/
│   ├── namespaces/
│   │   └── cryptkeeper-namespace.yaml
│   ├── deployments/
│   │   ├── postgres-deployment.yaml
│   │   ├── redis-deployment.yaml
│   │   ├── spicedb-deployment.yaml
│   │   ├── kafka-deployment.yaml
│   │   ├── zookeeper-deployment.yaml
│   │   ├── prometheus-deployment.yaml
│   │   ├── grafana-deployment.yaml
│   │   ├── cryptkeeper-backend-deployment.yaml
│   │   └── cryptkeeper-ui-deployment.yaml
│   ├── services/
│   │   ├── postgres-service.yaml
│   │   ├── redis-service.yaml
│   │   ├── spicedb-service.yaml
│   │   ├── kafka-service.yaml
│   │   ├── zookeeper-service.yaml
│   │   ├── prometheus-service.yaml
│   │   ├── grafana-service.yaml
│   │   ├── cryptkeeper-backend-service.yaml
│   │   └── cryptkeeper-ui-service.yaml
│   ├── configmaps/
│   │   ├── init-sql-configmap.yaml
│   │   ├── prometheus-configmap.yaml
│   │   └── cryptkeeper-configmap.yaml
│   ├── persistent-volume-claims/
│   │   ├── postgres-pvc.yaml
│   │   └── grafana-pvc.yaml
│   └── README.md
```


### How to?
`kubectl apply -f k8s/` to apply the configurations
## Kubernetes Deployment

### Prerequisites

- Kubernetes cluster (local or remote)
- kubectl configured to interact with your cluster

### Building the application
```sh
cd cryptkeeper-ui && docker build -t cryptkeeper-cryptkeeper-ui:latest -f Dockerfile .
docker build -t cryptkeeper-backend:latest -f Dockerfile .
```

Tag for k8s
```sh
docker tag cryptkeeper-cryptkeeper-ui:latest localhost:5000/cryptkeeper-cryptkeeper-ui:latest
docker tag cryptkeeper-backend:latest localhost:5000/cryptkeeper-backend:latest

docker run -d -p 5000:5000 --restart=always --name registry registry:2

docker push localhost:5000/cryptkeeper-cryptkeeper-ui:latest
docker push localhost:5000/cryptkeeper-backend:latest

```


### Deploying the Application

1. **Create the Namespace**:
   ```sh
   kubectl apply -f k8s/namespace.yaml
   kubectl apply -f k8s/resourcequota.yaml
   ```

2. **Create the ConfigMap**:
   ```sh
   kubectl apply -f k8s/configmap.yaml
   ```

3. **Deploy Postgres**:
   ```sh
   kubectl apply -f k8s/postgres/postgres-deployment.yaml
   kubectl apply -f k8s/postgres/postgres-service.yaml
   ```

4. **Deploy CryptKeeper Backend**:
   ```sh
   kubectl apply -f k8s/backend/cryptkeeper-backend-deployment.yaml
   kubectl apply -f k8s/backend/cryptkeeper-backend-service.yaml
   ```

5. **Deploy CryptKeeper UI**:
   ```sh
   kubectl apply -f k8s/ui/cryptkeeper-ui-deployment.yaml
   kubectl apply -f k8s/ui/cryptkeeper-ui-service.yaml
   ```

6. **Deploy Redis (Optional)**:
   ```sh
   kubectl apply -f k8s/redis/redis-deployment.yaml
   kubectl apply -f k8s/redis/redis-service.yaml
   ```

### Accessing the Application
Check if pods are running
```sh
kubectl get pods -n cryptkeeper
```

```sh
kubectl port-forward svc/cryptkeeper-backend 8000:8000 -n cryptkeeper
kubectl port-forward svc/cryptkeeper-ui 8080:80 -n cryptkeeper
```

- **CryptKeeper UI**: [http://localhost:8080](http://localhost:8080)
- **CryptKeeper API**: [http://localhost:8000](http://localhost:8000)


### Checking for logs
```sh
kubectl logs cryptkeeper-backend-df9b96cdc-wsg6r -n cryptkeeper
kubectl logs -n kube-system -l component=kube-apiserver
kubectl rollout restart deployment -n cryptkeeper
```

### Cleaning things up
```sh
kubectl delete namespace cryptkeeper
```

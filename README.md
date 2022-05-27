# eko-wkshp
A microservice workshop presented at EKO 2022

![Overview of System](https://raw.githubusercontent.com/ytaragin/eko-wkshp/main/overview/workshop_layout.png)

![Sequence Flow](https://raw.githubusercontent.com/ytaragin/eko-wkshp/main/overview/workshop.png)


[Rest documentation for the various services in the demo](https://editor.swagger.io/?url=https://raw.githubusercontent.com/ytaragin/eko-wkshp/main/swagger.yaml)





# Useful Commands
## Kubernetes
``` shell

# set default namespace
kubectl config set-context --current --namespace=workshop

# See all running pods
kubectl get pods

# See all services
kubectl get services

# Restart the protection pods
kubectl rollout restart deployment protection

```

## Docker Commands
```
# Build the protection service
docker build -t prot-container:l1 .

# simple run the docker locally


# Run protection docker while pointing to Kubernetes hosted services
docker run -it --rm  -e TASKSHOST='localhost:30003' -e TUNNELURL='http://localhost:30002' --network="host"  --name prot prot-container:l1 



```
## Kind Commands
```
# Upload docker image into cluster
kind load docker-image prot-container:l1 --name workshop

```

## Useful busybox pod
```
# Download the image locally
docker pull  radial/busyboxplus:curl

# Push image into kubernetes cluster
kind load docker-image radial/busyboxplus:curl --name workshop

# To run the pod
kubectl run -i --rm --tty debug --image=radial/busyboxplus:curl --restart=Never -- sh


```


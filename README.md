# eko-wkshp
A microservice workshop presented at EKO 2022

![Overview of System](https://raw.githubusercontent.com/ytaragin/eko-wkshp/main/overview/workshop_layout.png)

![Sequence Flow](https://raw.githubusercontent.com/ytaragin/eko-wkshp/main/overview/workshop.png)


[Rest documentation for the various services in the demo](https://editor.swagger.io/?url=https://raw.githubusercontent.com/ytaragin/eko-wkshp/main/swagger.yaml)







# Description of Environment

| Service      | Access inside of Kuberentes | Access outside of Kubernetes |
| ----------- | ----------- | ----------- |
| Tasks Service Rest endpoint | http://tasks-svc:8080  | http://localhost:30001  |
| Tasks Service GRPC endpoint | tasks-grpc:9001  | localhost:30003  |
| Tunnel Service | http://tunnel-svc:8080   | http://localhost:3002  |
| Protection Service |  http://protection-svc:8080 |  http://localhost:30004 |

# Workshop Stage Breakdown

| Stage      | Description | Sequence Steps | Expected Duration |
| ----------- | ----------- | ----------- | ----------- |
| [ Stage 1 ](steps.md#stage-1) | Create basic Go service | Setup | 8 minutes |
| [ Stage 2 ](steps.md#stage-2) | Add a rest endpoint | Setup | 12 Minutes |
| [ Stage 3 ](steps.md#stage-3) | Wrap service in docker container and run in Kuberentes | Setup | 10 minutes |
| [ Stage 4 ](steps.md#stage-4)| Add VPG creation endpoint and create tasks | Steps 1-3 | 15 minutes |
| [ Stage 5 ](steps.md#stage-5) | Call tunnel to create VPG | Steps 4-5 | 15 minutes |
| [ Stage 6 ](steps.md#stage-6) | Update task status to In Progress | Steps 6 | 10 minutes |
| [ Stage 7 ](steps.md#stage-7)| Wait for VPG completion and update task when done | Steps 8-9 | 12 minutes |




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


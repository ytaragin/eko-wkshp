# Detailed Walkthrough - Stage 3

Create a file called Dockerfile in the protection directory that will build the container
```docker
FROM golang:1.18



WORKDIR /workdir

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify


COPY . .

RUN go build -o protection

CMD ["/workdir/protection"]


```

Now we can build the docker container and push to kubernetes
Run these commands in the protection directory
```shell
# This builds the docker container
docker build -t prot-container:l1 .

# This pushes the docker container into kubernetes where it can be found
kind load docker-image prot-container:l1 --name workshop
```

Copy the provided protection.yaml file to protection-workshop/templates
```shell
cp ../protection.yaml ../protection-workshop/templates

```

Now we can push our service into kubernetes
```shell
helm upgrade wkshp ../protection-workshop

```

Can now test against kuberentes
```shell
curl localhost:30004/ping
```

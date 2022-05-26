# Step 1 - Create basic Go service
## Step Goals

At the end of this step, we will have the basic structure for the new protection service we will be creating.

The service should just print a start up message.

## Guidance

<details>
  <summary>Tips and Hints</summary>

- Create a directory called protection
- You will need to run
```shell
go mod init protection
```
to initialize the go module

</details>

## Walkthrough
<details>
  <summary>Detailed Walkthrough</summary>

### Create the structure
Run these commands in the workshop directory
```shell
# make directory
mkdir protection

# all future commands should be run from within the protection directory
cd protection

# Initialize the go package 
go mod init protection
```
### Create go file
Make a basic go file prot-svc.go
```go 
package main

import "log"

func main() {
	log.Println("Protection Service starting up ")

}
```
###  Run the program
```shell
go run .
```

</details>


# Step 2 - Add a Rest endpoint
## Step Goals
In this step add the simple ping rest endpoint so we can start running this as a service.

At the end of this step you should have running service that answers curl requests to the ping endpoint.



## Guidance

<details>
  <summary>Tips and Hints</summary>



</details>

## Walkthrough
<details>
 <summary>Detailed Walkthrough</summary>

Run to add gonic to project
```shell
go get github.com/gin-gonic/gin

```
Add import to your file
```
import (
	"log"

	"github.com/gin-gonic/gin"
)
```

Add a handler to you go program
```go
func getping(c *gin.Context) {
	//return simple json with a 200 code
	c.JSON(200, gin.H{
		"message": "protection pong",
	})
}
```

Add to main initialization
```go
func main() {
	log.Println("Protection Service starting up ")

	r := gin.Default()

	r.GET("/ping", getping) 

	log.Println("Starting to listen...")
	r.Run() // listen and serve on 0.0.0.0:8080

}
```

Run the program and test it out
```
go run .

# in another window
curl localhost:8080/ping
```

</details>

# Step 3 - Wrap with docker and run in kubernetes

## Step Goals

## Guidance

<details>
  <summary>Tips and Hints</summary>



</details>

## Walkthrough
<details>
 <summary>Detailed Walkthrough</summary>


Create a file Dockerfile to build the container
```docker
FROM golang:1.18


WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
RUN go build -v -o /usr/local/bin ./...

CMD ["protection"]

```

Build the docker container and push to kubernetes

*Note:* Due to DNS configuration in the default kind it's important to always use a tag
```shell
docker build -t prot-container:l1 .

kind load docker-image prot-container:l1 --name workshop
```

Copy the file protection.yaml to protection-workshop/templates
```shell
cp ../protection.yaml ../protection-workshop/templates

```


Then add to kubernetes
```shell
helm upgrade wkshp ../protection-workshop

```

Can now test against kuberentes
```shell
curl localhost:30004/ping
```

</details>

# Step 4 Add new endpoint and call to Tasks Service
## Step Goals

## Guidance

<details>
  <summary>Tips and Hints</summary>



</details>

## Walkthrough
<details>
 <summary>Detailed Walkthrough</summary>

We will now start  accessing services.
To make testing easier - let's make those endpoints configurable

Add these defintions in at ths top of your prot-svc.go file (after the imports)
``` go
var TUNNEL_URL string
var TASKS_GRPC_HOST string


```

```go
func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	log.Println("Protection Service starting up ")

	TASKS_GRPC_HOST = getEnv("TASKSHOST", "tasks-grpc:9001")
	TUNNEL_URL = getEnv("TUNNELURL", "http://tunnel-svc:8080")

```

Add a new handler to the Rest configuration 

```go
r.GET("/ping", getping)
r.POST("/vpg", createVPG)

```



and then add a new function which implements the handler.
This function is the primary function we will be expanding over the workshop to handle all the VPG creation logic and orchestration.

``` go
func createVPG(c *gin.Context) {
	type VPGRequest struct {
		VPGName string
	}

	var requestBody VPGRequest

	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": "Invalid Body"})
		return
	}

    log.Printf("Received request to create VPG %s", requestBody.VPGName)
   	c.JSON(200, gin.H{})

}
```
This code gets and parses the JSON from the request
But we want to start the process of creating a VPG which starts with creating a Task

So we need to add code to make GRPC call to the Tasks service

Add the grpc to the project. 
On the command line:
```shell
go get google.golang.org/grpc
go get google.golang.org/grpc/credentials/insecure
```

We also need to add the GRPC client code to our project
So copy tasks/tasks to the protection folder
```shell
cp -r ../tasks/tasks .
```


Now we can add code to make the call to the Tasks Service
```go
func createTask() (string, error) {
	conn, err := grpc.Dial(TASKS_GRPC_HOST, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("did not connect: %v", err)
		return "", err
	}
	defer conn.Close()

	c := pb.NewTasksClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ret, err := c.CreateTask(ctx, &pb.CreateTaskRequest{})
	if err != nil {
		log.Printf("could not create tast: %v", err)
		return "", err
	}
	log.Printf("New Task Created %s with status %d", ret.GetTaskid(), ret.GetStatus())
	return ret.GetTaskid(), nil
}

```
Add to your imports the missing imports (we will name it pb to make it easier to reference)
```go
	pb "tasks/tasks"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

```



and in the function createVPG we can call that code after we parse the body of the incoming request (with BindJSON).
Lets also update the return to return the new taskid
```go
	taskid, err := createTask()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": "Unable to create task"})
		return
	}
    log.Printf("Task Created %s", taskid)
   	c.JSON(200, gin.H{ "taskid": taskid})


``` 



To run this we need to add the task code to our docker file
```dockerfile
FROM golang:1.18


#### THIS IS THE NEW LINE
COPY ./tasks /usr/local/go/src/tasks/tasks


WORKDIR /usr/src/app

# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY *.go ./
RUN go build -v -o /usr/local/bin ./...

CMD ["protection"]
```

You can build and upload your code to kubernetes
```shell
docker build -t prot-container:l1 .

kind load docker-image prot-container:l1 --name workshop

# You will need to restart the pod in kubernetes
kubectl rollout restart deployment protection

# test it
curl -X POST localhost:30004/vpg -d '{"vpgname": "VPG1"}'

# You can easily see the logs for the pod
kubectl get pods

# Take real name for pod
kubectl logs protection-<THE NAME FROM THE get pods COMMAND>
```

An alternative to uplaoding the image to K8S and  time we can also run our docker outside of k8s
```shell
docker run -it --rm  -e TASKSHOST='localhost:30003' -e TUNNELURL='http://localhost:30002' --network="host"  --name prot prot-container:l1 

# can then test using local port
curl -X POST localhost:8080/vpg -d '{"vpgname": "VPG1"}'

# You can see the Task at
curl localhost:30001/tasks


```

</details>

# Step 5  Call Tunnel to create VPG
## Step Goals

## Guidance

<details>
  <summary>Tips and Hints</summary>



</details>

## Walkthrough
<details>
 <summary>Detailed Walkthrough</summary>

We can add a function to make a POST call to the Tunnel Service

```go
func callTunnelToCreateVPG(name string) (string, error) {

	postBody := []byte(fmt.Sprintf(`{"vpgname": "%s"}`, name))
	reqBody := bytes.NewBuffer(postBody)

	log.Printf("Making request to Tunnel for %s", name)
	resp, err := http.Post(TUNNEL_URL+"/vpg", "application/json", reqBody)
	//Handle Error
	if err != nil {
		log.Print("An Error Occured %v", err)
		return "", err
	}
	defer resp.Body.Close()
	//Read the response body\
	decoder := json.NewDecoder(resp.Body)

	type CreateResponse struct {
		Vpgid     string `json:"vpgid"`
		Completed int    `json:"completed"`
	}
	var retobj CreateResponse
	err = decoder.Decode(&retobj)
	if err != nil {
		log.Printf("Error decoding object %v", retobj)
		return "", err
	}

	return retobj.Vpgid, nil

}
```

In our createVPG handler we can now call that function and expand our return response
```go

	vpgid, err := callTunnelToCreateVPG(requestBody.VPGName)
	if err != nil || vpgid == "" {
		c.AbortWithStatusJSON(400, gin.H{"msg": "Bad Request"})
		return
	}
   	c.JSON(200, gin.H{"vpgid": vpgid, "taskid": taskid})

```

Build the docker and test your code

```shell
docker build -t prot-container:l1 .

docker run -it --rm  -e TASKSHOST='localhost:30003' -e TUNNELURL='http://localhost:30002' --network="host"  --name prot prot-container:l1 

curl -X POST localhost:8080/vpg -d '{"vpgname": "VPG1"}'

# You can see the VPG at

# You can see the Task at
curl localhost:30002/vpgs

```
</details>

# Step 6 Update Task Status
## Step Goals

## Guidance

<details>
  <summary>Tips and Hints</summary>



</details>

## Walkthrough
<details>
 <summary>Detailed Walkthrough</summary>

We need a function to update the status of the task that it is in progress
```go

func UpdateTask(id string, status pb.TaskMessage_TaskStatus) error {
	conn, err := grpc.Dial(TASKS_GRPC_HOST, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		log.Printf("did not connect: %v", err)
		return err
	}
	defer conn.Close()

	c := pb.NewTasksClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	ret, err := c.UpdateTask(ctx, &pb.TaskMessage{Taskid: id, Status: status})
	if err != nil {
		log.Printf("could not update task: %v", err)
		return err
	}
	log.Printf("New Task Created %s with status %d", ret.GetTaskid(), ret.GetStatus())
	return nil

}


```

We can now update the status of the task in our createVPG function 

```go
err = UpdateTask(taskid, pb.TaskMessage_INPROGRESS)
if err != nil {
    c.AbortWithStatusJSON(400, gin.H{"msg": "Unable to create task"})
    return
}
log.Printf("Task %s status set to INPROGRESS", taskid)


```

Build and test your function.

</details>

# Step 7 - Monitor VPG and update Task when done
In this step we will put the finishing touches on our service.
When the status of the VPG is 100% complete, we can mark the Task as complete.
Ideally, there would be an event sent via Kafka when it's complete but in this workshop we will just poll the Tunnel service to check the status.

We will use the /vpg/<vpgid> API on the tunnel service to get the status of the VPG.

## Step Goals

## Guidance

<details>
  <summary>Tips and Hints</summary>



</details>

## Walkthrough
<details>
 <summary>Detailed Walkthrough</summary>



Lets create a function to get the completion percentage from the Tunnel SVC
```go
func getVPGCompletionPCT(vpgid string) (int, error) {
	url := TUNNEL_URL + "/vpg/" + vpgid
	resp, err := http.Get(url)
	if err != nil {
		log.Print("An Error Occured %v", err)
		return 0, err
	}
	defer resp.Body.Close()

	type VPGTunnelResponse struct {
		VpgID         string
		VpgName       string
		CompletionPCT int
	}
	var retObj VPGTunnelResponse

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}
	log.Printf("Parsed response %s", string(body))

	err = json.Unmarshal(body, &retObj)
	if err != nil {
		return 0, err
	}
	log.Printf("Parsed response %v", retObj)

	return retObj.CompletionPCT, nil

}


```

Let's make a function to call that function and update the Task when necessary
```go
func checkVPGAndUpdateTask(vpgid string, taskid string) bool {
	pct, err := getVPGCompletionPCT(vpgid)
	if err != nil {
		log.Printf("Error getting completion for %s: %v", vpgid, err)
		//put back for later
		return false
	}
	if pct < 100 {
		return false
	}

	err = UpdateTask(taskid, pb.TaskMessage_COMPLETE)

	return err == nil
}


```


We can then create a function that will loop and update the Task
```go
func pollAndCheckStatus(vpgid string, taskid string) {
	done := false
	for !done {
		log.Printf("Sleeping to check status of vpg %s", vpgid)
		time.Sleep(5 * time.Second)
		done = checkVPGAndUpdateTask(vpgid, taskid)
	}
	log.Printf("VPG is complete %s", vpgid)
}

```


Finally, we need to trigger that function as part of our VPG creation flow  at the end of the createVPG handler

```go
go pollAndCheckStatus(vpgid, taskid)

```

Build and test your function.
Make sure to check the Task service for the COMPLETE status
</details>






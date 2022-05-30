# Stage 1 
## Stage Goals: Create basic Go service

At the end of this Stage, we will have the basic structure for the new protection service we will be creating.

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


# Stage 2
## Stage Goals: Add a Rest endpoint

In this Stage add the simple ping rest endpoint so we can start running this as a service.

At the end of this Stage you should have running service that answers curl requests to the ping endpoint.

We will use the [gin-gonic library](https://github.com/gin-gonic/gin) for our http server as is used in DSCC services.


## Guidance

<details>
  <summary>Tips and Hints</summary>

- Add the library to you project using
```
go get github.com/gin-gonic/gin
```
- Don't forget to update the imports
- Using gin:
	- To create an the gin object
		```go
		r := gin.Default()
		```
	- Use the r object to register your event listeners
		```go
		r.GET("/ping", pingHandler)
		```
	- Handler functions are functions that take one parameter - c *gin.Context
	- To return JSON from an http request you can use the JSON function on the Context object
	```go
		c.JSON(200, gin.H{
		"key": "value",
	    })
	```
	- To start the listener
		```go
		r.Run() 
		```
	- The default port is 8080
	- Use curl to test your service

</details>

## Walkthrough
<details>
 <summary>Detailed Walkthrough</summary>

Run on the command line to add gonic to project
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

# Stage 3
## Stage Goals: Wrap with docker and run in kubernetes

In this Stage, we will take the service we created and wrap it in a Docker container and then deploy to Kubernetes using Helm.  

Due to the Kind configured Kubernetes, after building the docker container, it must be pushed into the Kubernetes using the Kind command line.

**Note:** You must use a tag with the container name such as prot-container:l1 so Kubernetes can find it locally.

You will need create a helm template file for the protection service (creating a deployment and a service that exposes port 30004 as a NodePort. Use the namespace variable defined in the values.yaml file. Call your deployment protection and your service protection-svc.
Add the file to the protection-workshop helm chart. Then use helm to upgrade our  deployment called "wkshp" with the new service.


## Guidance

<details>
  <summary>Tips and Hints</summary>

- Here is a Dockerfile that can be used to build and run the portection service:
<details>
  <summary>Dockerfile</summary>

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
</details>

- To create a helm template file, you can copy the file protection-workshop/templates/tunnel.yaml to a file called protection.yaml and modify accordingly.
Those two deployments should be structurally very similar. 
- Remember to use 30004 as the NodePort. It is the one the Kubernetes instance is configured to expose.
- Alternatively, in the root directory of the workshop is a file protection.yaml which can be used.
- The command to upgrade a deployed helm chart is
	```shell
	helm upgrade <deployment name> <Chart directory>
	```
- Use curl, to test the service

</details>

## Walkthrough
<details>
 <summary>Detailed Walkthrough</summary>


Create a file called Dockerfile in the protection directory that will build the container
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

</details>

# Stage 4
## Stage Goals: Add new endpoint and call to Tasks Service (steps 1-3) 
Now we will start getting in to the main parts of the workshop.

We will add the /vpg Post endpoint to the service. When that end point is called, we will run the flow descirbed in the sequence diagram. At this point of the workshop, we will only create the task in the Task Service. In subsequent Stages, we will fill in more logic.

Calls to the Task Service are made using GRPC. The protobuf definition file can be found in the [tasks/tasks/tasks.proto file](tasks/tasks/tasks.proto).
You will need to copy the go files into the protection directory (keep them in a tasks subdirectory) so they can be used by the protection service.

You will need to add the endpoints of the GRPC service into your code. 
You will need to add google.golang.org/grpc and google.golang.org/grpc/credentials/insecure to your project.

You will need to update the docker file to copy the tasks folder with the GRPC stubs into the docker file at the Go Root which is at /usr/local/go/src. (See the tips and hints Section for a full Docker file)
Also note that since we already have a pod instance running inside kubernetes, we will need to restart the pod after building the new image and using the kind command line to push the docker container into Kubernetes (See the guidance section for command lines how to do all that.
	
Inside Kubernetes, the task service is running at tasks-grpc:9001


## Guidance

<details>
  <summary>Tips and Hints</summary>

- In Go, to parse JSON you need a struct defined that matches the structure and names of the JSON object - except for the case of the letters.
  **NOTE:** The first character of your field name in the struct, *must* be capitilized.
- So here is how to parse the JSON of the incoming request that is handled by gin-gonic:	
	You can define and create an object like this:
	```go
	type RequestObj struct {
		Field1 string
	}

	var requestBody RequestObj

	c.BindJSON(&requestBody)
	```
	

- To call a GRPC client you must create an object like this
	```go
		conn, err := grpc.Dial(TASKS_GRPC_HOST, grpc.WithTransportCredentials(insecure.NewCredentials()))

		if err != nil {
			log.Printf("did not connect: %v", err)
			return "", err
		}
		defer conn.Close()

		c := pb.NewTasksClient(conn)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

	```
	You can then use the object c as a client to the tasks service
- Here is an updated Dockerfile that can be used to build and run the portection service:

<details>
  <summary>Dockerfile</summary>

```docker
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
</details>

	
- Here are commands on how to build and deploy the docker file into kubernetes

<details>
  <summary>How to build and run docker in Kubernetes</summary>

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
</details>

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

In your main function, add a new handler to the Rest configuration 

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

	// Add the code below steps here

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
Add those two packages to the imports in your file:
```go

	pb "tasks/tasks"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

```



and in the function createVPG, after the  we can call that code after we parse the body of the incoming request (with BindJSON).
Lets also update the return to return the new taskid
```go
	taskid, err := createTask()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": "Unable to create task"})
		return
	}
    log.Printf("Task Created %s", taskid)
   	c.JSON(200, gin.H{ "taskid": taskid})  // only return task id so far


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
</details>

	
## Stage 4 Addendum
	
An alternative to uploading the image to K8S and  time we can also run our docker outside of k8s
```shell
docker run -it --rm  -e TASKSHOST='localhost:30003' -e TUNNELURL='http://localhost:30002' --network="host"  --name prot prot-container:l1 

# can then test using local port
curl -X POST localhost:8080/vpg -d '{"vpgname": "VPG1"}'

# You can see the Task at
curl localhost:30001/tasks


```


# Stage 5  
## Stage Goals: Call Tunnel to create VPG  (Steps 4-5)

We will now extend our CreateVPG handler function to create a VPG after it created a task.

This will be done by calling the Tunnel service using the POST /vpg endpoint. The [swagger documentation](https://editor.swagger.io/?url=https://raw.githubusercontent.com/ytaragin/eko-wkshp/main/swagger.yaml) describes the endpoint.

Once you do that update the JSON object returned by the vpg function to include the vpgid as well.

Inside Kubernetes, the tunnel svc endpoint is at http://tunnel-svc:8080

## Guidance

<details>
  <summary>Tips and Hints</summary>

- Use the http.Post function to make a POST call to the service
- Here is how you can turn a string into a buffer to be used by the post function
```go
	reqBody := bytes.NewBuffer([]byte(mystring))
```

- The reponse from the tunnel service will be a JSON block which we will need to parse. Earlier, we used a function provided by gin-gonic to parse the JSON. 
  This is not a gin function so we will need to parse the JSON using standard go functionality - though the concept is similar.
- Create a struct that matches the fields of the JSON (remember first letter of each field must be capitilized). Then you can use the  decoder object to decode into the object.
```go
	// define and create an object
	type CreateResponse struct {
		Vpgid     string `json:"vpgid"`
		Completed int    `json:"completed"`
	}
	var retobj CreateResponse

   decoder := json.NewDecoder(resp.Body) // resp is from the POST call
   	err = decoder.Decode(&retobj) // decode into the defined object

```



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

# Stage 6 
## Stage Goals: Update Task Status (Step 6) 

Now that we have initiated the creation of the VPG we can update the Tasks service that the Task is in progress using the UpdateTask endpoint.

## Guidance

<details>
  <summary>Tips and Hints</summary>

- This is very similar to how we called CreateTask before, except you will need to create an instance of pb.TaskMessage to pass to the UpdateTask function.
- Here is how you create a struct in go
```go
	// if you have type like
	type MyStruct struct {
		Field1 int,
		Field2 string,
	}

	// you can create an instace
	myinst := MyStruct{Field1: 10, Field2: "Hello"}
```
- Set the status to pb.TaskMessage_INPROGRESS


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

# Stage 7
## Stage Goals: Monitor VPG and update Task when done (Steps 8-10)

In this Stage we will put the finishing touches on our service.

Every few seconds, you can check the /vpg/<vpgid> API on the tunnel service to get the completion status of the VPG.
When the status of the VPG is 100% complete, mark the Task as complete.

(Ideally, there would be an event sent via Kafka when it's complete but in this workshop we will just poll the Tunnel service to check the status. )

You can call the /tasks endpoint on the task service (using curl) to see when the status is complete.

## Guidance

<details>
  <summary>Tips and Hints</summary>

- Call the UpdateTask with status of pb.TaskMessage_COMPLETE
- In go you can sleep via
	```go
	time.Sleep(5 * time.Second)
	```
- To run a function in the background, use the go keyword
	```go
	// if you have a function
	func pollAndCheckStatus(vpgid string, taskid string) {}

	// you can call in a background thread
	go pollAndCheckStatus(vpgid, taskid)

	```

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
go pollAndCheckStatus(vpgid, taskid) // Calling with go runs the function in a thread

```

Build and test your function.
Make sure to check the Task service for the COMPLETE status
</details>






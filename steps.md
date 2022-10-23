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
Only click if you are sure you want to see more information

[ Stage 2 Detailed Walkthrough ](steps/stage2walkthrough.md)

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
FROM golang:1.18 AS builder
WORKDIR /workdir
# pre-copy/cache go.mod for pre-downloading dependencies and only redownloading them in subsequent builds if they change
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN CGO_ENABLED=0 go build -o protection

FROM scratch
COPY --from=builder /workdir/protection /
CMD ["/protection"]

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
Only click if you are sure you want to see more information

[ Stage 3 Detailed Walkthrough ](steps/stage3walkthrough.md)



</details>

# Stage 4
## Stage Goals: Add new endpoint and call to Tasks Service (steps 1-3) 
Now we will start getting in to the main parts of the workshop.

We will add the /vpg Post endpoint to the service. When that end point is called, we will run the flow descirbed in the sequence diagram. At this point of the workshop, we will only create the task in the Task Service. In subsequent Stages, we will fill in more logic.


## Option A - GRPC

Calls to the Task Service are made using GRPC. The protobuf definition file can be found in the [tasks/tasks/tasks.proto file](tasks/tasks/tasks.proto).
You will need to copy the go files into the protection directory (keep them in a tasks subdirectory) so they can be used by the protection service.

You will need to add the endpoints of the GRPC service into your code. 
You will need to add google.golang.org/grpc and google.golang.org/grpc/credentials/insecure to your project.

You will need to update the docker file to copy the tasks folder with the GRPC stubs into the docker file at the Go Root which is at /usr/local/go/src. (See the tips and hints Section for a full Docker file)
Also note that since we already have a pod instance running inside kubernetes, we will need to restart the pod after building the new image and using the kind command line to push the docker container into Kubernetes (See the guidance section for command lines how to do all that.
	
Inside Kubernetes, the task service grpc endpoint is running at tasks-grpc:9001

	



### Guidance

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

	
- Here are commands on how to build and deploy the docker file into kubernetes
<details>
  <summary>How to build and run docker in Kubernetes</summary>

```shell
../utils/deploy_protection.sh

# test it
curl -X POST localhost:30004/vpg -d '{"vpgname": "VPG1"}'

```
</details>

</details>

	
## Option B - Rest

Rather then using the Task service's GRPC endpoint, it can also be called using the Rest endpoint as defined in 
the [swagger documentation](https://editor.swagger.io/?url=https://raw.githubusercontent.com/ytaragin/eko-wkshp/main/swagger.yaml)


Calls to the Task Service are made using GRPC. The protobuf definition file can be found in the [tasks/tasks/tasks.proto file](tasks/tasks/tasks.proto).
You will need to copy the go files into the protection directory (keep them in a tasks subdirectory) so they can be used by the protection service.

You will need to add the endpoints of the GRPC service into your code. 
You will need to add google.golang.org/grpc and google.golang.org/grpc/credentials/insecure to your project.

You will need to update the docker file to copy the tasks folder with the GRPC stubs into the docker file at the Go Root which is at /usr/local/go/src. (See the tips and hints Section for a full Docker file)
Also note that since we already have a pod instance running inside kubernetes, we will need to restart the pod after building the new image and using the kind command line to push the docker container into Kubernetes (See the guidance section for command lines how to do all that.
	
Inside Kubernetes, the task service rest endpoint is running at tasks-svc:8080

	



### Guidance

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
		Taskid     string `json:"vpgid"`
		Status int    `json:"status"`
	}
	var retobj CreateResponse

   decoder := json.NewDecoder(resp.Body) // resp is from the POST call
   	err = decoder.Decode(&retobj) // decode into the defined object

```

	
- Here are commands on how to build and deploy the docker file into kubernetes
<details>
  <summary>How to build and run docker in Kubernetes</summary>

```shell
# run from within the protection directory
../utils/deploy_protection.sh

# test it
curl -X POST localhost:30004/vpg -d '{"vpgname": "VPG1"}'

```
</details>

</details>

## Walkthrough - both options
<details>
 <summary>Detailed Walkthrough</summary>

Only click if you are sure you want to see more information

[ Stage 4 Detailed Walkthrough ](steps/stage4walkthrough.md)

</details>
	

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

Only click if you are sure you want to see more information

[ Stage 5 Detailed Walkthrough ](steps/stage5walkthrough.md)

</details>

# Stage 6 
## Stage Goals: Update Task Status (Step 6) 

Now that we have initiated the creation of the VPG we can update the Tasks service that the Task is in progress using the UpdateTask endpoint.
As in Step 4 we can use GRPC or Rest for this step.


## Option A - GRPC
### Guidance

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

## Option B - Rest
### Guidance
- The Update Task command, as described in the [swagger documentation](https://editor.swagger.io/?url=https://raw.githubusercontent.com/ytaragin/eko-wkshp/main/swagger.yaml) shows a PUT command is needed for this function
- The go http package does not have a default method for PUT the way it does for POST or GET
- This is the way to make a PUT call 
```go
	req, err := http.NewRequest("PUT", u, reqBody)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)

```

## Walkthrough
<details>
 <summary>Detailed Walkthrough</summary>

Only click if you are sure you want to see more information

[ Stage 6 Detailed Walkthrough ](steps/stage6walkthrough.md)

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
Only click if you are sure you want to see more information

[ Stage 7 Detailed Walkthrough ](steps/stage7walkthrough.md)


</details>





# Stage 8
## Stage Goals: Store VPG and task info in the DB

In this Stage we will store information in the database

- We need to create the tables in the database
   ```shell
   kubectl apply -f migrate.yaml

   ```
- This is the table that is created
```sql
CREATE TABLE vpg (
   vpg_id VARCHAR(50) PRIMARY KEY,
   task_id VARCHAR(50) NOT NULL,
   status INT NOT NULL
);
```

- There is a tool in the utils directory postgres.sh that connects you a psql client to see what is in the database
- Here are the connection parameters for the database
```go
	const (
		pgUser     = "postgres"
		pgPassword = "mysecret"
		pgHost     = "wkshp-postgresql"
		pgPort     = 5432
		pgDatabase = "protection"
	)

```
- This is the library to add to go.mod to support postgres db
```go
	github.com/jackc/pgx/v4 v4.17.2
```
- In DSCC, we use [ SQLC ](https://sqlc.dev/) for SQL queries.  
    - The folder protdb/db contains a module that contains a set of functions generated for queries.
	- Copy the db folder into your protection folder
- These are the functions provided by the provided module
	- AddVpg
	- GetNonReadyVPGs
	- UpdateStatus
- These are the status values to use toward the database
    - In Progress = 1
	- Complete = 2
	- Error = 3



## Guidance

<details>
  <summary>Tips and Hints</summary>

- This is the way to get a connection to a postgres db
```go
connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", pgUser, pgPassword, hostname, pgPort, pgDatabase)
db, err := sql.Open("pgx", connectionString)
```
- This is the way to create a repository object to query with:
```go
var repo repository.Querier = repository.New(db)
```
- The function GetNonReadyVPGs provides a list of VPGs not marked as complete.
- Here is how you can loop over the returned VPGS
```go
for _, v := range nonReadyVPGs {
		// do something with the vpg records
	}
```
- Context Object can be created like this
```go
var ctx context.Context = context.Background()
```



</details>

## Walkthrough
<details>
 <summary>Detailed Walkthrough</summary>
Only click if you are sure you want to see more information

[ Stage 8 Detailed Walkthrough ](steps/stage8walkthrough.md)


</details>

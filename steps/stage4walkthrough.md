# Detailed Walkthrough - Stage 4




## Walkthrough
<details>
 <summary>Detailed Walkthrough</summary>

We will now start  accessing services.

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


## 4A - GRPC Details


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

	pb "protection/tasks"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

```

## 4B - Rest Approach
This is a version of the CreateTask function that uses Rest instead

```go
func createTaskRest() (string, error) {

	log.Printf("Making REST call to create a task")

	resp, err := http.Post(TASKS_URL+"/task", "application/json", nil)
	//Handle Error
	if err != nil {
		log.Print("An Error Occured %v", err)
		return "", err
	}
	defer resp.Body.Close()
	//Read the response body\
	decoder := json.NewDecoder(resp.Body)

	type TaskResponse struct {
		Taskid string `json:"taskid"`
		Status int    `json:"status"`
	}
	var retobj TaskResponse
	err = decoder.Decode(&retobj)
	if err != nil {
		log.Printf("Error decoding object %v", retobj)
		return "", err
	}

	log.Printf("New Task Created %s with status %d", retobj.Taskid, retobj.Status)
	return retobj.Taskid, nil
}

```



## Both Approaches


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


You can build and upload your code to kubernetes
```shell
# run from within the protection directory
../utils/deploy_protection.sh

# test it
curl -X POST localhost:30004/vpg -d '{"vpgname": "VPG1"}'


# You can easily see the logs for the pod
kubectl get pods

# Take real name for pod
kubectl logs protection-<THE NAME FROM THE get pods COMMAND>
```

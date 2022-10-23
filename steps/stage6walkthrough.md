# Detailed Walkthrough - Stage 6

## Option A - GRPC
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

## Option B - Rest
A UpdateTask version that uses Rest
```go
func UpdateTask(id string, status int) error {
	log.Printf("Making REST call to update a task %s", id)

	postBody := []byte(fmt.Sprintf(`{"status": %d}`, status))
	reqBody := bytes.NewBuffer(postBody)

	u := fmt.Sprintf("%s/task/%s", TASKS_URL, id)

	req, err := http.NewRequest("PUT", u, reqBody)
	if err != nil {
		return err
	}

	_, err = http.DefaultClient.Do(req)

	//Handle Error
	if err != nil {
		log.Print("An Error Occured %v", err)
		return err
	}

	log.Printf("Status Updated")

	return nil
}

```


## Both Options


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

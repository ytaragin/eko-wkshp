package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	pb "protection/tasks"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var TUNNEL_URL string
var TASKS_URL string
var TASKS_GRPC_HOST string

func getping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "protection pong",
	})
}

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

func createTaskGPRC() (string, error) {
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

	taskid, err := createTask()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": "Unable to create task"})
		return
	}
	log.Printf("Task Created %s", taskid)

	vpgid, err := callTunnelToCreateVPG(requestBody.VPGName)
	if err != nil || vpgid == "" {
		c.AbortWithStatusJSON(400, gin.H{"msg": "Bad Request"})
		return
	}
	c.JSON(200, gin.H{"vpgid": vpgid, "taskid": taskid})

}

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
	TASKS_URL = getEnv("TASKSURL", "http://tasks-svc:8080")
	TUNNEL_URL = getEnv("TUNNELURL", "http://tunnel-svc:8080")

	r := gin.Default()

	r.GET("/ping", getping)
	r.POST("/vpg", createVPG)

	log.Println("Starting to listen...")
	r.Run() // listen and serve on 0.0.0.0:8080

}

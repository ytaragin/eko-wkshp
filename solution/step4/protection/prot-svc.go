package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "protection/tasks"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var TUNNEL_URL string
var TASKS_GRPC_HOST string

func getping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "protection pong",
	})
}

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
	c.JSON(200, gin.H{"taskid": taskid})

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
	TUNNEL_URL = getEnv("TUNNELURL", "http://tunnel-svc:8080")

	r := gin.Default()

	r.GET("/ping", getping)
	r.POST("/vpg", createVPG)

	log.Println("Starting to listen...")
	r.Run() // listen and serve on 0.0.0.0:8080

}

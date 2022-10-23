package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"

	pb "tasks/tasks"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type TaskManager struct {
	tasks map[string]pb.TaskMessage_TaskStatus
}

var taskMgr TaskManager

func (t TaskManager) CreateTask() (string, pb.TaskMessage_TaskStatus) {
	taskid := uuid.New().String()
	log.Println("Creating new task: " + taskid)
	t.tasks[taskid] = pb.TaskMessage_CREATED

	return taskid, t.tasks[taskid]
}

func (t TaskManager) UpdateTask(taskid string, status pb.TaskMessage_TaskStatus) error {
	curstatus, prs := t.tasks[taskid]
	if !prs {
		log.Println("Unknown task ", taskid)
		return fmt.Errorf("Unknown ID: %s", taskid)
	}

	log.Println("Updating status for task: "+taskid+" from ", curstatus.String(), " to ", status.String())
	t.tasks[taskid] = pb.TaskMessage_TaskStatus(status)

	return nil

}

func (t TaskManager) GetTaskStatus(taskid string) pb.TaskMessage_TaskStatus {
	curstatus, prs := t.tasks[taskid]
	if prs {
		return curstatus
	} else {
		log.Println("Unknown task ", taskid)
		return -1
	}

}

func (t TaskManager) GetAllTasks() string {
	str := "Task: Status\n"

	for taskid, status := range t.tasks {
		str += fmt.Sprintf("%s: %s\n", taskid, status.String())
	}

	return str
}

func getping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "tasks pong",
	})
}

func createTask(c *gin.Context) {
	id, status := taskMgr.CreateTask()

	c.JSON(200, gin.H{
		"taskid": id,
		"status": status,
	})
}

func putTask(c *gin.Context) {
	taskid := c.Param("taskid")
	fmt.Println("Updating status for task: " + taskid)

	type TaskRequest struct {
		Status int
	}

	var requestBody TaskRequest

	if err := c.BindJSON(&requestBody); err != nil {
		log.Println("Error parsing body")
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"msg": "Invalid Body"})
		return
	}

	log.Printf("Received request to update status to %d", requestBody.Status)

	err := taskMgr.UpdateTask(taskid, pb.TaskMessage_TaskStatus(requestBody.Status))
	if err != nil {
		c.JSON(404, gin.H{
			"taskid": taskid,
		})
		return
	}

	c.Writer.WriteHeader(200)
	return

}

func putTask2(c *gin.Context) {
	taskid := c.Param("taskid")
	fmt.Println("Updating status for task: " + taskid)

	type StatusRec struct {
		status int
	}
	statusStr := c.PostForm("status")
	statusInt, err := strconv.Atoi(statusStr)
	if err != nil {
		c.JSON(400, gin.H{
			"taskid": taskid,
		})
	}
	err = taskMgr.UpdateTask(taskid, pb.TaskMessage_TaskStatus(statusInt))
	if err != nil {
		c.JSON(404, gin.H{
			"taskid": taskid,
		})
		return
	}

	c.Writer.WriteHeader(200)
	return

}

func getTask(c *gin.Context) {
	taskid := c.Param("taskid")
	log.Println("Fetching status for task: " + taskid)

	status := taskMgr.GetTaskStatus(taskid)

	if status != -1 {
		c.JSON(200, gin.H{
			"taskid": taskid,
			"status": status,
		})
	} else {
		c.JSON(404, gin.H{
			"taskid": taskid,
			"status": -1,
		})
	}

}

func getAllTasks(c *gin.Context) {
	c.String(200, taskMgr.GetAllTasks())
}

type server struct {
	pb.UnimplementedTasksServer
}

func (s *server) CreateTask(ctx context.Context, in *pb.CreateTaskRequest) (*pb.TaskMessage, error) {
	log.Println("Received: Create Task")
	taskid, status := taskMgr.CreateTask()
	return &pb.TaskMessage{Taskid: taskid, Status: status}, nil
}

func (s *server) UpdateTask(ctx context.Context, in *pb.TaskMessage) (*pb.TaskMessage, error) {
	log.Println("Received: Update Task")
	err := taskMgr.UpdateTask(in.GetTaskid(), in.GetStatus())
	if err != nil {
		return nil, err
	}
	return in, nil
}

func startGRPC() {
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterTasksServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	taskMgr = TaskManager{
		tasks: make(map[string]pb.TaskMessage_TaskStatus),
	}

	go startGRPC()

	r := gin.Default()

	r.GET("/ping", getping)
	r.POST("/task", createTask)
	r.GET("/tasks", getAllTasks)
	r.PUT("/task/:taskid", putTask)
	r.GET("/task/:taskid", getTask)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

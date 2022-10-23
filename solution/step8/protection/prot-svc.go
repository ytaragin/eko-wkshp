package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"protection/db/repository"

	_ "github.com/jackc/pgx/v4/stdlib"

	pb "protection/tasks"

	"github.com/gin-gonic/gin"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	vpgInProgress = 1
	vpgReady      = 2
	vpgFailed     = 3
)

var TUNNEL_URL string
var TASKS_URL string
var TASKS_GRPC_HOST string

var DBCONN *sql.DB
var ctx context.Context = context.Background()
var DBCHECKTIMESECONDS int = 60

func getping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "protection pong",
	})
}

func connectToDB() *sql.DB {
	// ctx = context.Background()

	const (
		pgUser     = "postgres"
		pgPassword = "mysecret"
		pgHost     = "wkshp-postgresql"
		pgPort     = 5432
		pgDatabase = "protection"
	)

	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", pgUser, pgPassword, pgHost, pgPort, pgDatabase)
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		log.Printf("Failed to open DB", err)
		return nil
	}

	return db

}

func storeVPGinDB(vpgid string, taskid string, status int) error {
	var repo repository.Querier = repository.New(DBCONN)

	_, err := repo.AddVpg(ctx, repository.AddVpgParams{
		VpgID:  vpgid,
		TaskID: taskid,
		Status: int32(status),
	})
	if err != nil {
		fmt.Println("Unable to store to db")
		fmt.Println(err)
	}
	return err

}

func updateVPGinDB(vpgid string, status int) error {
	var repo repository.Querier = repository.New(DBCONN)

	err := repo.UpdateStatus(ctx, repository.UpdateStatusParams{
		VpgID:  vpgid,
		Status: vpgReady,
	})
	if err != nil {
		fmt.Println("Unable to store to db")
		fmt.Println(err)
	}
	return err

}

func checkOpenVPGs() bool {
	var repo repository.Querier = repository.New(DBCONN)

	nonReadyVPGs, err := repo.GetNonReadyVPGs(ctx)
	if err != nil {
		fmt.Println("Unable to read from db vpgs")
		fmt.Println(err)
		return false
	}

	cont := true
	for _, v := range nonReadyVPGs {
		cont = cont && checkVPGAndUpdateTask(v.VpgID, v.TaskID)
	}

	return cont
}

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

	err = UpdateTaskGRPC(taskid, pb.TaskMessage_COMPLETE)
	if err != nil {
		log.Printf("Could not call Tasks Service: %v", err)
		return false
	}

	updateVPGinDB(vpgid, vpgReady)

	return err == nil
}

func pollAndCheckStatus(vpgid string, taskid string) {
	done := false
	for !done {
		log.Printf("Sleeping to check status of vpg %s", vpgid)
		time.Sleep(5 * time.Second)
		done = checkVPGAndUpdateTask(vpgid, taskid)
	}
	log.Printf("VPG is complete %s", vpgid)
}

func UpdateTaskGRPC(id string, status pb.TaskMessage_TaskStatus) error {
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

func UpdateTaskRest(id string, status int) error {
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

func createTaskGRPC() (string, error) {
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

	taskid, err := createTaskGRPC()
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

	err = storeVPGinDB(vpgid, taskid, vpgInProgress)
	if err != nil {
		log.Printf("Unable to update VPG in DB %v", err)
	}

	// err = UpdateTaskGrpc(taskid, pb.TaskMessage_INPROGRESS)
	err = UpdateTaskGRPC(taskid, 1)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"msg": "Unable to create task"})
		return
	}
	log.Printf("Task %s status set to INPROGRESS", taskid)

	go pollAndCheckStatus(vpgid, taskid)

	c.JSON(200, gin.H{"vpgid": vpgid, "taskid": taskid})

}

func pollDBForVPGs() {
	done := false
	for !done {
		log.Printf("Sleeping to check status of vpgs ")
		time.Sleep(time.Duration(DBCHECKTIMESECONDS) * time.Second)
		done = !checkOpenVPGs()
	}
	log.Printf("Stopping to check for VPGs in DB")
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

func main() {
	log.Println("Protection Service starting up...")

	TASKS_GRPC_HOST = getEnv("TASKSHOST", "tasks-grpc:9001")
	TASKS_URL = getEnv("TASKSURL", "http://tasks-svc:8080")
	TUNNEL_URL = getEnv("TUNNELURL", "http://tunnel-svc:8080")

	DBCONN = connectToDB()

	go pollDBForVPGs()

	r := gin.Default()

	r.GET("/ping", getping)
	r.POST("/vpg", createVPG)

	log.Println("Starting to listen...")
	r.Run() // listen and serve on 0.0.0.0:8080

}

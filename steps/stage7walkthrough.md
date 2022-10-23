# Detailed Walkthrough - Stage 7

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

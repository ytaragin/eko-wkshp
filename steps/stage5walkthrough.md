# Detailed Walkthrough - Stage 5

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



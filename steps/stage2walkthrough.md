# Detailed Walkthrough - Stage 2

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

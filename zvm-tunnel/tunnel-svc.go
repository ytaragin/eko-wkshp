package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VPGRecord struct {
	Vpgid         string
	TimeRemaining int
	VpgName       string
}

var vpgs = make(map[string]*VPGRecord)

func getping(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "zvm pong",
	})
}

func createVPGGen(min int, max int) func(*gin.Context) {

	return func(c *gin.Context) {
		type VPGRequest struct {
			VPGName string
		}

		var requestBody VPGRequest

		if err := c.BindJSON(&requestBody); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"msg": "Invalid Body"})
			return
		}

		vpgid := uuid.New().String()
		vpgs[vpgid] = &VPGRecord{
			Vpgid:         vpgid,
			TimeRemaining: rand.Intn(max-min+1) + min,
			VpgName:       requestBody.VPGName,
		}

		fmt.Println("Creating new vpg: ", vpgid, " in ", vpgs[vpgid])

		c.JSON(200, gin.H{
			"vpgid":     vpgid,
			"completed": vpgs[vpgid].TimeRemaining,
		})
	}
}

func timeUpdater() {
	ticker := time.NewTicker(1 * time.Second)
	for _ = range ticker.C {
		// fmt.Println("Tock")
		for _, rec := range vpgs {
			if rec.TimeRemaining > 0 {
				rec.TimeRemaining = rec.TimeRemaining - 1
			}

		}
	}
}

func getAllVpgsJSON(c *gin.Context) {
	type msg struct {
		id     string
		status int
	}

	var v []*VPGRecord

	for _, rec := range vpgs {
		v = append(v, rec)
	}

	c.JSON(200, gin.H{
		"vpgs": v,
	})

}
func getAllVpgsStr(c *gin.Context) {
	str := "VPG: TimeLeft\n"

	for _, rec := range vpgs {
		str += fmt.Sprintf("%s(%s): %d\n", rec.Vpgid, rec.VpgName, rec.TimeRemaining)
	}

	c.String(200, str)

}

// func createVPG(c *gin.Context) {
// 	vpgid := uuid.New().String()
// 	vpgs[vpgid] = rand.Intn(MAXTIME - MINTIME + 1) + MINTIME
// 	fmt.Println("Creating new vpg: ",vpgid," in ", vpgs[vpgid])

// 	c.JSON(200, gin.H{
// 		"vpgid": vpgid,
// 		"status": 1,
// 	})
// }

func getVPG(c *gin.Context) {
	vpgid := c.Param("vpgid")
	fmt.Println("Fetching status for task: " + vpgid)

	rec, prs := vpgs[vpgid]

	fmt.Println("prs: ", prs)

	if prs {
		c.JSON(200, gin.H{
			"vpgid":         vpgid,
			"vpgname":       rec.VpgName,
			"completionpct": 100 - rec.TimeRemaining,
		})
	} else {
		c.JSON(404, gin.H{
			"vpgid": vpgid,
		})
	}

}

func getenv(envname string, def int) int {
	str := os.Getenv(envname)
	intVar, err := strconv.Atoi(str)
	if err != nil {
		intVar = def
	}
	return intVar
}

func main() {
	// rand.Seed(time.Now().UnixNano())
	min := getenv("MINTIME", 10)
	max := getenv("MAXTIME", 30)

	go timeUpdater()

	r := gin.Default()

	r.GET("/ping", getping)
	r.POST("/vpg", createVPGGen(min, max))
	r.GET("/vpg/:vpgid", getVPG)
	r.GET("/vpgs", getAllVpgsStr)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

//gotest
//testify

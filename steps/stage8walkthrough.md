# Detailed Walkthrough - Stage 8

- Copy the db folder to the protection folder
```shell
cp -r ../protdb/db .
```

- The go.mod should have this section
```go
    require (
        github.com/gin-gonic/gin v1.8.1
        github.com/jackc/pgx/v4 v4.17.2
        google.golang.org/grpc v1.50.1
        google.golang.org/protobuf v1.28.1
    )
```

- Create a function to connect to the database:
```go
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
```
THis function only needs to be called one time and store the result in a static variable
```go
var DBCONN *sql.DB = connectToDB()
var ctx context.Context = context.Background()
```

- Here is a function to create a record in the database
``` go
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

```

- A function to update a VPG record
```go
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

```




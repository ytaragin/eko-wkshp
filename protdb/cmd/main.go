package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/pkg/errors"
	"www.example.com/workshop/db/repository"
)

const (
	pgUser     = "postgres"
	pgPassword = "mysecret"
	pgHost     = "ccs-pg"
	pgPort     = 5432
	pgDatabase = "protection"
)

const (
	vpgInProgress = 1
	vpgReady      = 2
	vpgFailed     = 3
)

func main() {

	ctx := context.Background()

	err := doRepositoryStuff(ctx, "ccs-pg")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func doRepositoryStuff(ctx context.Context, hostname string) error {
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", pgUser, pgPassword, hostname, pgPort, pgDatabase)
	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		return errors.Wrap(err, "Failed to open DB")
	}

	var repo repository.Querier = repository.New(db)

	vpg1, err := repo.AddVpg(ctx, repository.AddVpgParams{
		VpgID:  uuid.NewString(),
		TaskID: uuid.NewString(),
		Status: vpgInProgress,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to add VPG")
	}

	err = repo.UpdateStatus(ctx, repository.UpdateStatusParams{
		VpgID:  vpg1.VpgID,
		Status: vpgReady,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to update VPG status")
	}

	_, err = repo.AddVpg(ctx, repository.AddVpgParams{
		VpgID:  uuid.NewString(),
		TaskID: uuid.NewString(),
		Status: vpgInProgress,
	})
	if err != nil {
		return errors.Wrap(err, "Failed to add VPG")
	}

	nonReadyVPGs, err := repo.GetNonReadyVPGs(ctx)
	if err != nil {
		return errors.Wrap(err, "Failed to get VPGs")
	}

	fmt.Printf("Found %d non-ready VPGs in the DB\n", len(nonReadyVPGs))

	return nil
}

package cmd

import (
	"github.com/cty002718/wow-arena-leaderboard/pkg/orm"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func migrateCmd() *cobra.Command {
	migrateCmd := &cobra.Command{
		Use:   "migrate",
		Short: "Migrate database",
	}

	migrateCmd.AddCommand(
		migrateUpCmd(),
		migrateDownCmd(),
	)

	return migrateCmd
}

func migrateUpCmd() *cobra.Command {
	migrateUpCmd := &cobra.Command{
		Use:   "up",
		Short: "Migrate up",
		Run: func(cmd *cobra.Command, args []string) {
			migrateUp()
		},
	}

	return migrateUpCmd
}

func migrateDownCmd() *cobra.Command {
	migrateDownCmd := &cobra.Command{
		Use:   "down",
		Short: "Migrate down",
		Run: func(cmd *cobra.Command, args []string) {
			migrateDown()
		},
	}

	return migrateDownCmd
}

func createMigrator() (*migrate.Migrate, error) {
	db, err := orm.OpenDB()
	if err != nil {
		return nil, errors.Errorf("Failed to open db: %v", err)
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, errors.Errorf("Failed to get sql db: %v", err)
	}

	dbDriver, err := postgres.WithInstance(sqlDb, &postgres.Config{})
	if err != nil {
		return nil, errors.Errorf("Failed to get db driver: %v", err)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file://migrations", "postgres", dbDriver)
	if err != nil {
		return nil, errors.Errorf("Failed to get migrator: %v", err)
	}

	return migrator, nil
}

func migrateUp() {
	migrator, err := createMigrator()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create migrator")
	}

	err = migrator.Up()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to migrate up")
	}
}

func migrateDown() {
	migrator, err := createMigrator()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to create migrator")
	}

	err = migrator.Down()
	if err != nil {
		logrus.WithError(err).Fatal("Failed to migrate down")
	}
}

package db

import (
	"database/sql"
	"fmt"
	"os"
	"strings"

	// postgres driver
	_ "github.com/lib/pq"
	migrate "github.com/rubenv/sql-migrate"
	"github.com/spf13/cobra"
)

const (
	dialect = "postgres"
)

var (
	dryRun           bool
	migrateTableName string
	migrateDir       string
	max              int
)

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "schema apply command",
	RunE: func(cmd *cobra.Command, args []string) error {
		smode := "disable"
		if Opt.SSLMode {
			smode = "verify-ca"
		}
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			Opt.Host, Opt.User, Opt.Password, Opt.DB, Opt.Port, smode, Opt.TimeZone,
		)
		db, err := sql.Open(dialect, dsn)
		if err != nil {
			return err
		}
		parts := strings.Split(migrateDir, "/")
		schemaName := parts[len(parts)-1]
		migrate.SetSchema(schemaName)
		migrate.SetTable(migrateTableName)
		src := migrate.FileMigrationSource{
			Dir: migrateDir,
		}

		if dryRun {
			migrations, _, err := migrate.PlanMigration(db, dialect, src, migrate.Up, max)
			if err != nil {
				return fmt.Errorf("Cannot plan migration: %s", err)
			}
			for _, m := range migrations {
				PrintMigration(m, migrate.Up)
			}
			return nil
		}

		n, err := migrate.ExecMax(db, dialect, src, migrate.Up, max)
		if err != nil {
			return err
		}
		fmt.Printf("Applied %d migrations", n)
		return nil
	},
}

var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Insert seed data into the database",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 正しい接続文字列を作成
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable TimeZone=%s",
			Opt.Host, Opt.User, Opt.Password, Opt.DB, Opt.Port, Opt.TimeZone,
		)
		db, err := sql.Open(dialect, dsn)
		if err != nil {
			return fmt.Errorf("failed to connect to database: %w", err)
		}
		defer db.Close()

		// seed_data.sql のパスを指定
		seedFile := "./database/migration/seed_data.sql"
		content, err := os.ReadFile(seedFile)
		if err != nil {
			return fmt.Errorf("failed to read seed file (%s): %w", seedFile, err)
		}

		// SQL を実行
		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute seed queries: %w", err)
		}

		fmt.Println("Seed data applied successfully!")
		return nil
	},
}

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "schema down command",
	RunE: func(cmd *cobra.Command, args []string) error {
		smode := "disable"
		if Opt.SSLMode {
			smode = "verify-ca"
		}
		dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
			Opt.Host, Opt.User, Opt.Password, Opt.DB, Opt.Port, smode, Opt.TimeZone,
		)
		db, err := sql.Open(dialect, dsn)
		if err != nil {
			return err
		}
		parts := strings.Split(migrateDir, "/")
		schemaName := parts[len(parts)-1]
		migrate.SetSchema(schemaName)
		migrate.SetTable(migrateTableName)
		src := migrate.FileMigrationSource{
			Dir: migrateDir,
		}

		if dryRun {
			migrations, _, err := migrate.PlanMigration(db, dialect, src, migrate.Down, max)
			if err != nil {
				return fmt.Errorf("Cannot plan migration: %s", err)
			}
			for _, m := range migrations {
				PrintMigration(m, migrate.Down)
			}
			return nil
		}

		n, err := migrate.ExecMax(db, dialect, src, migrate.Down, max)
		if err != nil {
			return err
		}
		fmt.Printf("Applied %d migrations", n)
		return nil
	},
}

func init() {
	applyCmd.PersistentFlags().StringVarP(&migrateTableName, "table", "t", "migrations", "migration table name")
	applyCmd.PersistentFlags().StringVarP(&migrateDir, "dir", "r", "", "migration directory")
	applyCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "dry run mode")
	applyCmd.PersistentFlags().IntVarP(&max, "max", "m", 0, "limit of apply migration")
	SchemaCmd.AddCommand(applyCmd)

	downCmd.PersistentFlags().StringVarP(&migrateTableName, "table", "t", "migrations", "migration table name")
	downCmd.PersistentFlags().StringVarP(&migrateDir, "dir", "r", "", "migration directory")
	downCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "dry run mode")
	downCmd.PersistentFlags().IntVarP(&max, "max", "m", 0, "limit of down migration")
	SchemaCmd.AddCommand(downCmd)

	RootCmd.AddCommand(seedCmd)
}

func PrintMigration(m *migrate.PlannedMigration, dir migrate.MigrationDirection) {
	fmt.Printf("==> Would apply migration %s (up)", m.Id)
	for _, q := range m.Up {
		fmt.Println(q)
	}
}

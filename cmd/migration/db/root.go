package db

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "migrate",
	Short: "DB Migrate Command",
}

var SchemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "DB Migrate Schema Command",
}

type Options struct {
	DB       string `envconfig:"POSTGRES_DB" required:"false"`
	User     string `envconfig:"POSTGRES_USER" required:"false" default:"admin"`
	Password string `envconfig:"POSTGRES_PASSWORD" required:"false"  default:"admin"`
	Host     string `envconfig:"POSTGRES_HOST" required:"false" default:"localhost"`
	Port     int    `envconfig:"POSTGRES_PORT" required:"false" default:"5432"`
	SSLMode  bool   `envconfig:"POSTGRES_SSL_MODE" required:"false" default:"false"`
	TimeZone string `envconfig:"POSTGRES_TIME_ZONE" required:"false" default:"Asia/Tokyo"`
}

var Opt Options

func init() {
	if err := envconfig.Process("", &Opt); err != nil {
		panic(err)
	}

	var params Options
	RootCmd.PersistentFlags().StringVarP(&params.Host, "host", "H", "", "db host")
	RootCmd.PersistentFlags().IntVarP(&params.Port, "port", "p", 0, "db port")
	RootCmd.PersistentFlags().StringVarP(&params.User, "user", "u", "", "db user")
	RootCmd.PersistentFlags().StringVarP(&params.Password, "password", "P", "", "db password")
	RootCmd.PersistentFlags().StringVarP(&params.DB, "dbname", "D", "", "db name")
	RootCmd.PersistentFlags().BoolVarP(&params.SSLMode, "ssl-mode", "s", false, "ssl mode")
	RootCmd.PersistentFlags().StringVarP(&params.TimeZone, "tz", "T", "Asia/Tokyo", "timezone")
	RootCmd.AddCommand(SchemaCmd)
	cobra.OnInitialize(func() {
		Opt.Host = getOverwriteString(params.Host, Opt.Host)
		Opt.Port = getOverwriteInt(params.Port, Opt.Port)
		Opt.User = getOverwriteString(params.User, Opt.User)
		Opt.Password = getOverwriteString(params.Password, Opt.Password)
		Opt.DB = getOverwriteString(params.DB, Opt.DB)
		Opt.SSLMode = getOverwriteBool(params.SSLMode, Opt.SSLMode)
		Opt.TimeZone = getOverwriteString(params.TimeZone, Opt.TimeZone)
	})
}

func getOverwriteString(v, defaultValue string) string {
	if v != "" {
		return v
	}
	return defaultValue
}

func getOverwriteInt(v, defaultValue int) int {
	if v != 0 {
		return v
	}
	return defaultValue
}

func getOverwriteBool(v, defaultValue bool) bool {
	if v != defaultValue {
		return v
	}
	return defaultValue
}

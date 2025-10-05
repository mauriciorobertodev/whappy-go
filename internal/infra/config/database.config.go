package config

type DatabaseDriver string

const (
	DatabaseDriverSQLite   DatabaseDriver = "sqlite"
	DatabaseDriverPostgres DatabaseDriver = "postgres"
)

func (d DatabaseDriver) IsValid() bool {
	switch d {
	case DatabaseDriverSQLite, DatabaseDriverPostgres:
		return true
	default:
		return false
	}
}

type DatabaseConfig struct {
	Driver DatabaseDriver
	DbName string
	DbUser string
	DbPass string
	DbHost string
	DbPort string
}

func (c *DatabaseConfig) GetDSN() string {
	switch c.Driver {
	case DatabaseDriverSQLite:
		if c.DbName == "" || c.DbName == ":memory:" {
			return ":memory:?_foreign_keys=on"
		}
		return "file:" + c.DbName + ".db?_foreign_keys=on"
	case DatabaseDriverPostgres:
		return "host=" + c.DbHost +
			" port=" + c.DbPort +
			" user=" + c.DbUser +
			" dbname=" + c.DbName +
			" password=" + c.DbPass +
			" sslmode=disable"
	default:
		panic("Unsupported Database Driver: " + c.Driver)
	}
}

func LoadWhappyDatabaseConfig() *DatabaseConfig {
	cfg := &DatabaseConfig{
		Driver: DatabaseDriver(GetEnvString("DB_DRIVER", "sqlite")),
		DbName: GetEnvString("DB_NAME", "whappy"),
		DbUser: GetEnvString("DB_USER", ""),
		DbPass: GetEnvString("DB_PASS", ""),
		DbHost: GetEnvString("DB_HOST", ""),
		DbPort: GetEnvString("DB_PORT", ""),
	}

	cfg.Validate()

	return cfg
}

func LoadWhatsmeowDatabaseConfig() *DatabaseConfig {
	cfg := &DatabaseConfig{
		Driver: DatabaseDriver(GetEnvString("WHATSMEOW_DB_DRIVER", "sqlite")),
		DbName: GetEnvString("WHATSMEOW_DB_NAME", "whatsmeow"),
		DbUser: GetEnvString("WHATSMEOW_DB_USER", ""),
		DbPass: GetEnvString("WHATSMEOW_DB_PASS", ""),
		DbHost: GetEnvString("WHATSMEOW_DB_HOST", ""),
		DbPort: GetEnvString("WHATSMEOW_DB_PORT", ""),
	}

	cfg.Validate()

	return cfg
}

func (c *DatabaseConfig) Validate() {
	if !c.Driver.IsValid() {
		panic("Invalid database driver: " + c.Driver)
	}

	if c.Driver == DatabaseDriverPostgres {
		if c.DbName == "" || c.DbUser == "" || c.DbPass == "" || c.DbHost == "" || c.DbPort == "" {
			panic("Invalid Postgres database config")
		}
	}

	if c.Driver == DatabaseDriverSQLite {
		if c.DbName == "" {
			panic("Invalid SQLite database config: DbName cannot be empty")
		}
	}
}

func (c *DatabaseConfig) CodeDriver() string {
	switch c.Driver {
	case DatabaseDriverSQLite:
		return "sqlite3"
	case DatabaseDriverPostgres:
		return "postgres"
	default:
		panic("Unsupported Database Driver: " + c.Driver)
	}
}

func (c *DatabaseConfig) IsSQLite() bool {
	return c.Driver == DatabaseDriverSQLite
}

func (c *DatabaseConfig) IsPostgres() bool {
	return c.Driver == DatabaseDriverPostgres
}

// connectionWithDataBase.go
package connectionWithDataBase

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

// Configuración de la estructura con valores fijos
type Config struct {
	DBUser                 string
	DBPass                 string
	DBName                 string
	InstanceConnectionName string
	PrivateIP              string
}

// Función para obtener la configuración con valores fijos
func getConfig() *Config {
	return &Config{
		DBUser:                 "postgres",
		DBPass:                 "cv-manager-db",
		DBName:                 "cvmanager",
		InstanceConnectionName: "cv-manager-432700:us-east1:cv-manager-db",
		PrivateIP:              "",
	}
}

// Función para conectar a la base de datos y devolver la conexión y un error si ocurre
func ConnectToDataBase() (*sql.DB, error) {
	config := getConfig()

	dsn := fmt.Sprintf("user=%s password=%s database=%s", config.DBUser, config.DBPass, config.DBName)
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgx.ParseConfig: %w", err)
	}

	var opts []cloudsqlconn.Option
	if config.PrivateIP != "" {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}

	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
	}

	cfg.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, config.InstanceConnectionName)
	}

	dbURI := stdlib.RegisterConnConfig(cfg)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Verificar la conexión
	if err := dbPool.Ping(); err != nil {
		return nil, fmt.Errorf("dbPool.Ping: %w", err)
	}

	return dbPool, nil
}

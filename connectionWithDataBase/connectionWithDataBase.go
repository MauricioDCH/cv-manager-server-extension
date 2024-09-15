package connectionWithDataBase

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"cv-manager-server-extension/config"

	"cloud.google.com/go/cloudsqlconn"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
)

// ConnectToDataBase establece una conexión con la base de datos y devuelve la conexión y un error si ocurre.
func ConnectToDataBase() (*sql.DB, error) {
	// Cargar configuración desde el archivo de configuración
	config, err := config.LoadConfig()
	if err != nil {
		return nil, fmt.Errorf("getConfig: %w", err)
	}

	// Crear el DSN (Data Source Name) para la base de datos
	dsn := fmt.Sprintf("user=%s password=%s database=%s", config.DBUser, config.DBPassword, config.DBName)
	cfg, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("pgx.ParseConfig: %w", err)
	}

	// Configurar opciones para Cloud SQL
	var opts []cloudsqlconn.Option
	if config.PrivateIP != "" {
		opts = append(opts, cloudsqlconn.WithDefaultDialOptions(cloudsqlconn.WithPrivateIP()))
	}

	// Crear un nuevo dialer para Cloud SQL
	d, err := cloudsqlconn.NewDialer(context.Background(), opts...)
	if err != nil {
		return nil, fmt.Errorf("cloudsqlconn.NewDialer: %w", err)
	}

	// Establecer la función de dial para conectar con Cloud SQL
	cfg.DialFunc = func(ctx context.Context, network, instance string) (net.Conn, error) {
		return d.Dial(ctx, config.InstanceConnectionName)
	}

	// Registrar la configuración de conexión y abrir la conexión con el pool
	// Asegúrate de que el URI de la base de datos esté correctamente formado.
	dbURI := stdlib.RegisterConnConfig(cfg)
	dbPool, err := sql.Open("pgx", dbURI)
	if err != nil {
		return nil, fmt.Errorf("sql.Open: %w", err)
	}

	// Verificar la conexión con un ping
	if err := dbPool.Ping(); err != nil {
		dbPool.Close() // Cerrar la conexión si el ping falla
		return nil, fmt.Errorf("dbPool.Ping: %w", err)
	}

	return dbPool, nil
}

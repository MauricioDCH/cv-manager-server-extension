package dbQueries

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"
)

// Columns representa la estructura de los campos de una tabla
type Columns struct {
	ColumnID   int    // ID autogenerado
	ColumnName string // Nombre de la columna
	DataType   string // Tipo de dato de la columna
}

type User struct {
	ID         int
	Nombre     string
	Correo     string
	Contraseña string
	CreatedAt  time.Time
}

// GetFieldsAsJSON obtiene los campos de una tabla y devuelve el resultado en formato JSON
func GetFieldsAsJSON(db *sql.DB, table string) (string, error) {
	fields, err := getFields(db, table)
	if err != nil {
		return "", err
	}

	columnsInfo := make([]map[string]string, 0, len(fields))
	for _, field := range fields {
		columnInfo := map[string]string{
			"ColumnField": field.ColumnName,
			"DataType":    field.DataType,
		}
		columnsInfo = append(columnsInfo, columnInfo)
	}

	jsonDataFields, err := json.MarshalIndent(columnsInfo, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error al convertir a JSON: %w", err)
	}

	return string(jsonDataFields), nil
}

// getFields obtiene los campos de una tabla y sus tipos de datos
func getFields(db *sql.DB, table string) ([]Columns, error) {
	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1", table)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var fields []Columns
	idCounter := 1
	for rows.Next() {
		var fieldName, dataType string
		if err := rows.Scan(&fieldName, &dataType); err != nil {
			return nil, err
		}

		field := Columns{
			ColumnID:   idCounter,
			ColumnName: fieldName,
			DataType:   dataType,
		}
		fields = append(fields, field)
		idCounter++
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return fields, nil
}

// Modificamos la función para aceptar una query como parámetro
func ReadUsers(db *sql.DB, query string) ([]User, error) {
	rows, err := db.Query(query) // Ejecutar la consulta que se pasa como parámetro
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []User
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Correo, &user.Contraseña); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

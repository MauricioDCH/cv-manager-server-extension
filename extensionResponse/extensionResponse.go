// extensionResponse.go
package extensionResponse

import (
	"encoding/json"
	"fmt"
)

// Estructura para representar cada campo en la respuesta JSON
type Field struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

// Estructura para representar el JSON completo
type ResponseToExtension struct {
	ResponseFromServer []Field `json:"responseFromServer"`
}

// Función para generar la respuesta JSON
func GenerateResponse(userEmail, userPassword string) ([]byte, error) {
	// Construir el JSON de respuesta
	responseToExtension := ResponseToExtension{
		ResponseFromServer: []Field{
			{
				Name:  "email",
				Type:  "text",
				Value: userEmail,
			},
			{
				Name:  "password",
				Type:  "pass",
				Value: userPassword,
			},
		},
	}

	// Convertir la respuesta a JSON
	responseJSON, err := json.Marshal(responseToExtension)
	if err != nil {
		return nil, fmt.Errorf("error al generar la respuesta JSON: %v", err)
	}

	return responseJSON, nil
}

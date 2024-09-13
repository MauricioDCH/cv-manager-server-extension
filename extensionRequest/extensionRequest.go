package extensionRequest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Estructura para representar cada input individual
type InputField struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

// Estructura que representa el JSON completo recibido
type RequestData struct {
	Inputs json.RawMessage `json:"inputs"`
	Email  string          `json:"email"`
}

// Función que maneja las solicitudes POST y devuelve los datos
func HandlePostRequest(responseWriter http.ResponseWriter, request *http.Request) (*RequestData, error) {
	bodyBytes, err := io.ReadAll(request.Body)
	if err != nil {
		return nil, fmt.Errorf("error al leer el cuerpo de la solicitud: %v", err)
	}

	var requestData RequestData
	err = json.Unmarshal(bodyBytes, &requestData)
	if err != nil {
		return nil, fmt.Errorf("error al decodificar el JSON: %v", err)
	}

	// Compactar el JSON crudo en `inputs` para mostrarlo en una sola línea
	var compactedInputs bytes.Buffer
	err = json.Compact(&compactedInputs, requestData.Inputs)
	if err != nil {
		return nil, fmt.Errorf("error al compactar el JSON de inputs: %v", err)
	}

	//fmt.Println("Correo electrónico (fuera de inputs):", requestData.Email)
	//fmt.Println("Inputs (JSON crudo en una línea):", compactedInputs.String())

	var compactedBody bytes.Buffer
	err = json.Compact(&compactedBody, bodyBytes)
	if err != nil {
		return nil, fmt.Errorf("error al compactar el JSON de la solicitud: %v", err)
	}

	responseWriter.Header().Set("Content-Type", "application/json")
	responseWriter.Write(compactedBody.Bytes())

	return &requestData, nil
}

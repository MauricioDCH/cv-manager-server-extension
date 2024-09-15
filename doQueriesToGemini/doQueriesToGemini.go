package doQueriesToGemini

import (
	"context"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

// FirstQuery realiza la primera consulta a Gemini para identificar qué inputs necesitan un valor
// y también están en la lista de campos disponibles.
func FirstQuery(ctx context.Context, model *genai.GenerativeModel, jsonDataInputs, jsonFieldsData string) (string, error) {
	// Construir el mensaje para la consulta
	message1 := fmt.Sprintf(
		"Necesito saber qué inputs de la siguiente lista requieren un valor: %s. "+
			"Estos inputs deben estar presentes en esta otra lista de campos disponibles: %s. "+
			"Por favor, enumera solo los inputs que coincidan, sin texto adicional.",
		jsonDataInputs, jsonFieldsData)

	// Realizar la consulta a Gemini
	responseInputs, err := model.GenerateContent(ctx, genai.Text(message1))
	if err != nil {
		return "", err
	}

	// Procesar la respuesta de Gemini
	responseInputsForQuery := returnResponse(responseInputs)
	return responseInputsForQuery, nil
}

// SecondQuery realiza la segunda consulta a Gemini para obtener una consulta SQL específica
func SecondQuery(ctx context.Context, model *genai.GenerativeModel, tableToQuery, responseInputsForQuery, email string) (string, error) {
	// Construir el mensaje para la consulta SQL

	message2 := fmt.Sprintf(
		"Necesito una consulta SQL para PostgreSQL. La consulta debe ser para la tabla %s "+
			" y los campos para la petición a la base de datos debe ser SOLO los campos que muestre esta respuesta: %s."+
			" y el id. La búsqueda debe ser para el correo %s"+
			"Devuelveme SOLO la consulta SQL completa sin ningún formato adicional ni texto extra. OJO: Los campos de la query deben ser los ColumnField del json llamado inputsTabla.",
		tableToQuery, responseInputsForQuery, email)

	/*
	 */
	// Realizar la consulta a Gemini
	mergeForQuery, err := model.GenerateContent(ctx, genai.Text(message2))
	if err != nil {
		return "", err
	}

	// Obtener la respuesta como un string
	query := returnResponseQuery(mergeForQuery)
	return query, nil
}

// returnResponse procesa la respuesta de Gemini y devuelve el texto completo.
func returnResponse(resp *genai.GenerateContentResponse) string {
	var responseText strings.Builder

	// Recorrer los candidatos y partes para construir la respuesta
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			// Convertir `part` a texto si es de tipo `genai.Text`
			if text, ok := part.(genai.Text); ok {
				responseText.WriteString(string(text))
			} else {
				log.Printf("Unexpected type: %T\n", part)
			}
		}
	}

	// Obtener la respuesta completa como string
	fullResponse := responseText.String()
	return fullResponse
}

// returnResponseQuery procesa la respuesta para extraer solo la consulta SQL
func returnResponseQuery(resp *genai.GenerateContentResponse) string {
	var responseText strings.Builder

	// Recorrer los candidatos y partes para construir la respuesta
	for _, candidate := range resp.Candidates {
		for _, part := range candidate.Content.Parts {
			// Convertir `part` a texto si es de tipo `genai.Text`
			if text, ok := part.(genai.Text); ok {
				responseText.WriteString(string(text))
			} else {
				log.Printf("Unexpected type: %T\n", part)
			}
		}
	}

	// Obtener la respuesta completa como string
	fullResponse := responseText.String()

	// Limpiar la respuesta para dejar solo la consulta SQL
	cleanedResponse := strings.TrimSpace(fullResponse)

	// Usar una expresión regular para buscar una consulta SQL en la respuesta
	re := regexp.MustCompile(`(?i)SELECT\s+.*?\s+FROM\s+\S+.*;`)
	match := re.FindString(cleanedResponse)

	// Devolver la consulta SQL si se encuentra; de lo contrario, devolver la respuesta completa
	if match != "" {
		return match
	}

	return cleanedResponse
}

package doQueriesToGemini

import (
	"context"
	"log"
	"strings"

	"github.com/google/generative-ai-go/genai"
)

// FirstQuery realiza la primera consulta a Gemini para identificar qué inputs necesitan un valor
// y también están en la lista de campos disponibles.
func FirstQuery(ctx context.Context, model *genai.GenerativeModel, jsonDataInputs, jsonFieldsData string) (string, error) {
	// Construir el mensaje para la consulta
	message1_part_1 := "Necesito que me digas cuales de estos inputs " + jsonDataInputs
	message1_part_2 := " necesitan un valor y también están en esta otra lista " + jsonFieldsData
	message1_part_3 := " para poder llenarlos, no es necesario una pregunta larga."
	message1 := message1_part_1 + message1_part_2 + message1_part_3

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
	message2_part_1 := "Necesito una query SQL para PostgreSQL. La consulta debe ser para la tabla "
	message2_part_2 := tableToQuery
	message2_part_3 := " SÓLO con los campos que coincidan con la respuesta y el id: " + responseInputsForQuery
	message2_part_4 := ", y debe buscar el correo "
	message2_part_5 := email
	message2_part_6 := ". Por favor, devuelve SOLO la consulta SQL completa usando los datos en español, sin ningún formato adicional o texto extra."
	message2 := message2_part_1 + message2_part_2 + message2_part_3 + message2_part_4 + message2_part_5 + message2_part_6

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
	// Eliminar cualquier texto no deseado antes y después de la consulta
	cleanedResponse := strings.TrimSpace(fullResponse)

	// Buscar la parte de la consulta SQL
	lines := strings.Split(cleanedResponse, "\n")
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "SELECT") {
			return strings.TrimSpace(line)
		}
	}

	// Si no se encuentra la consulta SQL específica, devolver la respuesta completa
	return cleanedResponse
}

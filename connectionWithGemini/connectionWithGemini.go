// connectionWithGemini.go
package connectionWithGemini

import (
	"context"
	"log"

	"github.com/google/generative-ai-go/genai" // Dependencias para Google Generative AI
	"google.golang.org/api/option"
)

// Función para conectarse a Gemini
func ConnectToGemini() (*genai.Client, context.Context) {
	// Crear un contexto
	ctx := context.Background()

	// Clave API de Google Generative AI
	apiKey := "AIzaSyAWOVQpKOtCnjGq22aXznUPJdyn2Upk7iE" // Reemplaza con tu clave de API
	if apiKey == "" {
		log.Fatal("La variable de entorno API_KEY no está configurada")
	}

	// Crear el cliente con la clave API
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		log.Fatal("Error al crear el cliente de Gemini:", err)
	}

	// Retornar el cliente y el contexto para su uso posterior
	return client, ctx
}

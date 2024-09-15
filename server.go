// server.go
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	dbConnection "cv-manager-server-extension/connectionWithDataBase"
	connectionWithGemini "cv-manager-server-extension/connectionWithGemini"
	dbQueries "cv-manager-server-extension/doQueriesToDB"
	doQueriesToGemini "cv-manager-server-extension/doQueriesToGemini"
	extensionRequest "cv-manager-server-extension/extensionRequest"
	extensionResponse "cv-manager-server-extension/extensionResponse"
	service "cv-manager-server-extension/security"
)

func main() {
	// Ruta para manejar las solicitudes POST en "/submit-data"
	http.HandleFunc("/submit-data", func(responseWriter http.ResponseWriter, request *http.Request) {
		// Llama a HandlePostRequest y maneja el error si ocurre
		requestData, err := extensionRequest.HandlePostRequest(responseWriter, request)
		if err != nil {
			http.Error(responseWriter, "Error al procesar la solicitud", http.StatusInternalServerError)
			return
		}

		if requestData == nil {
			http.Error(responseWriter, "Datos de solicitud vacíos", http.StatusBadRequest)
			return
		}

		// Extraer el JSON de inputs y el email
		inputsJSON := requestData.Inputs
		email := requestData.Email

		// Convertir inputsJSON a una cadena JSON compacta
		var compactedInputs bytes.Buffer
		if err := json.Compact(&compactedInputs, inputsJSON); err != nil {
			http.Error(responseWriter, "Error al compactar el JSON de inputs", http.StatusInternalServerError)
			return
		}

		// Mostrar los valores obtenidos
		fmt.Print("\nSOLICITUD POST RECIBIDA.\n\n")
		fmt.Println("Json de inputs:\t\t", compactedInputs.String())
		fmt.Println("Correo electrónico:\t", email)

		// Convertir inputsJSON a una cadena JSON antes de usarlo en la consulta
		inputsJSONStr, err := json.Marshal(inputsJSON)
		if err != nil {
			http.Error(responseWriter, "Error al convertir JSON a cadena", http.StatusInternalServerError)
			return
		}

		// Consultar la base de datos
		fmt.Print("\nCONECTANDO A LA BASE DE DATOS DE GPC...\n")
		db, err := dbConnection.ConnectToDataBase()
		if err != nil {
			http.Error(responseWriter, "Error al conectar a la base de datos", http.StatusInternalServerError)
			return
		}
		defer db.Close()

		fmt.Print(" \nCONEXIÓN EXITOSA A LA BASE DE DATOS!\n")
		tableToQuery := "users"
		// Llama a la función GetFieldsAsJSON
		jsonFields, err := dbQueries.GetFieldsAsJSON(db, tableToQuery)
		if err != nil {
			http.Error(responseWriter, "Error al obtener campos de la base de datos", http.StatusInternalServerError)
			return
		}

		fmt.Printf("\nCAMPOS DE LA TABLA '%s' EN FORMATO JSON:\n\n%s\n", tableToQuery, jsonFields)

		// Consultar Gemini
		fmt.Print("\nCONECTANDO A GEMINI...\n")
		client, ctx := connectionWithGemini.ConnectToGemini()
		defer client.Close()

		// Crear el modelo generativo
		model := client.GenerativeModel("gemini-1.5-flash")

		fmt.Print("CONEXIÓN EXITOSA A GEMINI!\nCREANDO LA PRIMERA CONSULTA...\nENVIANDO LA PRIMERA CONSULTA A GEMINI...\n")
		// Llama a la nueva función FirstQuery desde doQueriesToGemini
		responseInputsForQuery, err := doQueriesToGemini.FirstQuery(ctx, model, string(inputsJSONStr), jsonFields)
		if err != nil {
			http.Error(responseWriter, "Error en la consulta a Gemini", http.StatusInternalServerError)
			return
		}

		// Imprimir la respuesta de los inputs
		fmt.Print("\nRESPUESTA DE GEMINI A LA PRIMERA PETICIÓN A GEMINI:\n")
		fmt.Print("---------------------------------------------------\n")
		fmt.Println("\n" + responseInputsForQuery)

		fmt.Print("----------------------------------------\n")

		fmt.Print("ENVIANDO LA SEGUNDA CONSULTA A GEMINI...\n")
		fmt.Print("---------------------------------------------------\n")
		// Realizar la segunda consulta a Gemini
		query, err := doQueriesToGemini.SecondQuery(ctx, model, tableToQuery, responseInputsForQuery, email)
		if err != nil {
			http.Error(responseWriter, "Error en la segunda consulta a Gemini", http.StatusInternalServerError)
			return
		}

		// Imprimir la consulta SQL obtenida
		fmt.Print("\nRESPUESTA DE GEMINI A LA SEGUNDA PETICIÓN A GEMINI:\n")
		fmt.Print("----------------------------------------------------\n")
		fmt.Println("\nEsta es la query que se va a realizar: --> ||||\t\t" + query + "\t|||| <--")
		fmt.Print("----------------------------------------------------\n")

		// Realizar la consulta a la base de datos.
		users, err := dbQueries.ReadUsers(db, query)
		if err != nil {
			http.Error(responseWriter, "Error al leer usuarios de la base de datos", http.StatusInternalServerError)
			return
		}

		// Desencriptar la contraseña
		for _, user := range users {
			decryptedPassword, err := service.Decrypt(user.Contraseña, db)
			if err != nil {
				http.Error(responseWriter, "Error al desencriptar la contraseña", http.StatusInternalServerError)
				return
			}
			fmt.Printf("Respuesta de la petición a la base de datos: --> |||\t\tID: %d, Correo: %s, Contraseña: %s\t\t\n", user.ID, user.Correo, decryptedPassword)
		}

		// Generar el JSON de respuesta
		var responseJSON []byte
		for _, user := range users {
			responseJSON, err = extensionResponse.GenerateResponse(user.Correo, user.Contraseña)
			if err != nil {
				http.Error(responseWriter, "Error al generar la respuesta JSON", http.StatusInternalServerError)
				return
			}
		}

		// Enviar la respuesta JSON
		fmt.Print("\nENVIANDO RESPUESTA A LA EXTENSIÓN...\n")
		fmt.Printf("Respuesta JSON:\t\t%s\n\n", string(responseJSON))
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write(responseJSON)
		fmt.Print("\nRESPUESTA ENVIADA.\nFIN DE PROCESO\n")
		fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Println("----------------------------------------------------------------------------------------------------------------------------------------------------")
		fmt.Print("\nESPERANDO OTRA SOLICITUD DESDE EXTENSIÓN...\n\n")
	})

	// Iniciar el servidor en el puerto 5522
	serverAddress := ":5522"
	fmt.Printf("Servidor corriendo en http://localhost%s\n", serverAddress)
	fmt.Print("\nESPERANDO SOLICITUD DESDE EXTENSIÓN...\n\n")
	log.Fatal(http.ListenAndServe(serverAddress, nil))
}

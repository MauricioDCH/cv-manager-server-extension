package service

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"

	"cv-manager-server-extension/config"
)

// Función para obtener la clave cifrada desde la base de datos
func getEncryptedKeyFromDB(db *sql.DB) (string, error) {
	var encryptedKey string
	query := "SELECT key FROM encrypted_keys LIMIT 1" // Ajusta la consulta según tu necesidad
	row := db.QueryRow(query)
	if err := row.Scan(&encryptedKey); err != nil {
		return "", err
	}
	return encryptedKey, nil
}

// Función para desencriptar la clave cifrada utilizando la clave maestra
func decryptKey(encryptedKey string, masterKey []byte) ([]byte, error) {
	ciphertext, err := hex.DecodeString(encryptedKey)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(masterKey)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < aes.BlockSize {
		return nil, errors.New("ciphertext too short")
	}

	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	return ciphertext, nil
}

// Función para desencriptar el texto usando la clave desencriptada
func Decrypt(encryptedText string, db *sql.DB) (string, error) {
	// Cargar la configuración desde el archivo de configuración
	configs, err := config.LoadConfig()
	if err != nil {
		return "", fmt.Errorf("error loading configuration: %w", err)
	}

	// Obtener la clave maestra desde la configuración
	masterKeyBase64 := configs.MasterKey
	if len(masterKeyBase64) == 0 {
		return "", errors.New("master key not set in configuration")
	}

	// Decodificar la clave maestra desde Base64
	masterKey, err := base64.StdEncoding.DecodeString(masterKeyBase64)
	if err != nil {
		return "", fmt.Errorf("error decoding master key: %w", err)
	}

	// Recuperar la clave cifrada desde la base de datos
	encryptedKey, err := getEncryptedKeyFromDB(db)
	if err != nil {
		return "", err
	}

	// Desencriptar la clave cifrada para obtener la clave real
	decryptedKey, err := decryptKey(encryptedKey, masterKey)
	if err != nil {
		return "", err
	}

	// Decodificar el texto cifrado en hexadecimal
	ciphertext, err := hex.DecodeString(encryptedText)
	if err != nil {
		return "", err
	}

	// Crear el cifrador AES usando la clave desencriptada
	block, err := aes.NewCipher(decryptedKey)
	if err != nil {
		return "", err
	}

	// Verificar que el texto cifrado sea al menos del tamaño del bloque AES
	if len(ciphertext) < aes.BlockSize {
		return "", errors.New("ciphertext too short")
	}

	// Extraer el IV del texto cifrado (primer bloque)
	iv := ciphertext[:aes.BlockSize]
	ciphertext = ciphertext[aes.BlockSize:]

	// Crear el flujo de desencriptación en modo CFB usando el bloque y el IV
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(ciphertext, ciphertext)

	// Imprimir valores para depuración
	fmt.Printf("Master Key: %x\n", masterKey)
	fmt.Printf("Encrypted Key: %s\n", encryptedKey)
	fmt.Printf("Decrypted Key: %x\n", decryptedKey)
	fmt.Printf("Ciphertext: %x\n", ciphertext)

	// Devolver el texto desencriptado como string
	return string(ciphertext), nil
}

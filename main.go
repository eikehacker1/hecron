package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// User-Agent global
const userAgent = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/96.0.4664.45 Safari/537.36"

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: script.go arquivo.txt")
		return
	}

	filename := os.Args[1]

	// Cria o diretório 'out' se ele não existir
	err := os.MkdirAll("out", os.ModePerm)
	if err != nil {
		fmt.Println("Erro ao criar o diretório 'out':", err)
		return
	}

	processFile(filename)
}

func processFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Erro ao abrir o arquivo %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		url := scanner.Text()
		processURL(url)
	}
	if err := scanner.Err(); err != nil {
		fmt.Printf("Erro ao ler o arquivo %s: %v\n", filename, err)
	}
}

func processURL(url string) {
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		url = "http://" + url
	}

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("Erro ao criar a requisição para %s: %v\n", url, err)
		return
	}

	// Defina o User-Agent global
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Erro ao fazer a requisição para %s: %v\n", url, err)
		return
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Erro ao ler o corpo da resposta para %s: %v\n", url, err)
		return
	}

	// Cria um arquivo com a extensão .head no diretório 'out'
	filename := "out/" + strings.Replace(url, "://", ".", -1) + ".head"

	// Crie uma string com o cabeçalho HTTP
	headerString := "Status Code: " + resp.Status + "\n"
	for key, values := range resp.Header {
		for _, value := range values {
			headerString += key + ": " + value + "\n"
		}
	}

	// Combine o cabeçalho com o corpo da resposta
	fileContent := headerString + "\n" + string(body)

	err = ioutil.WriteFile(filename, []byte(fileContent), 0644)
	if err != nil {
		fmt.Printf("Erro ao salvar o arquivo %s: %v\n", filename, err)
	} else {
		fmt.Printf("Salvo o conteúdo de %s em %s\n", url, filename)
	}
}

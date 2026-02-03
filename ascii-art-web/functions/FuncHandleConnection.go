package artweb

import (
	artweb "ascii-art-web/functions/functionsascii"
	"bufio"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"strconv"
	"strings"
)

func HandleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// Lire la première ligne : ex. "GET / HTTP/1.1"
	requestLine, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Erreur de lecture:", err)
		return
	}
	fmt.Print("Requête: ", requestLine)

	// Extraire méthode et chemin
	parts := strings.Fields(requestLine)
	if len(parts) < 2 {
		return
	}
	method := parts[0]
	path := parts[1]

	// Lire les headers pour récupérer Content-Length si POST
	var contentLength int
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}
		if line == "\r\n" {
			break // fin des headers
		}
		// Chercher Content-Length
		if strings.HasPrefix(line, "Content-Length:") {
			cl := strings.TrimSpace(strings.Split(line, ":")[1])
			contentLength, _ = strconv.Atoi(cl)
		}
	}

	// Si GET /
	if method == "GET" && path == "/" {
		// Lire le fichier HTML (formulaire)
		bodyBytes, err := os.ReadFile("index.html")
		if err != nil {
			fmt.Println("Erreur lecture fichier HTML:", err)
			return
		}
		body := string(bodyBytes)
		response := "HTTP/1.1 200 OK\r\n" +
			"Content-Length: " + strconv.Itoa(len(body)) + "\r\n" +
			"Content-Type: text/html\r\n" +
			"Connection: close\r\n\r\n" +
			body
		conn.Write([]byte(response))
		return
	}

	if method == "GET" && path == "/favicon.ico" {
		iconBytes, err := os.ReadFile("favicon.ico")
		if err != nil {
			// Fallback si le fichier est manquant
			response := "HTTP/1.1 204 No Content\r\n" +
				"Connection: close\r\n\r\n"
			conn.Write([]byte(response))
			return
		}
		response := "HTTP/1.1 200 OK\r\n" +
			"Content-Length: " + strconv.Itoa(len(iconBytes)) + "\r\n" +
			"Content-Type: image/x-icon\r\n" +
			"Connection: close\r\n\r\n"
		conn.Write([]byte(response))
		conn.Write(iconBytes)
		return
	}

	// Si POST /ascii-art
	if method == "POST" && path == "/ascii-art" {
		// Lire le corps POST
		bodyData := make([]byte, contentLength)
		_, err := io.ReadFull(reader, bodyData)
		fmt.Println(string(bodyData))
		if err != nil {
			fmt.Println("Erreur de lecture du body POST:", err)
			return
		}

		// Parse les données de formulaire
		params, err := url.ParseQuery(string(bodyData))
		if err != nil {
			fmt.Println("Erreur parse query:", err)
			return
		}

		text := params.Get("text")
		banner := params.Get("banner") + ".txt"
		// Générer une réponse simple (tu peux insérer du vrai ASCII art ici)
		asciiart := artweb.AsciiArt(text, banner)
		var responseBody string
		for i := 0; i < len(asciiart); i++ {
			for j := 0; j < len(asciiart[i]); j++ {
				responseBody = responseBody + asciiart[i][j] + "\n"
			}
		}

		response := "HTTP/1.1 200 OK\r\n" +
			"Content-Length: " + strconv.Itoa(len(responseBody)) + "\r\n" +
			"Content-Type: text/plain\r\n" +
			"Connection: close\r\n\r\n" +
			responseBody

		conn.Write([]byte(response))
		return
	}

	// Sinon, 404 Not Found
	notFound := "404 - Page non trouvée"
	response := "HTTP/1.1 404 Not Found\r\n" +
		"Content-Length: " + strconv.Itoa(len(notFound)) + "\r\n" +
		"Content-Type: text/plain\r\n" +
		"Connection: close\r\n\r\n" +
		notFound
	conn.Write([]byte(response))
}

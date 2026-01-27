package main

import (
	artweb "ascii-art-web/functions"
	"fmt"
	"net"
)

func main() {
	// Étape 1 : Créer un listener TCP sur le port 8080
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	fmt.Println("Serveur en écoute sur le port 8080...")

	// Étape 2 : Boucle infinie pour accepter les connexions
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Erreur de connexion:", err)
			continue
		}
		go artweb.HandleConnection(conn) // Goroutine pour gérer la connexion
	}
}

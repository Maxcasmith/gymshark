package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	cors "gymshark/server/middleware/Cors"
)

func main() {
	address := flag.String("addr", ":8080", "listen address")
	dbDir := flag.String("db", "../../db", "directory holding packs.json and orders.json")
	allowedOrigins := flag.String("origins", "http://localhost:5173", "comma-separated origins allowed by CORS")
	flag.Parse()

	handler := cors.Cors(strings.Split(*allowedOrigins, ","), routes(*dbDir))

	log.Printf("listening on %s (db: %s, origins: %s)", *address, *dbDir, *allowedOrigins)
	log.Fatal(http.ListenAndServe(*address, handler))
}

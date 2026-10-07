package cors

import (
	"net/http"
	"slices"
)

func Cors(allowedOrigins []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		origin := request.Header.Get("Origin")
		writer.Header().Add("Vary", "Origin")

		if slices.Contains(allowedOrigins, origin) {
			writer.Header().Set("Access-Control-Allow-Origin", origin)
			writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, OPTIONS")
			writer.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		}

		if request.Method == http.MethodOptions && request.Header.Get("Access-Control-Request-Method") != "" {
			writer.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(writer, request)
	})
}

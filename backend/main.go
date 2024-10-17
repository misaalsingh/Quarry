package main

import (
    "log"
    "net/http"
    "backend/routes"
    "backend/db"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

// Add the debug function here, at package level
func debugRoutes(router *mux.Router) {
    router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        pathTemplate, err := route.GetPathTemplate()
        if err == nil {
            methods, _ := route.GetMethods()
            log.Printf("Route: %s, Methods: %v", pathTemplate, methods)
        }
        return nil
    })
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Printf("Received request: %s %s", r.Method, r.URL.Path)
        next.ServeHTTP(w, r)
    })
}

func main() {
    // Setup CockroachDB connection if needed (for static uses)
    dbConn, err := db.ConnectDB()
    if err != nil {
        log.Fatal("Failed to connect to CockroachDB:", err)
    }

    // Initialize the router
    router := mux.NewRouter()
	

    // Apply logging middleware
    router.Use(loggingMiddleware)

    // Setup routes
    routes.SetupRoutes(router, dbConn)

    // Add debug call here, after routes are set up
    log.Println("Registered routes:")
    debugRoutes(router)

    // Setup CORS
    c := cors.New(cors.Options{
        AllowedOrigins: []string{"*"}, // Adjust this based on your needs
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"*"},
    })
    handler := c.Handler(router)

    // Start the HTTP server
    log.Println("Server starting on port 8080")
    log.Fatal(http.ListenAndServe(":8080", handler))
}
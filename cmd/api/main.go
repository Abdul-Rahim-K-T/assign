package main

import (
	"log"
	"net/http"
	httpRoute "recruitment-management/internal/user/delivery/http" // Import the common user route handler
	"recruitment-management/pkg/database"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize the database
	database.InitDB()

	// Create a new router
	r := mux.NewRouter()

	// User routes (both Admin and Applicant)
	httpRoute.SetupUserRoutes(r)

	// Start the server
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}

package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"recruitment-management/internal/models"
	"recruitment-management/pkg/jwt"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

// Handle job creation
func (h *UserHandler) CreateJob(w http.ResponseWriter, r *http.Request) {
	var job models.Job

	// Decode the incoming JSON request body into a job struct
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Extract the user_id from JWT claims in the request context
	claims := jwt.GetClaims(r) // This retrieves the claims from the request context
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Set the `posted_by_id` from the JWT claims
	job.PostedByID = claims.UserID

	// set POstedOn field to current time
	job.PostedOn = time.Now()

	// Initialize total applicantion to 0
	job.TotalApplications = 0

	var user models.User
	err = h.Usecase.GetUserByID(claims.UserID, &user)
	if err != nil {
		http.Error(w, "Failed to fectch user details", http.StatusInternalServerError)
		return
	}

	// Populate the `PostedBy` field with user details
	job.PostedBy = user

	// Call usecase to create the job opening
	err = h.Usecase.CreateJob(job)
	if err != nil {
		http.Error(w, "Failed to create job opening", http.StatusInternalServerError)
		return
	}

	// Send a success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(job)
}

// Handle job details fetching
func (h *UserHandler) GetJobDetails(w http.ResponseWriter, r *http.Request) {
	// Get job_id from the url parameters
	vars := mux.Vars(r)
	jobID, err := strconv.Atoi(vars["job_id"])
	if err != nil {
		http.Error(w, "Invalid job ID", http.StatusBadRequest)
		return
	}

	// Fetch job details using the usecase
	job, applicants, err := h.Usecase.GetJobDetails(jobID)
	if err != nil {
		http.Error(w, "Failed to fetch job details", http.StatusInternalServerError)
		return
	}

	// Prepare the response
	response := struct {
		Job        models.Job    `json:"job"`
		Applicants []models.User `json:"applicants"`
	}{
		Job:        job,
		Applicants: applicants,
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Handle fetching all users
func (h *UserHandler) GetAllUsers(w http.ResponseWriter, r *http.Request) {
	// Fetch all users using the usecase
	users, err := h.Usecase.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// GetApplicantByID fetches the details of a spcifc applicant
func (h *UserHandler) GetApplicantByID(w http.ResponseWriter, r *http.Request) {
	// Extract applicant_id from URL parameters
	vars := mux.Vars(r)
	applicantIDstr := vars["applicant_id"]

	// Convert the applicant_id string to an int
	applicantID, err := strconv.Atoi(applicantIDstr)
	if err != nil {
		http.Error(w, "Invalid applicant ID", http.StatusBadRequest)
		return
	}

	// Calll the usecase to get the applicant's details by id
	var applicant models.User
	err = h.Usecase.GetProfileByID(applicantID, &applicant)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error fetching applicant: %v", err), http.StatusInternalServerError)
		return
	}

	// Check if the user is an applicant
	if applicant.UserType == "Admin" {
		http.Error(w, "This is an admin account, not an applicant. Please try another ID.", http.StatusBadRequest)
		return
	}

	// Return the applicant's data as a JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applicant)

}

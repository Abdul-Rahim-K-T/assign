package http

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"recruitment-management/internal/models"
	"recruitment-management/internal/services"
	"recruitment-management/internal/user/usecase"
	"recruitment-management/pkg/database"
	"recruitment-management/pkg/jwt"
	"strconv"
	"strings"
	"time"
)

type UserHandler struct {
	Usecase usecase.UserUsecase
}

// Signup handles the user signup for both Admin and Applicant
func (h *UserHandler) Signup(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Check if the user type is valid (either Admin or Applicant)
	if user.UserType != "admin" && user.UserType != "applicant" {
		http.Error(w, "Invalid user type", http.StatusBadRequest)
		return
	}

	// Set creation timestamp
	user.CreatedAt = time.Now()

	// Call the usecase to create the user
	if err := h.Usecase.CreateUser(user); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"})
}

// Login handles user login for both Admin and Applicant
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var loginData struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&loginData); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Call the usecase to authenticate the user and generate a token
	token, err := h.Usecase.Login(loginData.Email, loginData.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Respond with the JWT token
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func (h *UserHandler) uploadResume(w http.ResponseWriter, r *http.Request) {
	// Get the user claims from the context
	claims := jwt.GetClaims(r)
	if claims == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse the multipart form
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Retrieve file
	file, handler, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate file type
	ext := filepath.Ext(handler.Filename)
	if ext != ".pdf" && ext != ".docx" {
		http.Error(w, "Invalid file format. Only PDF or DOCX allowed", http.StatusBadRequest)
		return
	}

	// Define the file path
	dirPath := "uploads/resumes/"
	filePath := fmt.Sprintf("%s%s%s", dirPath, fmt.Sprintf("%v", claims.UserID), ext)

	// Create the directory if it doesn't exist
	if _, err := os.Stat(dirPath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirPath, os.ModePerm); err != nil {
			http.Error(w, "Failed to create directory", http.StatusInternalServerError)
			return
		}
	}

	// Save the file
	out, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Process the resume using APILayer API
	resumeData, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, "Failed to read resume file", http.StatusInternalServerError)
		return
	}

	parsedResume, err := services.ParseResumeWithAPILayer(resumeData)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to parse resume: %v", err), http.StatusInternalServerError)
		return
	}
	log.Println("Resume parsed successfully")

	// Serialize experience to JSON
	experienceJSON, err := json.Marshal(parsedResume.Experience)
	if err != nil {
		http.Error(w, "Failed to serialize experience", http.StatusInternalServerError)
		return
	}

	educationJSON, err := json.Marshal(parsedResume.Education)
	if err != nil {
		http.Error(w, "Failed to serialize education", http.StatusInternalServerError)
		return
	}

	// Count the number of profiles for the given user ID
	var profileCount int
	err = database.DB.QueryRow("SELECT COUNT(*) FROM profiles WHERE user_id = $1", claims.UserID).Scan(&profileCount)
	if err != nil {
		log.Printf("Failed to check profile existence: %v", err)
		http.Error(w, "Failed to check profile existence", http.StatusInternalServerError)
		return
	}

	// Create or update profile based on existence
	if profileCount == 0 {
		// No profile found, create a new one
		_, err := database.DB.Exec(
			`INSERT INTO profiles (user_id, resume_file, skills, education, experience, phone) 
			VALUES ($1, $2, $3, $4, $5, $6)`,
			claims.UserID, filePath, joinSkills(parsedResume.Skills), educationJSON,
			experienceJSON, parsedResume.Phone,
		)
		if err != nil {
			fmt.Printf("Failed to create profile: %v", err)
			http.Error(w, "Failed to create profile", http.StatusInternalServerError)
			return
		}
	} else {
		// Profile found, update it
		_, err := database.DB.Exec(
			`UPDATE profiles SET resume_file = $1, skills = $2, education = $3, 
			experience = $4, name = $5, email = $6, phone = $7 
			WHERE user_id = $8`,
			filePath, joinSkills(parsedResume.Skills), educationJSON, experienceJSON,
			parsedResume.Name, parsedResume.Email, parsedResume.Phone, claims.UserID,
		)
		if err != nil {
			fmt.Printf("Failed to update profile: %v", err)
			http.Error(w, "Failed to update profile", http.StatusInternalServerError)
			return
		}
	}

	// Respond with success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Resume uploaded successfully"})
}

// Function to join skills slice into a single string
func joinSkills(skills []string) string {
	return strings.Join(skills, ", ")
}

func (h *UserHandler) GetJobs(w http.ResponseWriter, r *http.Request) {
	// Call the usecase to get the list of jobs
	jobs, err := h.Usecase.GetJobs()
	if err != nil {
		http.Error(w, "Error fetching jobs", http.StatusInternalServerError)
		return
	}

	// return the jobs as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

// applyfor job handles the request to apply for a job
func (h *UserHandler) ApplyForJob(w http.ResponseWriter, r *http.Request) {
	// Extract job_id from the query parameters
	jobIDStr := r.URL.Query().Get("job_id")
	if jobIDStr == "" {
		http.Error(w, "Job ID is required", http.StatusBadRequest)
		return
	}

	jobID, err := strconv.Atoi(jobIDStr)
	if err != nil {
		http.Error(w, "Invalid Job ID", http.StatusBadRequest)
		return
	}

	// Get the applicant's user ID from the claims
	claims := jwt.GetClaims(r)
	if claims == nil || claims.UserType != "applicant" {
		http.Error(w, "Access forbidden: Only applicants can apply for jobs", http.StatusForbidden)
		return
	}

	// Check if the applicant has already applied for the job

	existingApplication, err := h.Usecase.GetApplicationByUserAndJob(strconv.Itoa(claims.UserID), strconv.Itoa(jobID))
	if err != nil {
		http.Error(w, "Failed to check existing application", http.StatusInternalServerError)
		return
	}

	if existingApplication != nil {
		http.Error(w, "You have already applied for this job", http.StatusBadRequest)
		return
	}

	// Call the usecase to apply for the job
	application := &models.Application{
		UserID: claims.UserID,
		JobID:  jobID,
	}
	if err := h.Usecase.ApplyForJob(application); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Update the toatal applications count for the job
	err = h.Usecase.UpdateJobApplicationsCount(jobIDStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return success response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Successfully applied for the job",
	})
}

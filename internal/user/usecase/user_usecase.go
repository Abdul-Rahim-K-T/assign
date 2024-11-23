package usecase

import (
	"errors"
	"fmt"
	"log"
	"recruitment-management/internal/models"
	"recruitment-management/internal/services"
	"recruitment-management/internal/user/repository/postgres"
	"recruitment-management/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type UserUsecase interface {
	CreateUser(user models.User) error
	Login(email, password string) (string, error)
	GetJobs() ([]models.Job, error)
	GetApplicationByUserAndJob(userID, jobID string) (*models.Application, error)
	ApplyForJob(application *models.Application) error
	UpdateJobApplicationsCount(jobID string) error

	PredictProfileScore(userID int, jobID int) (int, error)

	CreateJob(job models.Job) error
	GetUserByID(userID int, user *models.User) error
	GetJobDetails(jobID int) (models.Job, []models.User, error)
	GetAllUsers() ([]models.User, error)
	GetProfileByID(userID int, user *models.User) error
}

type UserUsecaseImpl struct {
	Repo postgres.UserRepository
}

func NewUserUsecase(repo postgres.UserRepository) *UserUsecaseImpl {
	return &UserUsecaseImpl{Repo: repo}
}

// CreateUser handles user creation for both Admin and Applicant
func (u *UserUsecaseImpl) CreateUser(user models.User) error {
	// Check if user with the same email already exists
	existingUser, err := u.Repo.GetUserByEmail(user.Email)
	if err == nil && existingUser.ID != 0 {
		return errors.New("user already exists with this email")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.PasswordHash = string(hashedPassword)

	// Save the user to the database
	return u.Repo.CreateUser(user)
}

// Login handles user authentication for both Admin and Applicant
func (u *UserUsecaseImpl) Login(email, password string) (string, error) {
	// Get user by email
	user, err := u.Repo.GetUserByEmail(email)
	if err != nil || user.ID == 0 {
		return "", errors.New("invalid email or password")
	}

	// Compare password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid email or password")
	}

	// Generate JWT token
	token, err := jwt.GeneratetToken(user.ID, user.UserType)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetJobs fetches the list of jobs from the repository
func (u *UserUsecaseImpl) GetJobs() ([]models.Job, error) {
	return u.Repo.FetchJobs()
}

// GetApplicationByUserAndJob checks if a user has already applied for a specific job
func (u *UserUsecaseImpl) GetApplicationByUserAndJob(userID, jobID string) (*models.Application, error) {
	application := &models.Application{}
	application, err := u.Repo.GetApplicationByUserAndJob(userID, jobID)
	if err != nil {
		return nil, fmt.Errorf("application not found")
	}
	return application, nil
}

// ApplyForJob creates a job application
func (u *UserUsecaseImpl) ApplyForJob(application *models.Application) error {
	return u.Repo.CreateApplication(application)
}

// Update the total application count for a job
func (u *UserUsecaseImpl) UpdateJobApplicationsCount(jobID string) error {
	return u.Repo.UpdateJobApplicationsCount(jobID)
}

func (u *UserUsecaseImpl) PredictProfileScore(userID int, jobID int) (int, error) {
	// Fetch job details
	job, _, err := u.Repo.GetJobDetails(jobID)
	if err != nil {
		return 0, err
	}
	log.Printf("Fetched Job: %+v", job)

	// Fetch user profile
	var user models.User
	err = u.Repo.GetProfileByID(userID, &user)
	if err != nil {
		return 0, err
	}
	log.Printf("Fetched User Profile: %+v", user.Profile)

	// calulate heuristic score
	score := services.CalculateHeuristicScore(job, user.Profile)
	log.Printf("Calculated Heuristic Score: %d", score)
	return score, nil
}

package postgres

import (
	"fmt"
	"recruitment-management/internal/models"
	"time"

	"database/sql"
)

// UserRepository defines the interface for interacting with the user database
type UserRepository interface {
	FetchJobs() ([]models.Job, error)
	CreateUser(user models.User) error
	GetUserByEmail(email string) (models.User, error)
	GetUserByID(userID int, user *models.User) error
	GetApplicationByUserAndJob(userID, jobID string, application *models.Application) error
	CreateApplication(application *models.Application) error
	UpdateJobApplicationsCount(JobID string) error

	GetJobByTitleAndUserID(title string, userID int) (models.Job, error)
	CreateJob(job models.Job) error
	GetJobDetails(jobID int) (models.Job, []models.User, error)
	GetAllUsers() ([]models.User, error)
	GetProfileByID(userID int, user *models.User) error
}

// UserRepositoryImpl implements the UserRepository interface
type UserRepositoryImpl struct {
	DB *sql.DB
}

// NewUserRepository creates a new instance of UserRepositoryImpl
func NewUserRepository(db *sql.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{DB: db}
}

// CreateUser creates a new user in the database
func (r *UserRepositoryImpl) CreateUser(user models.User) error {
	// Use a parameterized query to prevent SQL injection
	query := `INSERT INTO users (name, email, address, user_type, password_hash, profile_headline, created_at) 
	VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := r.DB.Exec(query, user.Name, user.Email, user.Address, user.UserType, user.PasswordHash, user.ProfileHeadline, time.Now())
	return err
}

// GetUserByEmail retrieves a user by email
func (r *UserRepositoryImpl) GetUserByEmail(email string) (models.User, error) {
	var user models.User
	query := `SELECT id, name, email, address, user_type, password_hash, profile_headline, created_at FROM users WHERE email = $1`
	row := r.DB.QueryRow(query, email)

	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.Address, &user.UserType, &user.PasswordHash, &user.ProfileHeadline, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, nil // No user found
		}
		return user, err // Database error
	}
	return user, nil
}

func (r *UserRepositoryImpl) FetchJobs() ([]models.Job, error) {
	var jobs []models.Job

	query := `SELECT id, title, description, posted_on, total_applications, company_name, posted_by_id FROM jobs`

	rows, err := r.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Loop through the rows and scan the data into the jobs slice
	for rows.Next() {
		var job models.Job
		if err := rows.Scan(&job.ID, &job.Title, &job.Description, &job.PostedOn, &job.TotalApplications, &job.CompanyName, &job.PostedByID); err != nil {
			return nil, err
		}
		jobs = append(jobs, job)
	}

	// Check for any errors that occured during the iteration
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return jobs, nil
}

func (r *UserRepositoryImpl) GetApplicationByUserAndJob(userID, JobID string, application *models.Application) error {
	query := "SELECT id, user_id, job_id FROM applications WHERE user_id = $1 AND job_id = $2 LIMIT 1"
	err := r.DB.QueryRow(query, userID, JobID).Scan(&application.ID, &application.UserID, &application.JobID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil // No existing application found
		}
		return fmt.Errorf("failed to get application: %v", err)
	}
	return nil
}

func (r *UserRepositoryImpl) CreateApplication(application *models.Application) error {
	query := "INSERT INTO applications (user_id, job_id) VALUES ($1, $2) RETURNING id"
	err := r.DB.QueryRow(query, application.UserID, application.JobID).Scan(&application.ID)
	if err != nil {
		return fmt.Errorf("failed to create application: %v", err)
	}
	return nil
}

func (r *UserRepositoryImpl) UpdateJobApplicationsCount(jobID string) error {
	query := "UPDATE jobs SET total_application = total_applications + 1 WHERE id = $1"
	_, err := r.DB.Exec(query, jobID)
	if err != nil {
		return fmt.Errorf("failed to update job application count: %v", err)
	}
	return nil
}

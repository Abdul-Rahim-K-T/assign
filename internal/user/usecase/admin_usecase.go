package usecase

import (
	"database/sql"
	"errors"
	"fmt"
	"recruitment-management/internal/models"
)

// CreateJob handles the job creation for admins
func (a *UserUsecaseImpl) CreateJob(job models.Job) error {
	// Check if the job title already exists in the database
	existingJob, err := a.Repo.GetJobByTitleAndUserID(job.Title, job.PostedByID)
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("failed to chek existing job: %w", err)
	}
	if existingJob.ID != 0 {
		return errors.New("job with the same title already exists for this admin")
	}

	// Create the job in the database
	err = a.Repo.CreateJob(job)
	if err != nil {
		return fmt.Errorf("failed to create job: %w", err)
	}

	return nil
}

func (u *UserUsecaseImpl) GetUserByID(userID int, user *models.User) error {
	// Query the database to fetch user details by userID
	err := u.Repo.GetUserByID(userID, user)
	if err != nil {
		return err
	}
	return nil
}

// get job details retrieves job details and list of applicants by job ID
func (u *UserUsecaseImpl) GetJobDetails(jobID int) (models.Job, []models.User, error) {
	return u.Repo.GetJobDetails(jobID)
}

func (u *UserUsecaseImpl) GetAllUsers() ([]models.User, error) {
	return u.Repo.GetAllUsers()
}

func (u *UserUsecaseImpl) GetProfileByID(userID int, user *models.User) error {
	return u.Repo.GetProfileByID(userID, user)
}

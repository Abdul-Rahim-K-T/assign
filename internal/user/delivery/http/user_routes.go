package http

import (
	"net/http"
	"recruitment-management/internal/user/repository/postgres"
	"recruitment-management/internal/user/usecase"
	"recruitment-management/pkg/database"

	"recruitment-management/pkg/jwt"

	"github.com/gorilla/mux"
)

func SetupUserRoutes(r *mux.Router) {
	// Initialize  repository, usecase, and handler
	repo := postgres.NewUserRepository(database.DB)
	usecase := usecase.NewUserUsecase(repo)
	handler := &UserHandler{Usecase: usecase}

	// Shared routes for both Admin and Applicant
	r.HandleFunc("/signup", handler.Signup).Methods("POST")
	r.HandleFunc("/login", handler.Login).Methods("POST")
	r.Handle("/jobs", jwt.AuthMiddleware(http.HandlerFunc(handler.GetJobs))).Methods("GET")

	// Routes for applicant
	r.Handle("/uploadResume", jwt.AuthMiddleware(jwt.ApplicantMiddleware(http.HandlerFunc(handler.uploadResume)))).Methods("POST")
	r.Handle("/jobs/apply", jwt.AuthMiddleware(jwt.ApplicantMiddleware(http.HandlerFunc(handler.ApplyForJob)))).Methods("GET")
	r.Handle("/predictProfileScore", jwt.AuthMiddleware(jwt.ApplicantMiddleware(http.HandlerFunc(handler.PredictProfileScore)))).Methods("GET")

	// Routes for admins - use AdminMiddleware for adminspecific rout
	r.Handle("/admin/job", jwt.AuthMiddleware(jwt.AdminMiddleware(http.HandlerFunc(handler.CreateJob)))).Methods("POST")
	r.Handle("/admin/job/{job_id}", jwt.AuthMiddleware(jwt.AdminMiddleware(http.HandlerFunc(handler.GetJobDetails)))).Methods("GET")
	r.Handle("/admin/applicants", jwt.AuthMiddleware(jwt.AdminMiddleware(http.HandlerFunc(handler.GetAllUsers)))).Methods("GET")
	r.Handle("/admin/applicant/{applicant_id}", jwt.AuthMiddleware(jwt.AdminMiddleware(http.HandlerFunc(handler.GetApplicantByID)))).Methods("GET")
}

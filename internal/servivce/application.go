package service

import (
	"log"

	"github.com/manan-sharma-5/seekajob-backend/internal/model"
	"github.com/manan-sharma-5/seekajob-backend/internal/repository"
)


func ApplyJobService(req *model.ApplyJobRequest) error {
	// Apply for the job
	log.Println("Applying for job with ID:", req.JobID, "by user with ID:", req.UserID)
	err := repository.ApplyJob(req.JobID, req.UserID)
	// Return success or error
	if err != nil {
		return err
	}
	// return nil
	return nil
}
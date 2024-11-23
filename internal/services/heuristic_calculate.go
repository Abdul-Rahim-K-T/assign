package services

import (
	"log"
	"recruitment-management/internal/models"
	"strconv"
	"strings"
)

func CalculateHeuristicScore(job models.Job, profile models.Profile) int {
	score := 0
	// Normalize skills by trimming spaces and converting to lowercas
	normalize := func(s string) []string {
		parts := strings.Split(s, ",")
		for i := range parts {
			parts[i] = strings.TrimSpace(strings.ToLower(parts[i]))
		}
		return parts
	}
	// Matching skills
	jobSkills := normalize(job.Description)
	profileSkills := normalize(profile.Skills)
	log.Printf("Normalized Job Skills: %v", jobSkills)
	log.Printf("Normalized Profile Skills: %v", profileSkills)

	matchedSkills := make([]string, 0)
	for _, jobSkill := range jobSkills {
		for _, profileSkill := range profileSkills {
			if jobSkill == profileSkill {
				score += 10 // Add 10 points for each matching skills
				matchedSkills = append(matchedSkills, jobSkill)
			}
		}
	}

	log.Printf("Matched Skills: %v", matchedSkills)

	// Matching education
	jobTitle := strings.ToLower(strings.TrimSpace(job.Title))
	for _, edu := range profile.Education {
		edu := strconv.Itoa(int(edu))
		if strings.Contains((strings.ToLower(edu)), jobTitle) {
			score += 20 // Add 20 points if education matches the job title
			break
		}
	}

	// Matching experience
	for _, exp := range profile.Experience {
		exp := strconv.Itoa(int(exp))
		if strings.Contains(strings.ToLower(exp), jobTitle) {
			score += 30 // Add 30 points if the experience matches the job title
			break
		}
	}

	log.Printf("Final Score: %d", score)

	return score
}

package model

import (
	"time"

	"github.com/google/uuid"
)

type ScanRequest struct {
	RepositoryURL string   `json:"repository_url" binding:"required"`
	Branch        string   `json:"branch" binding:"required"`
	Languages     []string `json:"languages"`
}

type ScanReport struct {
	ID             uuid.UUID  `json:"id"`
	UserID         uuid.UUID  `json:"user_id"`
	RepositoryURL  string     `json:"repository_url"`
	Branch         string     `json:"branch"`
	Status         string     `json:"status"`
	TotalFiles     int        `json:"total_files"`
	TotalIssues    int        `json:"total_issues"`
	CriticalIssues int        `json:"critical_issues"`
	HighIssues     int        `json:"high_issues"`
	MediumIssues   int        `json:"medium_issues"`
	LowIssues      int        `json:"low_issues"`
	CreatedAt      time.Time  `json:"created_at"`
	CompletedAt    *time.Time `json:"completed_at,omitempty"`
}

type Issue struct {
	ID          uuid.UUID `json:"id"`
	ReportID    uuid.UUID `json:"report_id"`
	FilePath    string    `json:"file_path"`
	Line        int       `json:"line"`
	Column      int       `json:"column"`
	Severity    string    `json:"severity"`
	Category    string    `json:"category"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Suggestion  string    `json:"suggestion"`
	CreatedAt   time.Time `json:"created_at"`
}
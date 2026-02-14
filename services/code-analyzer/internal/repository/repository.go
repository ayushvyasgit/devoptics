package repository

import (
	"database/sql"
	"time"

	"github.com/ayushvyasgit/devoptics/services/code-analyzer/internal/model"
	"github.com/google/uuid"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) CreateReport(report *model.ScanReport) error {
	query := `
		INSERT INTO scan_reports (id, user_id, repository_url, branch, status, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, report.ID, report.UserID, report.RepositoryURL,
		report.Branch, report.Status, report.CreatedAt)
	return err
}

func (r *Repository) UpdateReport(report *model.ScanReport) error {
	now := time.Now()
	query := `
		UPDATE scan_reports 
		SET status = $1, total_files = $2, total_issues = $3, 
		    critical_issues = $4, high_issues = $5, medium_issues = $6, 
		    low_issues = $7, completed_at = $8
		WHERE id = $9
	`
	_, err := r.db.Exec(query, report.Status, report.TotalFiles, report.TotalIssues,
		report.CriticalIssues, report.HighIssues, report.MediumIssues,
		report.LowIssues, now, report.ID)
	return err
}

func (r *Repository) GetReport(id uuid.UUID) (*model.ScanReport, error) {
	report := &model.ScanReport{}
	query := `
		SELECT id, user_id, repository_url, branch, status, total_files, total_issues,
		       critical_issues, high_issues, medium_issues, low_issues, created_at, completed_at
		FROM scan_reports WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&report.ID, &report.UserID, &report.RepositoryURL, &report.Branch, &report.Status,
		&report.TotalFiles, &report.TotalIssues, &report.CriticalIssues, &report.HighIssues,
		&report.MediumIssues, &report.LowIssues, &report.CreatedAt, &report.CompletedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return report, err
}

func (r *Repository) GetReportsByUser(userID uuid.UUID) ([]model.ScanReport, error) {
	query := `
		SELECT id, user_id, repository_url, branch, status, total_files, total_issues,
		       critical_issues, high_issues, medium_issues, low_issues, created_at, completed_at
		FROM scan_reports 
		WHERE user_id = $1 
		ORDER BY created_at DESC
		LIMIT 50
	`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var reports []model.ScanReport
	for rows.Next() {
		var report model.ScanReport
		if err := rows.Scan(
			&report.ID, &report.UserID, &report.RepositoryURL, &report.Branch, &report.Status,
			&report.TotalFiles, &report.TotalIssues, &report.CriticalIssues, &report.HighIssues,
			&report.MediumIssues, &report.LowIssues, &report.CreatedAt, &report.CompletedAt,
		); err != nil {
			continue
		}
		reports = append(reports, report)
	}

	return reports, nil
}

func (r *Repository) SaveIssue(issue *model.Issue) error {
	query := `
		INSERT INTO issues (id, report_id, file_path, line, column, severity, category, 
		                   title, description, suggestion, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	_, err := r.db.Exec(query, issue.ID, issue.ReportID, issue.FilePath, issue.Line,
		issue.Column, issue.Severity, issue.Category, issue.Title, issue.Description,
		issue.Suggestion, issue.CreatedAt)
	return err
}

func (r *Repository) GetIssuesByReport(reportID uuid.UUID) ([]model.Issue, error) {
	query := `
		SELECT id, report_id, file_path, line, column, severity, category, 
		       title, description, suggestion, created_at
		FROM issues 
		WHERE report_id = $1 
		ORDER BY severity DESC, file_path, line
	`
	rows, err := r.db.Query(query, reportID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var issues []model.Issue
	for rows.Next() {
		var issue model.Issue
		if err := rows.Scan(
			&issue.ID, &issue.ReportID, &issue.FilePath, &issue.Line, &issue.Column,
			&issue.Severity, &issue.Category, &issue.Title, &issue.Description,
			&issue.Suggestion, &issue.CreatedAt,
		); err != nil {
			continue
		}
		issues = append(issues, issue)
	}

	return issues, nil
}
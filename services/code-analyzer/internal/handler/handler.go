package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/ayushvyasgit/devoptics/services/code-analyzer/internal/analyzer"
	"github.com/ayushvyasgit/devoptics/services/code-analyzer/internal/model"
	"github.com/ayushvyasgit/devoptics/services/code-analyzer/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Handler struct {
	repo     *repository.Repository
	analyzer *analyzer.CodeAnalyzer
}

func NewHandler(db *sql.DB, workDir string) *Handler {
	return &Handler{
		repo:     repository.NewRepository(db),
		analyzer: analyzer.NewCodeAnalyzer(workDir),
	}
}

func (h *Handler) ScanRepository(c *gin.Context) {
	var req model.ScanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from header (set by API Gateway)
	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Create scan report
	report := &model.ScanReport{
		ID:            uuid.New(),
		UserID:        userID,
		RepositoryURL: req.RepositoryURL,
		Branch:        req.Branch,
		Status:        "in_progress",
		CreatedAt:     time.Now(),
	}

	// Save to database
	if err := h.repo.CreateReport(report); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create report"})
		return
	}

	// Process asynchronously
	go h.processRepository(report)

	c.JSON(http.StatusAccepted, gin.H{
		"message":   "Scan started",
		"report_id": report.ID,
	})
}

func (h *Handler) processRepository(report *model.ScanReport) {
	// For demo: create mock repository instead of cloning
	// In production, use git clone
	repoPath, err := h.analyzer.CreateMockRepository()
	if err != nil {
		report.Status = "failed"
		h.repo.UpdateReport(report)
		return
	}
	defer h.analyzer.Cleanup(repoPath)

	// Analyze code
	issues, fileCount, err := h.analyzer.AnalyzeDirectory(repoPath)
	if err != nil {
		report.Status = "failed"
		h.repo.UpdateReport(report)
		return
	}

	// Count issues by severity
	report.TotalFiles = fileCount
	report.TotalIssues = len(issues)
	for _, issue := range issues {
		issue.ReportID = report.ID
		issue.CreatedAt = time.Now()
		
		h.repo.SaveIssue(&issue)

		switch issue.Severity {
		case "critical":
			report.CriticalIssues++
		case "high":
			report.HighIssues++
		case "medium":
			report.MediumIssues++
		case "low":
			report.LowIssues++
		}
	}

	// Update report
	report.Status = "completed"
	h.repo.UpdateReport(report)
}

func (h *Handler) GetReport(c *gin.Context) {
	reportID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	report, err := h.repo.GetReport(reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get report"})
		return
	}

	if report == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Report not found"})
		return
	}

	c.JSON(http.StatusOK, report)
}

func (h *Handler) GetReports(c *gin.Context) {
	userIDStr := c.GetHeader("X-User-ID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	reports, err := h.repo.GetReportsByUser(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get reports"})
		return
	}

	if reports == nil {
		reports = []model.ScanReport{}
	}

	c.JSON(http.StatusOK, reports)
}

func (h *Handler) GetIssues(c *gin.Context) {
	reportID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid report ID"})
		return
	}

	issues, err := h.repo.GetIssuesByReport(reportID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get issues"})
		return
	}

	if issues == nil {
		issues = []model.Issue{}
	}

	c.JSON(http.StatusOK, issues)
}

func (h *Handler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "healthy",
		"service": "code-analyzer",
	})
}
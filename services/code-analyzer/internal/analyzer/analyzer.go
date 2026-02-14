package analyzer

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/ayushvyasgit/devoptics/services/code-analyzer/internal/model"
	"github.com/google/uuid"
)

type CodeAnalyzer struct {
	workDir string
}

func NewCodeAnalyzer(workDir string) *CodeAnalyzer {
	return &CodeAnalyzer{workDir: workDir}
}

// AnalyzeDirectory scans all files in a directory
func (a *CodeAnalyzer) AnalyzeDirectory(path string) ([]model.Issue, int, error) {
	var issues []model.Issue
	fileCount := 0

	err := filepath.WalkDir(path, func(filePath string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Skip directories and hidden files
		if d.IsDir() || strings.HasPrefix(d.Name(), ".") {
			return nil
		}

		// Skip non-code files
		if !isCodeFile(filePath) {
			return nil
		}

		fileCount++

		// Analyze individual file
		fileIssues := a.analyzeFile(filePath)
		issues = append(issues, fileIssues...)

		return nil
	})

	if err != nil {
		return nil, 0, err
	}

	return issues, fileCount, nil
}

// analyzeFile performs basic static analysis on a single file
func (a *CodeAnalyzer) analyzeFile(filePath string) []model.Issue {
	var issues []model.Issue
	relativePath := filepath.Base(filePath)

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return issues
	}

	lines := strings.Split(string(content), "\n")

	// Basic pattern matching
	for lineNum, line := range lines {
		line = strings.TrimSpace(line)

		// Check for TODO comments
		if strings.Contains(strings.ToUpper(line), "TODO") {
			issues = append(issues, model.Issue{
				ID:          uuid.New(),
				FilePath:    relativePath,
				Line:        lineNum + 1,
				Column:      1,
				Severity:    "low",
				Category:    "code-quality",
				Title:       "TODO comment found",
				Description: "TODO comment should be resolved",
				Suggestion:  "Complete the TODO or create a ticket",
			})
		}

		// Check for console.log (JavaScript)
		if strings.Contains(line, "console.log") {
			issues = append(issues, model.Issue{
				ID:          uuid.New(),
				FilePath:    relativePath,
				Line:        lineNum + 1,
				Column:      strings.Index(line, "console.log"),
				Severity:    "medium",
				Category:    "code-quality",
				Title:       "console.log statement",
				Description: "Remove console.log before production",
				Suggestion:  "Use proper logging framework",
			})
		}

		// Check for print statements (Python)
		if strings.Contains(line, "print(") && strings.HasSuffix(filePath, ".py") {
			issues = append(issues, model.Issue{
				ID:          uuid.New(),
				FilePath:    relativePath,
				Line:        lineNum + 1,
				Column:      strings.Index(line, "print("),
				Severity:    "low",
				Category:    "code-quality",
				Title:       "print statement found",
				Description: "Use proper logging instead of print",
				Suggestion:  "Replace with logging.info() or logging.debug()",
			})
		}

		// Check for hardcoded passwords
		if strings.Contains(strings.ToLower(line), "password") && 
		   (strings.Contains(line, "=") || strings.Contains(line, ":")) &&
		   (strings.Contains(line, "\"") || strings.Contains(line, "'")) {
			issues = append(issues, model.Issue{
				ID:          uuid.New(),
				FilePath:    relativePath,
				Line:        lineNum + 1,
				Column:      1,
				Severity:    "critical",
				Category:    "security",
				Title:       "Potential hardcoded password",
				Description: "Hardcoded credentials detected",
				Suggestion:  "Use environment variables for credentials",
			})
		}

		// Check for eval() usage
		if strings.Contains(line, "eval(") {
			issues = append(issues, model.Issue{
				ID:          uuid.New(),
				FilePath:    relativePath,
				Line:        lineNum + 1,
				Column:      strings.Index(line, "eval("),
				Severity:    "high",
				Category:    "security",
				Title:       "Dangerous eval() usage",
				Description: "eval() can execute arbitrary code",
				Suggestion:  "Avoid using eval(), find safer alternative",
			})
		}
	}

	return issues
}

// CreateMockRepository creates a mock repo for testing (no git required)
func (a *CodeAnalyzer) CreateMockRepository() (string, error) {
	// Create temp directory
	mockPath := filepath.Join(a.workDir, fmt.Sprintf("mock-%s", uuid.New().String()[:8]))
	if err := os.MkdirAll(mockPath, 0755); err != nil {
		return "", err
	}

	// Create sample files
	files := map[string]string{
		"main.js": `
console.log("Starting application");
const password = "hardcoded123"; // TODO: Move to env
function processData() {
    console.log("Processing");
}
`,
		"utils.py": `
print("Utility module loaded")
def calculate():
    password = "admin123"
    return 42
`,
		"config.go": `
package main
// TODO: Add configuration validation
func main() {
    // Clean code here
}
`,
	}

	for filename, content := range files {
		filePath := filepath.Join(mockPath, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return "", err
		}
	}

	return mockPath, nil
}

// Cleanup removes cloned repository
func (a *CodeAnalyzer) Cleanup(path string) error {
	return os.RemoveAll(path)
}

// isCodeFile checks if a file is a code file
func isCodeFile(path string) bool {
	codeExtensions := map[string]bool{
		".go":   true,
		".js":   true,
		".jsx":  true,
		".ts":   true,
		".tsx":  true,
		".py":   true,
		".java": true,
		".rs":   true,
		".cpp":  true,
		".c":    true,
		".rb":   true,
		".php":  true,
		".cs":   true,
	}

	ext := filepath.Ext(path)
	return codeExtensions[ext]
}
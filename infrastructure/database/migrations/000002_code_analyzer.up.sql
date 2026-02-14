-- Scan Reports table
CREATE TABLE scan_reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id),
    repository_url VARCHAR(500) NOT NULL,
    branch VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    total_files INT DEFAULT 0,
    total_issues INT DEFAULT 0,
    critical_issues INT DEFAULT 0,
    high_issues INT DEFAULT 0,
    medium_issues INT DEFAULT 0,
    low_issues INT DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    completed_at TIMESTAMP WITH TIME ZONE
);

CREATE INDEX idx_scan_reports_user_id ON scan_reports(user_id);
CREATE INDEX idx_scan_reports_status ON scan_reports(status);

-- Issues table
CREATE TABLE issues (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NOT NULL REFERENCES scan_reports(id) ON DELETE CASCADE,
    file_path VARCHAR(500) NOT NULL,
    line INT NOT NULL,
    column_number INT NOT NULL,
    severity VARCHAR(50) NOT NULL,
    category VARCHAR(100) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    suggestion TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

CREATE INDEX idx_issues_report_id ON issues(report_id);
CREATE INDEX idx_issues_severity ON issues(severity);
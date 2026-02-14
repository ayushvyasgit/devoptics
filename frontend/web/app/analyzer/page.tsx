"use client";

import { useState, useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/authStore";
import api from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Card, CardHeader, CardTitle, CardDescription, CardContent } from "@/components/ui/card";

interface Report {
  id: string;
  repository_url: string;
  branch: string;
  status: string;
  total_issues: number;
  critical_issues: number;
  high_issues: number;
  medium_issues: number;
  low_issues: number;
  created_at: string;
}

export default function AnalyzerPage() {
  const router = useRouter();
  const { token } = useAuthStore();
  const [repositoryUrl, setRepositoryUrl] = useState("");
  const [branch, setBranch] = useState("main");
  const [loading, setLoading] = useState(false);
  const [reports, setReports] = useState<Report[]>([]);

  useEffect(() => {
    if (!token) {
      router.push("/login");
      return;
    }
    loadReports();
  }, [token, router]);

  const loadReports = async () => {
    try {
      const response = await api.get("/analyzer/reports");
      setReports(response.data || []);
    } catch (error) {
      console.error("Failed to load reports:", error);
    }
  };

  const handleScan = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);

    try {
      const response = await api.post("/analyzer/scan", {
        repository_url: repositoryUrl,
        branch,
      });

      alert(`✅ ${response.data.message}\nReport ID: ${response.data.report_id}`);
      setRepositoryUrl("");
      setBranch("main");
      
      // Reload reports after 2 seconds
      setTimeout(loadReports, 2000);
    } catch (error: any) {
      alert(`❌ Scan failed: ${error.response?.data?.error || error.message}`);
    } finally {
      setLoading(false);
    }
  };

  const getSeverityColor = (severity: string) => {
    const colors: Record<string, string> = {
      critical: "bg-red-100 text-red-800",
      high: "bg-orange-100 text-orange-800",
      medium: "bg-yellow-100 text-yellow-800",
      low: "bg-blue-100 text-blue-800",
    };
    return colors[severity] || "bg-gray-100 text-gray-800";
  };

  if (!token) {
    return <div className="flex h-screen items-center justify-center">Loading...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="mx-auto max-w-6xl space-y-8">
        <div>
          <h1 className="text-3xl font-bold text-gray-900">Code Analyzer</h1>
          <p className="text-gray-600 mt-2">
            Scan your repositories for code quality issues
          </p>
        </div>

        {/* Scan Form */}
        <Card>
          <CardHeader>
            <CardTitle>New Scan</CardTitle>
            <CardDescription>Enter repository details to start scanning</CardDescription>
          </CardHeader>
          <CardContent>
            <form onSubmit={handleScan} className="space-y-4">
              <div>
                <label className="block text-sm font-medium mb-2">Repository URL</label>
                <Input
                  type="url"
                  placeholder="https://github.com/username/repository"
                  value={repositoryUrl}
                  onChange={(e) => setRepositoryUrl(e.target.value)}
                  required
                />
              </div>

              <div>
                <label className="block text-sm font-medium mb-2">Branch</label>
                <Input
                  type="text"
                  placeholder="main"
                  value={branch}
                  onChange={(e) => setBranch(e.target.value)}
                  required
                />
              </div>

              <Button type="submit" loading={loading}>
                Start Scan
              </Button>
            </form>
          </CardContent>
        </Card>

        {/* Reports List */}
        <Card>
          <CardHeader>
            <CardTitle>Scan Reports</CardTitle>
            <CardDescription>Your repository scan history</CardDescription>
          </CardHeader>
          <CardContent>
            {reports.length === 0 ? (
              <div className="text-center py-8 text-gray-500">
                No scans yet. Start your first scan above.
              </div>
            ) : (
              <div className="space-y-4">
                {reports.map((report) => (
                  <div key={report.id} className="border rounded-lg p-4 hover:bg-gray-50">
                    <div className="flex items-start justify-between">
                      <div className="flex-1">
                        <div className="font-medium">{report.repository_url}</div>
                        <div className="text-sm text-gray-500 mt-1">
                          Branch: {report.branch} • Status: {report.status}
                        </div>
                        
                        {report.status === "completed" && (
                          <div className="mt-3 flex gap-2">
                            {report.critical_issues > 0 && (
                              <span className={`px-2 py-1 text-xs rounded ${getSeverityColor("critical")}`}>
                                {report.critical_issues} Critical
                              </span>
                            )}
                            {report.high_issues > 0 && (
                              <span className={`px-2 py-1 text-xs rounded ${getSeverityColor("high")}`}>
                                {report.high_issues} High
                              </span>
                            )}
                            {report.medium_issues > 0 && (
                              <span className={`px-2 py-1 text-xs rounded ${getSeverityColor("medium")}`}>
                                {report.medium_issues} Medium
                              </span>
                            )}
                            {report.low_issues > 0 && (
                              <span className={`px-2 py-1 text-xs rounded ${getSeverityColor("low")}`}>
                                {report.low_issues} Low
                              </span>
                            )}
                          </div>
                        )}
                      </div>

                      {report.status === "completed" && (
                        <div className="text-right">
                          <div className="text-2xl font-bold">{report.total_issues}</div>
                          <div className="text-xs text-gray-500">Issues</div>
                        </div>
                      )}

                      {report.status === "in_progress" && (
                        <div className="text-blue-600">⏳ Scanning...</div>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            )}
          </CardContent>
        </Card>
      </div>
    </div>
  );
}
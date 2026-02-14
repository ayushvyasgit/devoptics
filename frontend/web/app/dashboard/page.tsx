"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export default function DashboardPage() {
  const router = useRouter();
  const { user, token, logout } = useAuthStore();

  useEffect(() => {
    // Check if user is logged in
    if (!token) {
      router.push("/login");
    }
  }, [token, router]);

  const handleLogout = () => {
    logout();
    router.push("/login");
  };

  if (!user) {
    return <div className="flex h-screen items-center justify-center">Loading...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="mx-auto max-w-4xl">
        <div className="mb-8 flex items-center justify-between">
          <h1 className="text-3xl font-bold text-gray-900">Dashboard</h1>
          <Button variant="outline" onClick={handleLogout}>
            Logout
          </Button>
        </div>

        <Card>
          <CardHeader>
            <CardTitle>Welcome, {user.first_name}! 🎉</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-2">
              <p><strong>Email:</strong> {user.email}</p>
              <p><strong>Name:</strong> {user.first_name} {user.last_name}</p>
              <p><strong>Role:</strong> {user.role}</p>
              <p className="text-sm text-green-600">✅ Auth service is working!</p>
            </div>
          </CardContent>
        </Card>

        <div className="mt-8">
          <Card>
            <CardHeader>
              <CardTitle>Next Steps</CardTitle>
            </CardHeader>
            <CardContent>
              <ul className="list-disc list-inside space-y-2 text-sm text-gray-600">
                <li>✅ Infrastructure running (Docker)</li>
                <li>✅ Database migrations complete</li>
                <li>✅ Auth Service working (Port 8091)</li>
                <li>✅ Frontend working (Port 3000)</li>
                <li>⏳ Next: Build API Gateway</li>
                <li>⏳ Next: Build Code Analyzer Service</li>
                <li>⏳ Next: Build AI Engine</li>
              </ul>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
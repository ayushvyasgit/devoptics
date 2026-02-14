import Link from "next/link";
import { Button } from "@/components/ui/button";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export default function HomePage() {
  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-blue-50 to-indigo-100">
      <Card className="w-full max-w-2xl">
        <CardHeader>
          <CardTitle className="text-center text-4xl text-blue-600">
            🚀 DevOptics
          </CardTitle>
          <p className="text-center text-xl text-gray-600 mt-2">
            AI-Powered DevOps Intelligence Platform
          </p>
        </CardHeader>
        <CardContent>
          <div className="space-y-6">
            <div className="text-center space-y-4">
              <p className="text-gray-700">
                Welcome to DevOptics! This platform helps you analyze code quality,
                detect security vulnerabilities, and monitor system health with AI.
              </p>
              
              <div className="flex gap-4 justify-center">
                <Link href="/register">
                  <Button className="px-8">Get Started</Button>
                </Link>
                <Link href="/login">
                  <Button variant="outline" className="px-8">Sign In</Button>
                </Link>
              </div>
            </div>

            <div className="border-t pt-6 mt-6">
              <h3 className="font-semibold text-lg mb-3">✨ Features:</h3>
              <ul className="space-y-2 text-sm text-gray-600">
                <li>✅ AI-powered code analysis with CodeBERT</li>
                <li>✅ Real-time anomaly detection</li>
                <li>✅ Security vulnerability scanning</li>
                <li>✅ Incident summarization with LLMs</li>
                <li>✅ System metrics monitoring</li>
              </ul>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
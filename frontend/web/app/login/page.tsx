"use client";
import axios from "axios";
import { useState } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { authAPI } from "@/lib/api";
import { useAuthStore } from "@/store/authStore";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";

export default function LoginPage() {
  const router = useRouter();
  const { setAuth } = useAuthStore();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);

    try {
      const response = await authAPI.login({ email, password });
      setAuth(response.data.user, response.data.token);
      router.push("/");
    } catch (err) {
      if (axios.isAxiosError(err)) {
        setError(err.response?.data?.error || "Login failed");
      } else {
        setError("Something went wrong");
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="relative flex min-h-screen items-center justify-center overflow-hidden bg-gradient-to-br from-gray-50 via-blue-50 to-indigo-100">
      {/* Animated floating blobs (like Shopify / modern SaaS aesthetic) */}
      <div className="absolute inset-0 z-0">
        <div className="absolute -left-20 top-20 h-72 w-72 animate-blob rounded-full bg-blue-300 opacity-30 mix-blend-multiply blur-3xl"></div>
        <div className="absolute right-10 top-40 h-80 w-80 animate-blob animation-delay-2000 rounded-full bg-purple-300 opacity-30 mix-blend-multiply blur-3xl"></div>
        <div className="absolute -bottom-20 left-1/3 h-96 w-96 animate-blob animation-delay-4000 rounded-full bg-pink-300 opacity-30 mix-blend-multiply blur-3xl"></div>
        <div className="absolute bottom-10 right-20 h-64 w-64 animate-blob animation-delay-6000 rounded-full bg-cyan-300 opacity-25 mix-blend-multiply blur-3xl"></div>
      </div>

      {/* Main Card – with entrance animation */}
      <Card
        className="relative z-10 w-full max-w-md transform border border-gray-200/60 bg-white/80 shadow-2xl backdrop-blur-sm transition-all duration-700 animate-in fade-in slide-in-from-bottom-12 zoom-in-95"
      >
        <CardHeader className="space-y-1 text-center">
          <CardTitle className="bg-gradient-to-r from-blue-600 to-indigo-600 bg-clip-text text-4xl font-extrabold text-transparent">
            DevOptics
          </CardTitle>
          <CardDescription className="text-base text-gray-600">
            Sign in to continue your journey
          </CardDescription>
        </CardHeader>

        <CardContent>
          <form onSubmit={handleSubmit} className="space-y-6">
            {error && (
              <div className="animate-in fade-in slide-in-from-top-5 rounded-lg bg-red-50/90 p-4 text-sm text-red-800 shadow-sm backdrop-blur-sm">
                {error}
              </div>
            )}

            {/* Email – floating label style */}
            <div className="group relative">
              <label
                htmlFor="email"
                className={`pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 origin-left text-sm text-gray-500 transition-all duration-300 group-focus-within:-translate-y-9 group-focus-within:scale-90 group-focus-within:font-medium group-focus-within:text-blue-600 ${
                  email ? "-translate-y-9 scale-90 font-medium text-blue-600" : ""
                }`}
              >
                Email address
              </label>
              <Input
                id="email"
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                required
                className="peer h-12 border-gray-300 bg-white/70 pl-4 pt-4 shadow-sm transition-all focus:border-blue-500 focus:ring-blue-500/30"
              />
            </div>

            {/* Password – floating label */}
            <div className="group relative">
              <label
                htmlFor="password"
                className={`pointer-events-none absolute left-3 top-1/2 -translate-y-1/2 origin-left text-sm text-gray-500 transition-all duration-300 group-focus-within:-translate-y-9 group-focus-within:scale-90 group-focus-within:font-medium group-focus-within:text-blue-600 ${
                  password ? "-translate-y-9 scale-90 font-medium text-blue-600" : ""
                }`}
              >
                Password
              </label>
              <Input
                id="password"
                type="password"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                required
                className="peer h-12 border-gray-300 bg-white/70 pl-4 pt-4 shadow-sm transition-all focus:border-blue-500 focus:ring-blue-500/30"
              />
            </div>

            <Button
              type="submit"
              className="group h-12 w-full overflow-hidden bg-gradient-to-r from-blue-600 to-indigo-600 text-lg font-semibold shadow-lg transition-all duration-300 hover:scale-[1.02] hover:shadow-xl hover:from-blue-700 hover:to-indigo-700 disabled:opacity-70"
              disabled={loading}
            >
              <span className="relative z-10">
                {loading ? "Signing in..." : "Sign in"}
              </span>
              <span className="absolute inset-0 translate-x-full bg-white/20 transition-transform duration-500 group-hover:translate-x-0"></span>
            </Button>

            <p className="text-center text-sm text-gray-600">
              Don’t have an account?{" "}
              <Link
                href="/register"
                className="font-medium text-blue-600 transition-colors hover:text-blue-800"
              >
                Sign up
              </Link>
            </p>
          </form>
        </CardContent>
      </Card>
    </div>
  );
}
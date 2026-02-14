from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware
from app.api import code_analysis, anomaly
from app.core.config import settings

app = FastAPI(
    title="DevOptics AI Engine",
    description="AI-powered code analysis and anomaly detection",
    version="1.0.0"
)

# CORS
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000", "http://localhost:8080"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Include routers
app.include_router(code_analysis.router)
app.include_router(anomaly.router)

@app.get("/health")
async def health_check():
    return {
        "status": "healthy",
        "service": "ai-engine",
        "version": "1.0.0"
    }

@app.on_event("startup")
async def startup_event():
    print("🚀 DevOptics AI Engine starting...")
    print(f"📊 Environment: {settings.environment}")
    print(f"🔧 Bug threshold: {settings.bug_threshold}")
    print(f"⚡ Anomaly threshold: {settings.anomaly_threshold}")
    print("✅ AI Engine ready!")

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(
        "app.main:app",
        host="0.0.0.0",
        port=settings.port,
        reload=settings.environment == "development"
    )
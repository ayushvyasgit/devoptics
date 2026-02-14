import os
from typing import Optional
from pydantic_settings import BaseSettings

class Settings(BaseSettings):
    # Server
    port: int = int(os.getenv("AI_ENGINE_PORT", "8083"))
    environment: str = os.getenv("ENVIRONMENT", "development")
    
    # Model settings
    bug_threshold: float = float(os.getenv("BUG_THRESHOLD", "0.7"))
    anomaly_threshold: float = float(os.getenv("ANOMALY_THRESHOLD", "0.8"))
    
    # Optional API keys
    openai_api_key: Optional[str] = os.getenv("OPENAI_API_KEY")
    anthropic_api_key: Optional[str] = os.getenv("ANTHROPIC_API_KEY")
    
    # Model paths
    codebert_model: str = "microsoft/codebert-base"
    
    class Config:
        case_sensitive = False

settings = Settings()
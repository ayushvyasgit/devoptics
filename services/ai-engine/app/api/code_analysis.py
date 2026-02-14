from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import List, Optional, Dict
from app.models.code_analyzer import get_code_analyzer
from app.models.bug_detector import get_bug_detector
from app.models.security_scanner import get_security_scanner

router = APIRouter(prefix="/code", tags=["code-analysis"])

class CodeAnalysisRequest(BaseModel):
    code: str
    language: str = "python"
    check_security: bool = True
    check_bugs: bool = True

class CodeAnalysisResult(BaseModel):
    bug_probability: Optional[float] = None
    security_issues: List[Dict] = []
    complexity_score: float = 0.0
    recommendations: List[str] = []

class CodeSimilarityRequest(BaseModel):
    code1: str
    code2: str

class CodeSimilarityResult(BaseModel):
    similarity: float
    is_clone: bool
    message: str

@router.post("/analyze", response_model=CodeAnalysisResult)
async def analyze_code(request: CodeAnalysisRequest):
    """Analyze code for bugs and security issues"""
    try:
        result = CodeAnalysisResult()
        recommendations = []
        
        # Bug detection
        if request.check_bugs:
            bug_detector = get_bug_detector()
            bug_prob = bug_detector.predict_bug_probability(request.code)
            result.bug_probability = bug_prob
            
            if bug_prob > 0.7:
                recommendations.append("High bug probability detected. Review code carefully.")
        
        # Security scanning
        if request.check_security:
            scanner = get_security_scanner()
            vulnerabilities = scanner.scan_code(request.code, request.language)
            result.security_issues = vulnerabilities
            
            if vulnerabilities:
                recommendations.append(f"Found {len(vulnerabilities)} security issues. Address them immediately.")
        
        result.recommendations = recommendations
        
        return result
        
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))

@router.post("/similarity", response_model=CodeSimilarityResult)
async def check_similarity(request: CodeSimilarityRequest):
    """Check similarity between two code snippets"""
    try:
        analyzer = get_code_analyzer()
        similarity = analyzer.calculate_similarity(request.code1, request.code2)
        is_clone = similarity > 0.9
        
        return CodeSimilarityResult(
            similarity=similarity,
            is_clone=is_clone,
            message=f"Similarity: {similarity:.2%}. {'Code clone detected!' if is_clone else 'Codes are different.'}"
        )
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
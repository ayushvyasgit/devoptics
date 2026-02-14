from fastapi import APIRouter, HTTPException
from pydantic import BaseModel
from typing import List
from app.models.anomaly_detector import get_anomaly_detector

router = APIRouter(prefix="/anomaly", tags=["anomaly-detection"])

class MetricPoint(BaseModel):
    timestamp: str
    value: float
    service_name: str
    metric_type: str

class AnomalyResult(BaseModel):
    timestamp: str
    value: float
    is_anomaly: bool
    anomaly_score: float
    severity: str

@router.post("/detect", response_model=List[AnomalyResult])
async def detect_anomalies(metrics: List[MetricPoint]):
    """Detect anomalies in metrics"""
    try:
        results = []
        
        for metric in metrics:
            detector = get_anomaly_detector(metric.service_name, metric.metric_type)
            result = detector.detect_anomaly(
                metric.value,
                metric.service_name,
                metric.metric_type
            )
            
            results.append(AnomalyResult(
                timestamp=metric.timestamp,
                value=metric.value,
                is_anomaly=result['is_anomaly'],
                anomaly_score=result['anomaly_score'],
                severity=result['severity']
            ))
        
        return results
    
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))
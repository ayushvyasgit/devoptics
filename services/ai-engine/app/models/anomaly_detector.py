import numpy as np
from pyod.models.iforest import IForest
from typing import List, Dict
from collections import deque

class AnomalyDetector:
    def __init__(self, contamination: float = 0.1):
        self.model = IForest(
            n_estimators=100,
            contamination=contamination,
            random_state=42
        )
        self.history = deque(maxlen=100)  # Keep last 100 data points
        self.is_fitted = False
    
    def detect_anomaly(self, value: float, service_name: str, metric_type: str) -> Dict:
        """Detect if a metric value is anomalous"""
        self.history.append(value)
        
        if len(self.history) < 10:
            # Not enough data yet
            return {
                'is_anomaly': False,
                'anomaly_score': 0.0,
                'severity': 'low',
            }
        
        # Extract features
        features = self._extract_features(value)
        
        # Fit model if needed
        if not self.is_fitted and len(self.history) >= 20:
            X = np.array(list(self.history)).reshape(-1, 1)
            self.model.fit(X)
            self.is_fitted = True
        
        if not self.is_fitted:
            # Use statistical method as fallback
            return self._statistical_detect(value)
        
        # Predict
        score = self.model.decision_function(features.reshape(1, -1))[0]
        is_anomaly = self.model.predict(features.reshape(1, -1))[0] == 1
        
        # Normalize score to 0-1
        anomaly_score = min(max(abs(score), 0), 1)
        
        severity = self._get_severity(anomaly_score)
        
        return {
            'is_anomaly': bool(is_anomaly),
            'anomaly_score': float(anomaly_score),
            'severity': severity,
        }
    
    def _extract_features(self, value: float) -> np.ndarray:
        """Extract features for anomaly detection"""
        history_array = np.array(list(self.history))
        
        features = [
            value,
            np.mean(history_array),
            np.std(history_array),
            np.min(history_array),
            np.max(history_array),
        ]
        
        # Add trend if enough data
        if len(self.history) >= 5:
            recent = history_array[-5:]
            trend = np.mean(np.diff(recent))
            features.append(trend)
        else:
            features.append(0)
        
        return np.array(features)
    
    def _statistical_detect(self, value: float) -> Dict:
        """Fallback statistical anomaly detection (3-sigma rule)"""
        history_array = np.array(list(self.history))
        mean = np.mean(history_array)
        std = np.std(history_array)
        
        if std == 0:
            return {'is_anomaly': False, 'anomaly_score': 0.0, 'severity': 'low'}
        
        z_score = abs((value - mean) / std)
        is_anomaly = z_score > 3
        anomaly_score = min(z_score / 3, 1.0)
        
        return {
            'is_anomaly': is_anomaly,
            'anomaly_score': float(anomaly_score),
            'severity': self._get_severity(anomaly_score),
        }
    
    def _get_severity(self, score: float) -> str:
        """Determine severity based on anomaly score"""
        if score > 0.8:
            return 'critical'
        elif score > 0.6:
            return 'high'
        elif score > 0.4:
            return 'medium'
        else:
            return 'low'

# Store per-service detectors
_detectors = {}

def get_anomaly_detector(service_name: str, metric_type: str) -> AnomalyDetector:
    key = f"{service_name}:{metric_type}"
    if key not in _detectors:
        _detectors[key] = AnomalyDetector()
    return _detectors[key]
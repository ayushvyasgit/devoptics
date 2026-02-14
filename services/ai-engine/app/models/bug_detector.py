import numpy as np
from sklearn.ensemble import GradientBoostingClassifier
from sklearn.preprocessing import StandardScaler
from radon.complexity import cc_visit
from typing import Dict, List
import re

class BugDetector:
    def __init__(self):
        # Initialize ML model
        self.model = GradientBoostingClassifier(
            n_estimators=100,
            learning_rate=0.1,
            max_depth=5,
            random_state=42
        )
        self.scaler = StandardScaler()
        self.is_trained = False
        
        # Pre-train with dummy data for demo
        self._pretrain()
    
    def _pretrain(self):
        """Pre-train with synthetic data"""
        # Generate synthetic training data
        np.random.seed(42)
        X = np.random.randn(1000, 15)
        y = (X[:, 0] + X[:, 1] > 0).astype(int)
        
        X_scaled = self.scaler.fit_transform(X)
        self.model.fit(X_scaled, y)
        self.is_trained = True
        print("✅ Bug detector pre-trained")
    
    def extract_features(self, code: str) -> np.ndarray:
        """Extract features from code"""
        features = []
        
        # 1. Code length features
        features.append(len(code))
        features.append(len(code.split('\n')))
        
        # 2. Cyclomatic complexity
        try:
            complexity_results = cc_visit(code)
            avg_complexity = np.mean([c.complexity for c in complexity_results]) if complexity_results else 0
        except:
            avg_complexity = 0
        features.append(avg_complexity)
        
        # 3. Pattern-based features
        features.append(code.count('if'))
        features.append(code.count('for'))
        features.append(code.count('while'))
        features.append(code.count('try'))
        features.append(code.count('except') + code.count('catch'))
        features.append(code.count('TODO') + code.count('FIXME'))
        
        # 4. Code smells
        features.append(1 if len(code.split('\n')) > 100 else 0)  # Long method
        features.append(1 if avg_complexity > 10 else 0)  # High complexity
        features.append(len(re.findall(r'\bdef\b|\bfunction\b', code)))  # Function count
        
        # 5. Additional patterns
        features.append(code.count('class'))
        features.append(code.count('return'))
        features.append(1 if 'global ' in code else 0)  # Global variables
        
        return np.array(features).reshape(1, -1)
    
    def predict_bug_probability(self, code: str) -> float:
        """Predict probability of bugs in code"""
        if not self.is_trained:
            return 0.5
        
        features = self.extract_features(code)
        features_scaled = self.scaler.transform(features)
        
        # Get probability
        proba = self.model.predict_proba(features_scaled)[0][1]
        
        return float(proba)

# Singleton
_bug_detector = None

def get_bug_detector() -> BugDetector:
    global _bug_detector
    if _bug_detector is None:
        _bug_detector = BugDetector()
    return _bug_detector
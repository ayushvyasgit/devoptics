import torch
from transformers import AutoTokenizer, AutoModel
import numpy as np
from typing import List, Tuple

class CodeBERTAnalyzer:
    def __init__(self):
        print("🔄 Loading CodeBERT model...")
        self.tokenizer = AutoTokenizer.from_pretrained("microsoft/codebert-base")
        self.model = AutoModel.from_pretrained("microsoft/codebert-base")
        
        # Move to GPU if available
        self.device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
        self.model.to(self.device)
        self.model.eval()
        
        print(f"✅ CodeBERT loaded on {self.device}")
    
    def get_code_embedding(self, code: str) -> np.ndarray:
        """Generate semantic embedding for code"""
        # Tokenize
        inputs = self.tokenizer(
            code,
            return_tensors="pt",
            max_length=512,
            truncation=True,
            padding=True
        )
        
        # Move to device
        inputs = {k: v.to(self.device) for k, v in inputs.items()}
        
        # Get embeddings
        with torch.no_grad():
            outputs = self.model(**inputs)
            # Use [CLS] token embedding
            embedding = outputs.last_hidden_state[:, 0, :].cpu().numpy()
        
        return embedding[0]
    
    def calculate_similarity(self, code1: str, code2: str) -> float:
        """Calculate similarity between two code snippets"""
        emb1 = self.get_code_embedding(code1)
        emb2 = self.get_code_embedding(code2)
        
        # Cosine similarity
        similarity = np.dot(emb1, emb2) / (np.linalg.norm(emb1) * np.linalg.norm(emb2))
        
        return float(similarity)
    
    def is_clone(self, code1: str, code2: str, threshold: float = 0.9) -> bool:
        """Check if two code snippets are clones"""
        similarity = self.calculate_similarity(code1, code2)
        return similarity > threshold

# Singleton instance
_codebert_analyzer = None

def get_code_analyzer() -> CodeBERTAnalyzer:
    global _codebert_analyzer
    if _codebert_analyzer is None:
        _codebert_analyzer = CodeBERTAnalyzer()
    return _codebert_analyzer
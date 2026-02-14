import re
from typing import List, Dict

class SecurityScanner:
    def __init__(self):
        # Define vulnerability patterns
        self.vulnerability_patterns = {
            'sql_injection': [
                r'execute\s*\(\s*["\'].*%s',
                r'cursor\.execute\s*\([^?]*\+',
                r'query\s*=\s*["\'].*\+\s*\w+',
            ],
            'xss': [
                r'\.innerHTML\s*=',
                r'document\.write\s*\(',
                r'eval\s*\(',
            ],
            'hardcoded_secrets': [
                r'password\s*=\s*["\'][^"\']+["\']',
                r'api_key\s*=\s*["\'][^"\']+["\']',
                r'secret\s*=\s*["\'][^"\']+["\']',
                r'token\s*=\s*["\'][^"\']+["\']',
            ],
            'command_injection': [
                r'os\.system\s*\(',
                r'subprocess\.call\s*\(',
                r'exec\s*\(',
            ],
        }
        
        self.recommendations = {
            'sql_injection': 'Use parameterized queries or ORM',
            'xss': 'Sanitize user input, use textContent instead of innerHTML',
            'hardcoded_secrets': 'Use environment variables or secure vault',
            'command_injection': 'Use safe APIs, validate and sanitize input',
        }
    
    def scan_code(self, code: str, language: str = 'python') -> List[Dict]:
        """Scan code for security vulnerabilities"""
        vulnerabilities = []
        
        for vuln_type, patterns in self.vulnerability_patterns.items():
            for pattern in patterns:
                matches = re.finditer(pattern, code, re.IGNORECASE)
                
                for match in matches:
                    line_num = code[:match.start()].count('\n') + 1
                    
                    vulnerabilities.append({
                        'type': vuln_type,
                        'severity': self._get_severity(vuln_type),
                        'line': line_num,
                        'column': match.start() - code.rfind('\n', 0, match.start()),
                        'code_snippet': match.group(0),
                        'description': f'Potential {vuln_type.replace("_", " ")} vulnerability',
                        'recommendation': self.recommendations.get(vuln_type, 'Review code'),
                    })
        
        return vulnerabilities
    
    def _get_severity(self, vuln_type: str) -> str:
        """Determine severity of vulnerability"""
        high_severity = ['sql_injection', 'command_injection']
        medium_severity = ['xss', 'hardcoded_secrets']
        
        if vuln_type in high_severity:
            return 'critical'
        elif vuln_type in medium_severity:
            return 'high'
        else:
            return 'medium'

# Singleton
_security_scanner = None

def get_security_scanner() -> SecurityScanner:
    global _security_scanner
    if _security_scanner is None:
        _security_scanner = SecurityScanner()
    return _security_scanner
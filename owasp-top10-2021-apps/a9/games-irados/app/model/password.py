import hashlib
import os
import base64
import hmac
import re

class Password:

    PBKDF2_ITERATIONS = 100_000
    PBKDF2_ALGORITHM = 'sha256'
    PBKDF2_PREFIX = 'pbkdf2_sha256'

    def __init__(self, password):
        self.password = password

    def get_hashed_password(self):
        """Return a salted PBKDF2 hash of the password in the format:
        pbkdf2_sha256$iterations$salt$hash
        """
        salt = base64.b64encode(os.urandom(16)).decode('utf-8')
        dk = hashlib.pbkdf2_hmac(
            self.PBKDF2_ALGORITHM,
            self.password.encode('utf-8'),
            salt.encode('utf-8'),
            self.PBKDF2_ITERATIONS
        )
        hash_b64 = base64.b64encode(dk).decode('utf-8')
        return f"{self.PBKDF2_PREFIX}${self.PBKDF2_ITERATIONS}${salt}${hash_b64}"

    def validate_password(self, hashed_password):
        """Validate a stored hash against the current plaintext password.

        This function supports two formats for backward compatibility:
        - New: pbkdf2_sha256$iterations$salt$hash
        - Legacy: raw SHA-256 hex (64 hex chars)
        """
        return self._compare_password(hashed_password, self.password)

    def _make_hash(self, string):
        # Legacy raw SHA-256 hex (kept for compatibility checks)
        return hashlib.sha256(string.encode('utf-8')).hexdigest()

    def _compare_password(self, stored_hash, plain_password):
        # Use constant-time comparison for all checks
        if stored_hash is None:
            return False

        # New PBKDF2 format: pbkdf2_sha256$iterations$salt$hash
        if isinstance(stored_hash, str) and stored_hash.startswith(self.PBKDF2_PREFIX + "$"):
            try:
                _, iterations_str, salt, hash_b64 = stored_hash.split('$', 3)
                iterations = int(iterations_str)
                dk = hashlib.pbkdf2_hmac(
                    self.PBKDF2_ALGORITHM,
                    plain_password.encode('utf-8'),
                    salt.encode('utf-8'),
                    iterations
                )
                computed_b64 = base64.b64encode(dk).decode('utf-8')
                return hmac.compare_digest(computed_b64, hash_b64)
            except Exception:
                return False

        # Legacy raw SHA-256 hex (64 hex chars)
        if isinstance(stored_hash, str) and re.fullmatch(r'[0-9a-fA-F]{64}', stored_hash):
            computed = self._make_hash(plain_password)
            return hmac.compare_digest(computed, stored_hash)

        # Unknown/unsupported format
        return False

import hashlib
import hmac


class Password:

    def __init__(self, password):
        self.password = password

    def get_hashed_password(self):
        return self._make_hash(self.password)

    def validate_password(self, hashed_password):
        return self._compare_password(hashed_password, self._make_hash(self.password))

    def _make_hash(self, string):
        return hashlib.sha256(string).hexdigest()

    def _compare_password(self, password_1, password_2):
        # Use hmac.compare_digest for constant-time comparison to avoid timing attacks
        try:
            return hmac.compare_digest(password_1, password_2)
        except TypeError:
            # Fallback: ensure inputs are bytes for compare_digest
            if isinstance(password_1, str):
                password_1 = password_1.encode('utf-8')
            if isinstance(password_2, str):
                password_2 = password_2.encode('utf-8')
            return hmac.compare_digest(password_1, password_2)

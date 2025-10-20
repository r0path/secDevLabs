import hashlib
import os
import binascii
import hmac


class Password:

    def __init__(self, password):
        self.password = password

    def get_hashed_password(self):
        """Create a salted, iterated PBKDF2-HMAC-SHA256 hash for storage.

        Format: pbkdf2_sha256$<iterations>$<salt_hex>$<dk_hex>
        """
        return self._make_hash(self.password)

    def validate_password(self, hashed_password):
        """Validate a stored password hash against the current plaintext password.

        This function is backward-compatible: it will verify both the new
        pbkdf2_sha256 format and legacy unsalted SHA-256 hex digests.
        """
        return self._compare_password(hashed_password, self.password)

    def _make_hash(self, string):
        if isinstance(string, str):
            pwd = string.encode('utf-8')
        else:
            pwd = string
        salt = os.urandom(16)
        iterations = 100000
        dk = hashlib.pbkdf2_hmac('sha256', pwd, salt, iterations)
        salt_hex = binascii.hexlify(salt).decode('ascii')
        dk_hex = binascii.hexlify(dk).decode('ascii')
        return f"pbkdf2_sha256${iterations}${salt_hex}${dk_hex}"

    def _compare_password(self, stored_hash, candidate):
        if isinstance(candidate, str):
            candidate_bytes = candidate.encode('utf-8')
        else:
            candidate_bytes = candidate

        # New format: pbkdf2_sha256$iterations$salt_hex$dk_hex
        parts = stored_hash.split('$')
        if len(parts) == 4 and parts[0] == 'pbkdf2_sha256':
            try:
                iterations = int(parts[1])
                salt = binascii.unhexlify(parts[2])
                dk_stored = binascii.unhexlify(parts[3])
            except (ValueError, binascii.Error):
                return False
            dk_candidate = hashlib.pbkdf2_hmac('sha256', candidate_bytes, salt, iterations)
            return hmac.compare_digest(dk_stored, dk_candidate)

        # Legacy unsalted SHA-256 (hex)
        if isinstance(stored_hash, str) and len(stored_hash) == 64:
            try:
                # Ensure it's hex
                binascii.unhexlify(stored_hash)
            except (binascii.Error, TypeError):
                return False
            digest = hashlib.sha256(candidate_bytes).hexdigest()
            return hmac.compare_digest(digest, stored_hash)

        return False

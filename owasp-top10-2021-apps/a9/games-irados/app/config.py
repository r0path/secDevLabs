import os

DEBUG = os.environ.get('DEBUG', 'False')
SECRET_KEY = os.environ.get('SECRET_KEY')

# Security-related configuration (configurable via environment variables).
# These values are provided so the application can apply secure headers and
# cookie attributes in a central, configurable place. Defaults are chosen to
# be minimally disruptive; operators should override via environment
# variables for stricter policies in production.

# SESSION_COOKIE_SECURE: when True, cookies will only be sent over HTTPS.
# Parsed from environment variable SESSION_COOKIE_SECURE (e.g. 'true', '1').
_session_cookie_secure_env = os.environ.get('SESSION_COOKIE_SECURE')
if _session_cookie_secure_env is None:
    SESSION_COOKIE_SECURE = None
else:
    SESSION_COOKIE_SECURE = str(_session_cookie_secure_env).lower() in ('1', 'true', 'yes')

# SESSION_COOKIE_SAMESITE: e.g. 'Lax', 'Strict', or 'None'. Leave unset to keep framework defaults.
SESSION_COOKIE_SAMESITE = os.environ.get('SESSION_COOKIE_SAMESITE')

# SECURITY_HEADERS: a dict of headers that the application can apply to responses.
# Keep conservative defaults that are unlikely to break development setups. Operators
# should set CONTENT_SECURITY_POLICY and STRICT_TRANSPORT_SECURITY in production.
SECURITY_HEADERS = {
    'X-Frame-Options': os.environ.get('X_FRAME_OPTIONS', 'SAMEORIGIN'),
    'X-Content-Type-Options': os.environ.get('X_CONTENT_TYPE_OPTIONS', 'nosniff'),
    'Referrer-Policy': os.environ.get('REFERRER_POLICY', 'no-referrer'),
    # Content-Security-Policy left blank by default to avoid inadvertently breaking
    # pages that load resources from CDNs. Set CONTENT_SECURITY_POLICY in env for production.
    'Content-Security-Policy': os.environ.get('CONTENT_SECURITY_POLICY', ''),
    # Strict-Transport-Security should only be set if the app is served over HTTPS.
    # Example value: 'max-age=63072000; includeSubDomains; preload'
    'Strict-Transport-Security': os.environ.get('STRICT_TRANSPORT_SECURITY', ''),
}

MYSQL_ENDPOINT = os.environ.get('MYSQL_ENDPOINT')
MYSQL_PASSWORD = os.environ.get('MYSQL_PASSWORD')
MYSQL_USER = os.environ.get('MYSQL_USER')
MYSQL_DB = os.environ.get('MYSQL_DB')

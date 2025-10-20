import os

DEBUG = os.environ.get('DEBUG', 'False')
SECRET_KEY = os.environ.get('SECRET_KEY', 'ooops,algo errado!')

MYSQL_ENDPOINT = os.environ.get('MYSQL_ENDPOINT')
MYSQL_PASSWORD = os.environ.get('MYSQL_PASSWORD')
MYSQL_USER = os.environ.get('MYSQL_USER')
MYSQL_DB = os.environ.get('MYSQL_DB')

# Security: ensure session cookies are only sent over HTTPS and have a SameSite policy.
# These default to secure values but can be overridden by environment variables if needed.
SESSION_COOKIE_SECURE = os.environ.get('SESSION_COOKIE_SECURE', 'True') == 'True'
SESSION_COOKIE_SAMESITE = os.environ.get('SESSION_COOKIE_SAMESITE', 'Lax')

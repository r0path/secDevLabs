import os

DEBUG = os.environ.get('DEBUG', 'False')
# Use an environment-provided SECRET_KEY if available; otherwise generate a secure random key at startup.
# NOTE: this will invalidate persisted sessions across restarts but avoids using a hardcoded, guessable secret.
SECRET_KEY = os.environ.get('SECRET_KEY') or os.urandom(24)

MYSQL_ENDPOINT = os.environ.get('MYSQL_ENDPOINT')
MYSQL_PASSWORD = os.environ.get('MYSQL_PASSWORD')
MYSQL_USER = os.environ.get('MYSQL_USER')
MYSQL_DB = os.environ.get('MYSQL_DB')

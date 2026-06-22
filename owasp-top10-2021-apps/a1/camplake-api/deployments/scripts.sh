#!/bin/bash
#
# This script will generate environment passwords to build the environment.
#

# Generating cryptographically secure random credentials
generate_random() {
  local length=$1
  local rand
  if command -v openssl >/dev/null 2>&1; then
    # Calculate required base64 bytes to produce at least `length` alphanumeric chars
    rand=$(openssl rand -base64 $(( (length*3 + 3) / 4 )) | tr -dc 'A-Za-z0-9' | cut -c1-$length)
  else
    # Fallback to /dev/urandom if openssl is not available
    rand=$(LC_ALL=C tr -dc 'A-Za-z0-9' </dev/urandom | head -c "$length")
  fi
  printf '%s' "$rand"
}

# Username starts with 'User' followed by 12 random alphanumeric chars
MONGO_DATABASE_USERNAME_TMP="User$(generate_random 12)"
# Password starts with 'Pass' followed by 24 random alphanumeric chars
MONGO_DATABASE_PASSWORD_TMP="Pass$(generate_random 24)"

# Writing passwords into dockers.env file to be used by docker compose
echo "MONGO_DATABASE_USERNAME=$MONGO_DATABASE_USERNAME_TMP" > deployments/dockers.env
echo "MONGO_DATABASE_PASSWORD=$MONGO_DATABASE_PASSWORD_TMP" >> deployments/dockers.env

# Writing passwords into .env to be used by run_create_cert.sh and to send to STDOUT
echo "MONGO_DATABASE_USERNAME=\"$MONGO_DATABASE_USERNAME_TMP\"" > .env
echo "MONGO_DATABASE_PASSWORD=\"$MONGO_DATABASE_PASSWORD_TMP\"" >> .env

# Preparing script to create mongoDB default user
cat << EOF > deployments/mongo-init.js
var db = connect("mongodb://localhost/DB");
db.createUser(
    {
        user: "${MONGO_DATABASE_USERNAME_TMP}",
        pwd: "${MONGO_DATABASE_PASSWORD_TMP}",
        roles: [{ role: "userAdminAnyDatabase", db: "admin" }]
    }
);
EOF
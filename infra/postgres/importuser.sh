#!/bin/bash

# Keycloak configuration
KEYCLOAK_URL="http://keycloak:8080"
REALM_NAME="myrealm"
CLIENT_ID="admin-cli"
ADMIN_USER="admin"
ADMIN_PASSWORD="admin"

# PostgreSQL configuration
DB_HOST="postgres"
DB_PORT="5432"
DB_NAME="yuno_db"
DB_USER="root"
DB_PASSWORD="root"

# Authenticate with Keycloak to get an access token
TOKEN=$(curl -s -X POST "$KEYCLOAK_URL/realms/master/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=$ADMIN_USER" \
  -d "password=$ADMIN_PASSWORD" \
  -d "grant_type=password" \
  -d "client_id=$CLIENT_ID" | jq -r .access_token)

# Check if the token was obtained successfully
if [ -z "$TOKEN" ]; then
  echo "Error: Failed to obtain Keycloak token."
  exit 1
fi

# Get the list of users in the specified realm
USERS=$(curl -s -X GET "$KEYCLOAK_URL/admin/realms/$REALM_NAME/users" \
  -H "Authorization: Bearer $TOKEN" -H "Content-Type: application/json")

# Insert users into PostgreSQL
for row in $(echo "${USERS}" | jq -r '.[] | @base64'); do
  # Decode the user information
  _jq() {
    echo "${row}" | base64 --decode | jq -r "${1}"
  }

  USER_ID=$(_jq '.id')
  USERNAME=$(_jq '.username')
  EMAIL=$(_jq '.email')

  # Connect to PostgreSQL and execute the insert
  PGPASSWORD=$DB_PASSWORD psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME <<EOF
  INSERT INTO users (user_id, username, email) VALUES ('$USER_ID', '$USERNAME', '$EMAIL')
  ON CONFLICT (user_id) DO NOTHING;
EOF
done

echo "Users synchronized with PostgreSQL database."

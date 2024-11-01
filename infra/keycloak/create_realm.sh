#!/bin/bash

# Keycloak server configuration
KEYCLOAK_URL="http://localhost:8080"
REALM_NAME="myrealm"
CLIENT_ID="myclient"
CLIENT_SECRET="your-client-secret"
ADMIN_USER="admin"
ADMIN_PASSWORD="admin"

# Path to Keycloak CLI in the container
KCADM="/opt/keycloak/bin/kcadm.sh"



# Authenticate with Keycloak
echo "Authenticating with Keycloak..."
$KCADM config credentials --server $KEYCLOAK_URL --realm master \
  --user $ADMIN_USER --password $ADMIN_PASSWORD

# Check if authentication was successful
if [ $? -ne 0 ]; then
  echo "Error: Failed to authenticate with Keycloak. Check admin credentials."
  exit 1
fi

# Create the realm if it doesn’t exist
if ! $KCADM get realms/$REALM_NAME --server $KEYCLOAK_URL; then
  echo "Creating realm '$REALM_NAME'..."
  $KCADM create realms -s realm=$REALM_NAME -s enabled=true --server $KEYCLOAK_URL
fi

# Create the client
echo "Creating client '$CLIENT_ID'..."
$KCADM create clients -r $REALM_NAME -s clientId=$CLIENT_ID -s enabled=true \
  -s secret=$CLIENT_SECRET -s "redirectUris=[\"http://localhost:3000/*\"]" \
  -s publicClient=false -s directAccessGrantsEnabled=true --server $KEYCLOAK_URL

# Create a role in the client
echo "Creating role 'user' in client '$CLIENT_ID'..."
$KCADM create clients/$CLIENT_ID/roles -r $REALM_NAME -s name="user" --server $KEYCLOAK_URL

# Create a user in the realm with additional fields
USER_NAME="testuser"
USER_PASSWORD="password"
USER_EMAIL="testuser@example.com"
USER_FIRST_NAME="Test"
USER_LAST_NAME="User"

echo "Creating user '$USER_NAME'..."
$KCADM create users -r $REALM_NAME -s username=$USER_NAME -s enabled=true \
  -s email=$USER_EMAIL -s firstName=$USER_FIRST_NAME -s lastName=$USER_LAST_NAME \
  --server $KEYCLOAK_URL

# Set the user's password
$KCADM set-password -r $REALM_NAME --username $USER_NAME --new-password $USER_PASSWORD --server $KEYCLOAK_URL



# Makefile

# Variables
VAULT_TOKEN=root

enable-transit:
	- docker-compose exec -e VAULT_TOKEN=$(VAULT_TOKEN) vault vault secrets enable transit

create-realm:
	- docker-compose exec keycloak /opt/keycloak/create_realm.sh

import-users:
	- docker-compose exec postgres /opt/importuser.sh

up:
	- docker-compose up -d

down:
	- docker-compose down -v

restart: down up


all: enable-transit create-realm import-users

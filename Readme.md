Here's the updated `README.md` with an example of how to obtain a token using `curl`:

---

# Project Setup

This project uses Docker Compose to run multiple services, including HashiCorp Vault, Keycloak, and PostgreSQL. The `Makefile` provides convenient commands to start, stop, and configure these services.

## Prerequisites

- **Docker** and **Docker Compose** must be installed on your machine.
- Ensure the `make` command is available in your environment.

## Quick Start

To start the services and set up the environment, follow these steps:

### 1. Start the Services

Run the following command to start all services:

```bash
make up
```

This will initialize all services defined in the `docker-compose.yml` file and run them in the background.

### 2. Wait for All Services to Start

Give the services a moment to initialize completely. You can check their status using:

```bash
docker-compose ps
```

Ensure all services are listed as `running` before proceeding to the next step.

### 3. Configure the Environment

Once all services are running, execute the following command to configure Vault, Keycloak, and PostgreSQL:

```bash
make all
```

This command performs the following actions in sequence:
- Enables the **transit secrets engine** in Vault.
- Creates the required **realm** in Keycloak.
- Imports **users** into the PostgreSQL database.

### 4. Obtain a Token from Keycloak

To authenticate a user and obtain a token from Keycloak, use the following `curl` command:

```bash
curl -X POST "http://localhost:8080/realms/myrealm/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "grant_type=password" \
  -d "client_id=myclient" \
  -d "client_secret=your-client-secret" \
  -d "username=testuser" \
  -d "password=password" | jq .
```

Replace the following parameters as needed:
- `myrealm`: The name of your Keycloak realm.
- `myclient`: The client ID in Keycloak.
- `your-client-secret`: The client secret for your Keycloak client.
- `testuser`: The username of the user you want to authenticate.
- `password`: The password of the user.

The command will return a JSON response containing the access token and other token details, which are parsed with `jq` for readability.

### Additional Commands

For convenience, here are a few other helpful `Makefile` commands:

- **Stop Services**:

  ```bash
  make down
  ```

  This stops all services and removes volumes to clear any stored data.

- **Restart Services**:

  ```bash
  make restart
  ```

  This command restarts the services, removing volumes first to ensure a clean start.

Aquí tienes el fragmento en inglés para agregar al `README.md` sobre cómo cargar una tarjeta en la aplicación:

---

### Storing a Card

To load a card, follow these steps:

1. **Obtain a Keycloak Token**: First, you need to obtain an access token from Keycloak. You can do this by running the `curl` command provided earlier in this README.

2. **Create a Key Pair**: Once you have the access token, make a request to `[POST] /keys` to generate a new key pair. Use the access token in the request headers to authenticate.

3. **Encrypt the PAN**: Use the returned public key to encrypt the PAN (Primary Account Number). There is an example in the `examples/` directory showing how to perform this encryption.

4. **Submit the Encrypted Card Data**: After encrypting the PAN, make a request to `[POST] /cards` to submit the card data securely.

### Swagger for API Testing

All internal endpoints of the application are available in `/swagger`, allowing you to test them directly in the API documentation interface.

--- 

## Notes

- **Vault Transit Engine**: The `enable-transit` step configures Vault with the transit secrets engine, which is required for encryption services.
- **Keycloak Realm Setup**: The `create-realm` command sets up Keycloak with the necessary authentication realm for your application.
- **User Import**: The `import-users` command loads predefined user data into PostgreSQL.

These steps ensure that all required services and configurations are properly initialized.

---

This updated `README.md` includes instructions for obtaining a Keycloak token, along with the overall project setup and additional Makefile commands.

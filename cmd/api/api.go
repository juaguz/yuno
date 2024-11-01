package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	vault "github.com/hashicorp/vault/api"
	"github.com/joho/godotenv"
	_ "github.com/juaguz/yuno/docs"
	"github.com/juaguz/yuno/internal/cards"
	"github.com/juaguz/yuno/internal/cards/dtos"
	"github.com/juaguz/yuno/internal/cards/repositories"
	"github.com/juaguz/yuno/internal/keys"
	"github.com/juaguz/yuno/kit/database"
	"github.com/juaguz/yuno/kit/kms"
	"github.com/juaguz/yuno/kit/users/auth"
	"github.com/juaguz/yuno/kit/users/repository"
	kitvault "github.com/juaguz/yuno/kit/vault"
	"github.com/juaguz/yuno/pkg/cards/api"
	keysApi "github.com/juaguz/yuno/pkg/keys/api"
	httpSwagger "github.com/swaggo/http-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func jsonResponseMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		next.ServeHTTP(w, r)
	})
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	vaultAddress := os.Getenv("VAULT_ADDRESS")
	v, err := vault.NewClient(&vault.Config{
		Address: vaultAddress,
	})
	if err != nil {
		panic(err)
	}

	v.SetToken(os.Getenv("VAULT_TOKEN"))

	cardRepo := repositories.NewCardRepository(db)

	vaultService := kitvault.NewVaultService(v)

	kmsService := kms.NewVaultKmsService(v)

	cardService := cards.NewCardService(cardRepo, kmsService, vaultService)

	transactionalService := database.NewTransactionalRepository[dtos.Card](db, cardService)

	batchupdater := cards.NewBatchUpdater(cardService)

	cardsHandler := api.NewCardHandler(transactionalService, batchupdater)

	userRepo := repository.NewUserRepository(db)

	keycloakURL := os.Getenv("KEYCLOAK_URL")
	realm := os.Getenv("KEYCLOAK_REALM")

	jwtMiddleware := auth.JWTMiddleware(userRepo, keycloakURL, realm)

	keysProvider := keys.NewKeysProvider(kmsService)
	keysHandler := keysApi.NewKeysHandler(keysProvider)

	r := chi.NewRouter()
	r.Use(jwtMiddleware)
	r.Use(jsonResponseMiddleware)

	// @securityDefinitions.apikey Bearer
	// @in header
	// @name Authorization
	// @description Type "Bearer" followed by a space and JWT token.
	r.Mount("/cards", cardsHandler.Routes())
	r.Mount("/keys", keysHandler.Routes())
	r.Get("/swagger/*", httpSwagger.WrapHandler)
	if err := http.ListenAndServe(os.Getenv("APP_ADDRESS"), r); err != nil {
		log.Fatalf("error starting server: %s", err)
	}
}

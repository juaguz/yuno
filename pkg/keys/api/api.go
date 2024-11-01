package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/juaguz/yuno/internal/keys"
	"github.com/juaguz/yuno/kit/users/auth"
)

type KeysHandler struct {
	Service *keys.KeysProvider
}

func NewKeysHandler(service *keys.KeysProvider) *KeysHandler {
	return &KeysHandler{Service: service}
}

// Routes configures the routes for KeysHandler
func (h *KeysHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreateKey)

	return r
}

// CreateKey godoc
// @Summary Create a new key
// @Description Generates a new public key for the authenticated user
// @Tags keys
// @Produce json
// @Success 201 {object} KeysResponse
// @Failure 500 {string} string "Internal server error"
// @Router /keys [post]
// @Security Bearer
func (h *KeysHandler) CreateKey(w http.ResponseWriter, r *http.Request) {
	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	key, err := h.Service.CreateKey(r.Context(), user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	k := &KeysResponse{
		PublicKey: key,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(k)
}

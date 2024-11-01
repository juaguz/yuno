package api

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/juaguz/yuno/internal/cards/dtos"
	"github.com/juaguz/yuno/kit/errors/senital"
	"github.com/juaguz/yuno/kit/users/auth"
)

type CardHandler struct {
	Service            Service[dtos.Card]
	BatchUpdateService BatchUpdate
}

type Service[T any] interface {
	Create(ctx context.Context, entity *T) (*T, error)
	Update(ctx context.Context, entity *T) error
	Delete(ctx context.Context, entity *T) error
	Get(ctx context.Context, entity *T) (*T, error)
}

type BatchUpdate interface {
	Update(ctx context.Context, userID uuid.UUID, cards []*dtos.BatchUpdate) ([]*dtos.BatchUpdateStatus, error)
}

func NewCardHandler(service Service[dtos.Card], batchUpdate BatchUpdate) *CardHandler {
	return &CardHandler{Service: service, BatchUpdateService: batchUpdate}
}

// Routes configures the routes for CardHandler
func (h *CardHandler) Routes() chi.Router {
	r := chi.NewRouter()

	r.Post("/", h.CreateCard)
	r.Get("/{cardID}", h.GetCard)
	r.Put("/{cardID}", h.UpdateCard)
	r.Delete("/{cardID}", h.DeleteCard)
	r.Put("/batch", h.BatchUpdate)

	return r
}

// CreateCard godoc
// @Summary Create a new card
// @Description Create a card with card holder and PAN
// @Tags cards
// @Accept json
// @Produce json
// @Param card body CardCreation true "Card Creation Request"
// @Success 201 {object} dtos.Card
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /cards [post]
// @Security Bearer
func (h *CardHandler) CreateCard(w http.ResponseWriter, r *http.Request) {
	card := &CardCreation{}
	if err := json.NewDecoder(r.Body).Decode(card); err != nil {
		log.Printf("error %s", err.Error())
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	c := dtos.Card{
		CardHolder: card.CardHolder,
		Pan:        card.Pan,
		UserId:     user.ID,
	}
	res, err := h.Service.Create(r.Context(), &c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(res)
}

// GetCard godoc
// @Summary Get a card
// @Description Retrieve a card by its ID
// @Tags cards
// @Produce json
// @Param cardID path string true "Card ID"
// @Success 200 {object} dtos.Card
// @Failure 400 {string} string "Invalid card ID"
// @Failure 404 {string} string "Card not found"
// @Failure 500 {string} string "Internal server error"
// @Router /cards/{cardID} [get]
// @Security Bearer
func (h *CardHandler) GetCard(w http.ResponseWriter, r *http.Request) {
	cardID, err := uuid.Parse(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, "invalid card ID", http.StatusBadRequest)
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	c := &dtos.Card{
		ID:     cardID,
		UserId: user.ID,
	}

	card, err := h.Service.Get(r.Context(), c)
	if err != nil {
		if errors.Is(err, senital.ErrNotFound) {
			http.Error(w, "card not found", http.StatusNotFound)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	json.NewEncoder(w).Encode(card)
}

// UpdateCard godoc
// @Summary Update a card
// @Description Update a card's details by its ID
// @Tags cards
// @Accept json
// @Param cardID path string true "Card ID"
// @Param card body CardUpdate true "Card Update Request"
// @Success 204 "No content"
// @Failure 400 {string} string "Invalid request body or card ID"
// @Failure 500 {string} string "Internal server error"
// @Router /cards/{cardID} [put]
// @Security Bearer
func (h *CardHandler) UpdateCard(w http.ResponseWriter, r *http.Request) {
	var body *CardUpdate
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	cardID, err := uuid.Parse(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, "invalid card ID", http.StatusBadRequest)
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	card := dtos.Card{
		ID:         cardID,
		CardHolder: body.CardHolder,
		UserId:     user.ID,
	}

	if err := h.Service.Update(r.Context(), &card); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DeleteCard godoc
// @Summary Delete a card
// @Description Delete a card by its ID
// @Tags cards
// @Param cardID path string true "Card ID"
// @Success 204 "No content"
// @Failure 400 {string} string "Invalid card ID"
// @Failure 500 {string} string "Internal server error"
// @Router /cards/{cardID} [delete]
// @Security Bearer
func (h *CardHandler) DeleteCard(w http.ResponseWriter, r *http.Request) {
	cardID, err := uuid.Parse(chi.URLParam(r, "cardID"))
	if err != nil {
		http.Error(w, "invalid card ID", http.StatusBadRequest)
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	c := &dtos.Card{
		ID:     cardID,
		UserId: user.ID,
	}

	if err := h.Service.Delete(r.Context(), c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// BatchUpdate godoc
// @Summary Batch update cards
// @Description Update multiple cards in a single request
// @Tags cards
// @Accept json
// @Produce json
// @Param batch body []dtos.BatchUpdate true "Batch Update Request"
// @Success 200 {array} dtos.BatchUpdateStatus
// @Failure 400 {string} string "Invalid request body"
// @Failure 500 {string} string "Internal server error"
// @Router /cards/batch [put]
// @Security Bearer
func (h *CardHandler) BatchUpdate(w http.ResponseWriter, r *http.Request) {
	var batch []*dtos.BatchUpdate
	if err := json.NewDecoder(r.Body).Decode(&batch); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	user, err := auth.GetUserFromContext(r.Context())
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	statuses, err := h.BatchUpdateService.Update(r.Context(), user.ID, batch)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(statuses)
}

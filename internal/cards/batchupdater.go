package cards

import (
	"context"

	"github.com/google/uuid"
	"github.com/juaguz/yuno/internal/cards/dtos"
)

type BatchUpdater struct {
	CardService *CardService
}

func NewBatchUpdater(cardService *CardService) *BatchUpdater {
	return &BatchUpdater{
		CardService: cardService,
	}
}

func (b *BatchUpdater) Update(ctx context.Context, userID uuid.UUID, cards []*dtos.BatchUpdate) ([]*dtos.BatchUpdateStatus, error) {
	const numWorkers = 15
	jobs := make(chan *dtos.BatchUpdate, len(cards))
	results := make(chan *dtos.BatchUpdateStatus, len(cards))
	for w := 0; w < numWorkers; w++ {
		go b.worker(ctx, jobs, results, userID)
	}

	for _, card := range cards {
		jobs <- card
	}
	close(jobs)

	var updateStatus []*dtos.BatchUpdateStatus
	for i := 0; i < len(cards); i++ {
		updateStatus = append(updateStatus, <-results)
	}
	close(results)

	return updateStatus, nil
}

func (b *BatchUpdater) worker(ctx context.Context, jobs <-chan *dtos.BatchUpdate, results chan<- *dtos.BatchUpdateStatus, userID uuid.UUID) {
	for card := range jobs {
		c := &dtos.Card{
			ID:         card.ID,
			CardHolder: card.CardHolder,
			UserId:     userID,
		}

		err := b.CardService.Update(ctx, c)
		status := &dtos.BatchUpdateStatus{
			CardID: c.ID,
		}
		if err != nil {
			status.Status = dtos.Failed
		} else {
			status.Status = dtos.Succeeded
		}
		results <- status
	}
}

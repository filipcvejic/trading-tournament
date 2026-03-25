package trackedtrade

import (
	"context"
	"strings"
)

type Service struct {
	repository Repository
}

func NewService(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}

func (s *Service) IngestEvent(ctx context.Context, req IngestTrackedTradeEventRequest) error {
	switch req.EventType {
	case EventTypeOpen:
		if req.PositionID == 0 || req.Symbol == "" || req.Side == "" || req.OpenPrice == 0 || req.Volume <= 0 || req.OpenedAt == nil {
			return ErrMissingOpenFields
		}

		side := Side(strings.ToUpper(req.Side))
		if side != SideBuy && side != SideSell {
			return ErrMissingOpenFields
		}

		trade := TrackedTrade{
			PositionID: req.PositionID,
			Symbol:     req.Symbol,
			Side:       side,
			OpenPrice:  req.OpenPrice,
			Volume:     req.Volume,
			StopLoss:   req.StopLoss,
			OpenedAt:   *req.OpenedAt,
		}

		return s.repository.Create(ctx, trade)

	case EventTypeUpdate:
		if req.PositionID == 0 {
			return ErrMissingUpdateFields
		}

		return s.repository.UpdateStopLoss(ctx, req.PositionID, req.StopLoss)

	case EventTypeClose:
		if req.PositionID == 0 || req.ClosedAt == nil {
			return ErrMissingCloseFields
		}

		return s.repository.Close(ctx, req.PositionID, *req.ClosedAt)

	default:
		return ErrInvalidEventType
	}
}

func (s *Service) List(ctx context.Context) ([]TrackedTrade, error) {
	return s.repository.List(ctx)
}

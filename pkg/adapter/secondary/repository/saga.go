package repository

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/haandol/hexagonal/pkg/constant"
	"github.com/haandol/hexagonal/pkg/constant/status"
	"github.com/haandol/hexagonal/pkg/dto"
	"github.com/haandol/hexagonal/pkg/entity"
	"github.com/haandol/hexagonal/pkg/message/command"
	"github.com/haandol/hexagonal/pkg/message/event"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SagaRepository struct {
	db *gorm.DB
}

func NewSagaRepository(db *gorm.DB) *SagaRepository {
	return &SagaRepository{
		db: db,
	}
}

func (r *SagaRepository) Start(ctx context.Context, cmd *command.StartSaga) (dto.Saga, error) {
	history, err := json.Marshal(&[]any{
		cmd,
	})
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{
		CorrelationID: cmd.CorrelationID,
		TripID:        cmd.Body.TripID,
		CarID:         cmd.Body.CarID,
		HotelID:       cmd.Body.HotelID,
		FlightID:      cmd.Body.FlightID,
		History:       history,
		Status:        status.SagaStarted,
	}
	result := r.db.WithContext(ctx).Create(row)
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}

	return row.DTO()
}

func (r *SagaRepository) ProcessCarBooking(ctx context.Context, evt *event.CarBooked) (dto.Saga, error) {
	saga, err := r.GetByCorrelationID(ctx, evt.CorrelationID)
	if err != nil {
		return dto.Saga{}, err
	}

	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return dto.Saga{}, err
	}
	v = append(v, evt)

	history, err := json.Marshal(&v)
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{}
	result := r.db.WithContext(ctx).
		Model(row).
		Clauses(clause.Returning{}).
		Where("id = ?", saga.ID).
		Updates(&entity.Saga{
			CarBookingID: evt.Body.BookingID,
			Status:       evt.Name,
			History:      history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) CompensateCarBooking(ctx context.Context, evt *event.CarBookingCancelled) (dto.Saga, error) {
	saga, err := r.GetByCorrelationID(ctx, evt.CorrelationID)
	if err != nil {
		return dto.Saga{}, err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return dto.Saga{}, err
	}
	v = append(v, evt)

	history, err := json.Marshal(&v)
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{}
	result := r.db.WithContext(ctx).
		Model(row).
		Clauses(clause.Returning{}).
		Where("correlation_id = ?", evt.CorrelationID).
		Updates(&entity.Saga{
			CarBookingID: 0,
			Status:       evt.Name,
			History:      history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) ProcessHotelBooking(ctx context.Context, evt *event.HotelBooked) (dto.Saga, error) {
	saga, err := r.GetByCorrelationID(ctx, evt.CorrelationID)
	if err != nil {
		return dto.Saga{}, err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return dto.Saga{}, err
	}
	v = append(v, evt)

	history, err := json.Marshal(&v)
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{}
	result := r.db.WithContext(ctx).
		Model(row).
		Clauses(clause.Returning{}).
		Where("correlation_id = ?", evt.CorrelationID).
		Updates(&entity.Saga{
			HotelBookingID: evt.Body.BookingID,
			Status:         evt.Name,
			History:        history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) CompensateHotelBooking(ctx context.Context, evt *event.HotelBookingCancelled) (dto.Saga, error) {
	saga, err := r.GetByCorrelationID(ctx, evt.CorrelationID)
	if err != nil {
		return dto.Saga{}, err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return dto.Saga{}, err
	}
	v = append(v, evt)

	history, err := json.Marshal(&v)
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{}
	result := r.db.WithContext(ctx).
		Model(row).
		Clauses(clause.Returning{}).
		Where("correlation_id = ?", evt.CorrelationID).
		Updates(&entity.Saga{
			HotelBookingID: 0,
			Status:         evt.Name,
			History:        history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) ProcessFlightBooking(ctx context.Context, evt *event.FlightBooked) (dto.Saga, error) {
	saga, err := r.GetByCorrelationID(ctx, evt.CorrelationID)
	if err != nil {
		return dto.Saga{}, err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return dto.Saga{}, err
	}
	v = append(v, evt)

	history, err := json.Marshal(&v)
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{}
	result := r.db.WithContext(ctx).
		Model(row).
		Clauses(clause.Returning{}).
		Where("correlation_id = ?", evt.CorrelationID).
		Updates(&entity.Saga{
			FlightBookingID: evt.Body.BookingID,
			Status:          evt.Name,
			History:         history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) CompensateFlightBooking(ctx context.Context, evt *event.FlightBookingCancelled) (dto.Saga, error) {
	saga, err := r.GetByCorrelationID(ctx, evt.CorrelationID)
	if err != nil {
		return dto.Saga{}, err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return dto.Saga{}, err
	}
	v = append(v, evt)

	history, err := json.Marshal(&v)
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{}
	result := r.db.WithContext(ctx).
		Model(row).
		Clauses(clause.Returning{}).
		Where("correlation_id = ?", evt.CorrelationID).
		Updates(&entity.Saga{
			FlightBookingID: 0,
			Status:          evt.Name,
			History:         history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) End(ctx context.Context, cmd *command.EndSaga) (dto.Saga, error) {
	saga, err := r.GetByCorrelationID(ctx, cmd.CorrelationID)
	if err != nil {
		return dto.Saga{}, err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return dto.Saga{}, err
	}
	v = append(v, cmd)
	history, err := json.Marshal(&v)
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{}
	result := r.db.WithContext(ctx).
		Model(row).
		Clauses(clause.Returning{}).
		Where("correlation_id = ?", cmd.CorrelationID).
		Updates(&entity.Saga{
			Status:  status.SagaEnded,
			History: history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) Abort(ctx context.Context, cmd *command.AbortSaga) (dto.Saga, error) {
	saga, err := r.GetByCorrelationID(ctx, cmd.CorrelationID)
	if err != nil {
		return dto.Saga{}, err
	}
	v := []any{}
	if err := json.Unmarshal([]byte(saga.History), &v); err != nil {
		return dto.Saga{}, err
	}
	v = append(v, cmd)
	history, err := json.Marshal(&v)
	if err != nil {
		return dto.Saga{}, err
	}

	row := &entity.Saga{}
	result := r.db.WithContext(ctx).
		Model(row).
		Clauses(clause.Returning{}).
		Where("correlation_id = ?", cmd.CorrelationID).
		Updates(&entity.Saga{
			Status:  status.SagaAborted,
			History: history,
		})
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}
	if result.RowsAffected == 0 {
		return dto.Saga{}, errors.New("no rows affected")
	}

	return row.DTO()
}

func (r *SagaRepository) GetById(ctx context.Context, id uint) (dto.Saga, error) {
	row := &entity.Saga{}

	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	result := db.
		Where("id = ?", id).
		Take(row)
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}

	return row.DTO()
}

func (r *SagaRepository) GetByCorrelationID(ctx context.Context, id string) (dto.Saga, error) {
	row := &entity.Saga{}

	var db *gorm.DB
	if tx, ok := ctx.Value(constant.TX("tx")).(*gorm.DB); ok {
		db = tx
	} else {
		db = r.db.WithContext(ctx)
	}

	result := db.
		Where("correlation_id = ?", id).
		Find(row)
	if result.Error != nil {
		return dto.Saga{}, result.Error
	}

	return row.DTO()
}

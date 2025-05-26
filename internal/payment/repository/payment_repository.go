package repository

import (
	"context"
	"time"

	"govo/internal/payment/model"

	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(ctx context.Context, payment *model.Payment) (*model.Payment, error) {
	if err := r.db.WithContext(ctx).Create(payment).Error; err != nil {
		return nil, err
	}
	return payment, nil
}

func (r *PaymentRepository) GetByID(ctx context.Context, id uint) (*model.Payment, error) {
	var payment model.Payment
	if err := r.db.WithContext(ctx).First(&payment, id).Error; err != nil {
		return nil, err
	}
	return &payment, nil
}

func (r *PaymentRepository) List(ctx context.Context, customerID uint, status string, startDate, endDate *time.Time) ([]*model.Payment, error) {
	var payments []*model.Payment
	query := r.db.WithContext(ctx).Where("customer_id = ?", customerID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if startDate != nil {
		query = query.Where("created_at >= ?", startDate)
	}

	if endDate != nil {
		query = query.Where("created_at <= ?", endDate)
	}

	if err := query.Find(&payments).Error; err != nil {
		return nil, err
	}

	return payments, nil
}

func (r *PaymentRepository) Update(ctx context.Context, payment *model.Payment) error {
	return r.db.WithContext(ctx).Save(payment).Error
}

func (r *PaymentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Payment{}, id).Error
}

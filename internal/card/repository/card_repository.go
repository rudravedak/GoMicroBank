package repository

import (
	"govo/internal/card/model"

	"gorm.io/gorm"
)

type CardRepository struct {
	db *gorm.DB
}

func NewCardRepository(db *gorm.DB) *CardRepository {
	return &CardRepository{db: db}
}

func (r *CardRepository) Create(card *model.Card) error {
	return r.db.Create(card).Error
}

func (r *CardRepository) GetByID(id uint) (*model.Card, error) {
	var card model.Card
	err := r.db.First(&card, id).Error
	if err != nil {
		return nil, err
	}
	return &card, nil
}

func (r *CardRepository) GetByCustomerID(customerID uint) ([]*model.Card, error) {
	var cards []*model.Card
	err := r.db.Where("customer_id = ?", customerID).Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *CardRepository) Delete(id uint) error {
	return r.db.Delete(&model.Card{}, id).Error
}

func (r *CardRepository) GetCustomerCards(customerID uint) ([]*model.Card, error) {
	var cards []*model.Card
	err := r.db.Where("customer_id = ?", customerID).Find(&cards).Error
	if err != nil {
		return nil, err
	}
	return cards, nil
}

func (r *CardRepository) AddCard(card *model.Card) error {
	return r.db.Create(card).Error
}

func (r *CardRepository) RemoveCard(customerID uint, cardNumber string) error {
	return r.db.Where("customer_id = ? AND card_number = ?", customerID, cardNumber).Delete(&model.Card{}).Error
}

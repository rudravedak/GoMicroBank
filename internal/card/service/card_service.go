package service

import (
	"govo/internal/card/model"
	"govo/internal/card/repository"
)

type CardService struct {
	repo *repository.CardRepository
}

func NewCardService(repo *repository.CardRepository) *CardService {
	return &CardService{repo: repo}
}

func (s *CardService) GetCustomerCards(customerID uint) ([]*model.Card, error) {
	return s.repo.GetCustomerCards(customerID)
}

func (s *CardService) AddCard(customerID uint, cardNumber string, cardType string, expiryDate string, cvv string, creditLimit float64, balance float64) error {
	card := &model.Card{
		CustomerID:  customerID,
		CardNumber:  cardNumber,
		CardType:    cardType,
		ExpiryDate:  expiryDate,
		CVV:         cvv,
		CreditLimit: creditLimit,
		Balance:     balance,
		IsActive:    true,
	}
	return s.repo.Create(card)
}

func (s *CardService) RemoveCard(customerID uint, cardNumber string) error {
	return s.repo.RemoveCard(customerID, cardNumber)
}

func (s *CardService) GetCardByID(id uint) (*model.Card, error) {
	return s.repo.GetByID(id)
}

func (s *CardService) GetCardsByCustomerID(customerID uint) ([]*model.Card, error) {
	return s.repo.GetByCustomerID(customerID)
}

func (s *CardService) DeleteCard(id uint) error {
	return s.repo.Delete(id)
}

package service

import (
	"govo/internal/customer/model"
	"govo/internal/customer/repository"
)

type CustomerService struct {
	repo *repository.CustomerRepository
}

func NewCustomerService(repo *repository.CustomerRepository) *CustomerService {
	return &CustomerService{repo: repo}
}

func (s *CustomerService) CreateCustomer(customer *model.Customer) error {
	return s.repo.Create(customer)
}

func (s *CustomerService) GetCustomer(id uint) (*model.Customer, error) {
	return s.repo.GetByID(id)
}

func (s *CustomerService) UpdateCustomer(customer *model.Customer) error {
	return s.repo.Update(customer)
}

func (s *CustomerService) DeleteCustomer(id uint) error {
	return s.repo.Delete(id)
}

func (s *CustomerService) ListCustomers() ([]model.Customer, error) {
	return s.repo.List()
}

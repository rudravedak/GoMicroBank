package main

import (
	"context"
	"log"
	"net"

	"govo/api/proto/customer"
	"govo/internal/customer/model"
	"govo/internal/customer/repository"
	"govo/internal/customer/service"

	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CustomerServer struct {
	customer.UnimplementedCustomerServiceServer
	service *service.CustomerService
}

func (s *CustomerServer) CreateCustomer(ctx context.Context, req *customer.CreateCustomerRequest) (*customer.CreateCustomerResponse, error) {
	c := &model.Customer{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		Balance:   float64(req.Balance),
	}

	if err := s.service.CreateCustomer(c); err != nil {
		return nil, err
	}

	return &customer.CreateCustomerResponse{
		Id:        uint32(c.ID),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		Phone:     c.Phone,
		Address:   c.Address,
		Balance:   float32(c.Balance),
		Cards:     c.Cards,
	}, nil
}

func (s *CustomerServer) GetCustomer(ctx context.Context, req *customer.GetCustomerRequest) (*customer.GetCustomerResponse, error) {
	c, err := s.service.GetCustomer(uint(req.Id))
	if err != nil {
		return nil, err
	}

	return &customer.GetCustomerResponse{
		Id:        uint32(c.ID),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		Phone:     c.Phone,
		Address:   c.Address,
		Balance:   float32(c.Balance),
		Cards:     c.Cards,
	}, nil
}

func (s *CustomerServer) UpdateCustomer(ctx context.Context, req *customer.UpdateCustomerRequest) (*customer.UpdateCustomerResponse, error) {
	c := &model.Customer{
		ID:        uint(req.Id),
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Address:   req.Address,
		Balance:   float64(req.Balance),
	}

	if err := s.service.UpdateCustomer(c); err != nil {
		return nil, err
	}

	return &customer.UpdateCustomerResponse{
		Id:        uint32(c.ID),
		FirstName: c.FirstName,
		LastName:  c.LastName,
		Email:     c.Email,
		Phone:     c.Phone,
		Address:   c.Address,
		Balance:   float32(c.Balance),
		Cards:     c.Cards,
	}, nil
}

func (s *CustomerServer) DeleteCustomer(ctx context.Context, req *customer.DeleteCustomerRequest) (*customer.DeleteCustomerResponse, error) {
	if err := s.service.DeleteCustomer(uint(req.Id)); err != nil {
		return nil, err
	}

	return &customer.DeleteCustomerResponse{
		Success: true,
	}, nil
}

func (s *CustomerServer) ListCustomers(ctx context.Context, req *customer.ListCustomersRequest) (*customer.ListCustomersResponse, error) {
	customers, err := s.service.ListCustomers()
	if err != nil {
		return nil, err
	}

	response := &customer.ListCustomersResponse{
		Customers: make([]*customer.GetCustomerResponse, len(customers)),
	}

	for i, c := range customers {
		response.Customers[i] = &customer.GetCustomerResponse{
			Id:        uint32(c.ID),
			FirstName: c.FirstName,
			LastName:  c.LastName,
			Email:     c.Email,
			Phone:     c.Phone,
			Address:   c.Address,
			Balance:   float32(c.Balance),
			Cards:     c.Cards,
		}
	}

	return response, nil
}

func main() {
	// PostgreSQL bağlantısı
	dsn := "host=postgres user=postgres password=postgres dbname=customerdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	// Tabloları oluştur
	if err := db.AutoMigrate(&model.Customer{}); err != nil {
		log.Fatalf("Tablo oluşturulamadı: %v", err)
	}

	// Dependency injection
	customerRepo := repository.NewCustomerRepository(db)
	customerService := service.NewCustomerService(customerRepo)
	customerServer := &CustomerServer{service: customerService}

	// gRPC server'ı başlat
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Port dinlenemedi: %v", err)
	}

	grpcServer := grpc.NewServer()
	customer.RegisterCustomerServiceServer(grpcServer, customerServer)

	log.Println("Customer servisi 50052 portunda başlatılıyor...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Server başlatılamadı: %v", err)
	}
}

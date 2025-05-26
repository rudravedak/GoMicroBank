package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	cardpb "govo/api/proto/card"
	"govo/internal/card/handler"
	"govo/internal/card/model"
	"govo/internal/card/repository"
	"govo/internal/card/service"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type CardServer struct {
	cardpb.UnimplementedCardServiceServer
	service *service.CardService
}

func (s *CardServer) GetCustomerCards(ctx context.Context, req *cardpb.GetCustomerCardsRequest) (*cardpb.GetCustomerCardsResponse, error) {
	cards, err := s.service.GetCustomerCards(uint(req.CustomerId))
	if err != nil {
		return nil, err
	}

	response := &cardpb.GetCustomerCardsResponse{
		Cards: make([]*cardpb.GetCardResponse, len(cards)),
	}

	for i, c := range cards {
		response.Cards[i] = &cardpb.GetCardResponse{
			Id:          uint32(c.ID),
			CustomerId:  uint32(c.CustomerID),
			CardNumber:  c.CardNumber,
			CardType:    c.CardType,
			ExpiryDate:  c.ExpiryDate,
			CreditLimit: float32(c.CreditLimit),
			Balance:     float32(c.Balance),
		}
	}

	return response, nil
}

func (s *CardServer) AddCard(ctx context.Context, req *cardpb.AddCardRequest) (*cardpb.AddCardResponse, error) {
	err := s.service.AddCard(
		uint(req.CustomerId),
		req.CardNumber,
		req.CardType,
		req.ExpiryDate,
		req.Cvv,
		float64(req.CreditLimit),
		float64(req.Balance),
	)
	if err != nil {
		return nil, err
	}

	return &cardpb.AddCardResponse{
		Success: true,
	}, nil
}

func (s *CardServer) RemoveCard(ctx context.Context, req *cardpb.RemoveCardRequest) (*cardpb.RemoveCardResponse, error) {
	err := s.service.RemoveCard(uint(req.CustomerId), req.CardNumber)
	if err != nil {
		return nil, err
	}

	return &cardpb.RemoveCardResponse{
		Success: true,
	}, nil
}

func main() {
	// PostgreSQL bağlantısı
	dsn := "host=postgres user=postgres password=postgres dbname=carddb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	// Tabloları oluştur
	if err := db.AutoMigrate(&model.Card{}); err != nil {
		log.Fatalf("Tablo oluşturulamadı: %v", err)
	}

	// Dependency injection
	cardRepo := repository.NewCardRepository(db)
	cardService := service.NewCardService(cardRepo)
	cardServer := &CardServer{service: cardService}
	cardHandler := handler.NewCardHandler(cardService)

	// HTTP router
	router := mux.NewRouter()
	router.HandleFunc("/api/cards", cardHandler.CreateCard).Methods("POST")
	router.HandleFunc("/api/cards", cardHandler.GetCard).Methods("GET")
	router.HandleFunc("/api/cards/list", cardHandler.ListCards).Methods("GET")
	router.HandleFunc("/api/cards", cardHandler.DeleteCard).Methods("DELETE")

	// HTTP server
	go func() {
		log.Println("HTTP server 8081 portunda başlatılıyor...")
		if err := http.ListenAndServe(":8081", router); err != nil {
			log.Fatalf("HTTP server başlatılamadı: %v", err)
		}
	}()

	// gRPC server'ı başlat
	lis, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Port dinlenemedi: %v", err)
	}

	grpcServer := grpc.NewServer()
	cardpb.RegisterCardServiceServer(grpcServer, cardServer)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("gRPC server 50054 portunda başlatılıyor...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server başlatılamadı: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down...")
	grpcServer.GracefulStop()
}

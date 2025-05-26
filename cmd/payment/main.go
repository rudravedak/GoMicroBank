package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	paymentpb "govo/api/proto/payment"
	"govo/internal/payment/handler"
	"govo/internal/payment/model"
	"govo/internal/payment/repository"
	"govo/internal/payment/service"
	"govo/kafka"

	"github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/timestamppb"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PaymentServer struct {
	paymentpb.UnimplementedPaymentServiceServer
	service *service.PaymentService
}

func (s *PaymentServer) CreatePayment(ctx context.Context, req *paymentpb.CreatePaymentRequest) (*paymentpb.CreatePaymentResponse, error) {
	payment, err := s.service.CreatePayment(
		ctx,
		uint(req.CustomerId),
		uint(req.CardId),
		req.Amount,
		req.PaymentType,
		req.Description,
	)
	if err != nil {
		return nil, err
	}

	return &paymentpb.CreatePaymentResponse{
		Payment: &paymentpb.Payment{
			Id:          uint32(payment.ID),
			CustomerId:  uint32(payment.CustomerID),
			CardId:      uint32(payment.CardID),
			Amount:      payment.Amount,
			PaymentType: payment.PaymentType,
			Status:      payment.Status,
			Description: payment.Description,
			CreatedAt:   timestamppb.New(payment.CreatedAt),
			UpdatedAt:   timestamppb.New(payment.UpdatedAt),
		},
	}, nil
}

func (s *PaymentServer) GetPayment(ctx context.Context, req *paymentpb.GetPaymentRequest) (*paymentpb.GetPaymentResponse, error) {
	payment, err := s.service.GetPayment(ctx, uint(req.PaymentId))
	if err != nil {
		return nil, err
	}

	return &paymentpb.GetPaymentResponse{
		Payment: &paymentpb.Payment{
			Id:          uint32(payment.ID),
			CustomerId:  uint32(payment.CustomerID),
			CardId:      uint32(payment.CardID),
			Amount:      payment.Amount,
			PaymentType: payment.PaymentType,
			Status:      payment.Status,
			Description: payment.Description,
			CreatedAt:   timestamppb.New(payment.CreatedAt),
			UpdatedAt:   timestamppb.New(payment.UpdatedAt),
		},
	}, nil
}

func (s *PaymentServer) ListPayments(ctx context.Context, req *paymentpb.ListPaymentsRequest) (*paymentpb.ListPaymentsResponse, error) {
	var startDate, endDate *time.Time
	if req.StartDate != nil {
		t := req.StartDate.AsTime()
		startDate = &t
	}
	if req.EndDate != nil {
		t := req.EndDate.AsTime()
		endDate = &t
	}

	payments, err := s.service.ListPayments(ctx, uint(req.CustomerId), req.Status, startDate, endDate)
	if err != nil {
		return nil, err
	}

	response := &paymentpb.ListPaymentsResponse{
		Payments: make([]*paymentpb.Payment, len(payments)),
	}

	for i, p := range payments {
		response.Payments[i] = &paymentpb.Payment{
			Id:          uint32(p.ID),
			CustomerId:  uint32(p.CustomerID),
			CardId:      uint32(p.CardID),
			Amount:      p.Amount,
			PaymentType: p.PaymentType,
			Status:      p.Status,
			Description: p.Description,
			CreatedAt:   timestamppb.New(p.CreatedAt),
			UpdatedAt:   timestamppb.New(p.UpdatedAt),
		}
	}

	return response, nil
}

func (s *PaymentServer) CancelPayment(ctx context.Context, req *paymentpb.CancelPaymentRequest) (*paymentpb.CancelPaymentResponse, error) {
	err := s.service.CancelPayment(ctx, uint(req.PaymentId), req.Reason)
	if err != nil {
		return nil, err
	}

	return &paymentpb.CancelPaymentResponse{
		Success: true,
	}, nil
}

func main() {
	// PostgreSQL bağlantısı
	dsn := "host=postgres user=postgres password=postgres dbname=paymentdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Veritabanına bağlanılamadı: %v", err)
	}

	// Tabloları oluştur
	if err := db.AutoMigrate(&model.Payment{}); err != nil {
		log.Fatalf("Tablo oluşturulamadı: %v", err)
	}

	// Kafka client
	kafkaClient := kafka.NewClient([]string{"kafka:9092"})
	defer kafkaClient.Close()

	// Kafka consumer
	consumer := kafka.NewConsumer([]string{"kafka:9092"})
	defer consumer.Close()

	// Consumer'ı başlat
	ctx, cancel := context.WithCancel(context.Background())
	go consumer.Start(ctx)

	// Dependency injection
	paymentRepo := repository.NewPaymentRepository(db)
	paymentService := service.NewPaymentService(paymentRepo, kafkaClient)
	paymentServer := &PaymentServer{service: paymentService}
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// HTTP router
	router := mux.NewRouter()
	router.HandleFunc("/api/payments", paymentHandler.CreatePayment).Methods("POST")
	router.HandleFunc("/api/payments", paymentHandler.GetPayment).Methods("GET")
	router.HandleFunc("/api/payments/list", paymentHandler.ListPayments).Methods("GET")
	router.HandleFunc("/api/payments/cancel", paymentHandler.CancelPayment).Methods("POST")

	// HTTP server
	go func() {
		log.Println("HTTP server 8080 portunda başlatılıyor...")
		if err := http.ListenAndServe(":8080", router); err != nil {
			log.Fatalf("HTTP server başlatılamadı: %v", err)
		}
	}()

	// gRPC server'ı başlat
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Port dinlenemedi: %v", err)
	}

	grpcServer := grpc.NewServer()
	paymentpb.RegisterPaymentServiceServer(grpcServer, paymentServer)

	// Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("gRPC server 50053 portunda başlatılıyor...")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("gRPC server başlatılamadı: %v", err)
		}
	}()

	<-sigChan
	log.Println("Shutting down...")
	cancel()
	grpcServer.GracefulStop()
}

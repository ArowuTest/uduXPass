package database

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/uduxpass/backend/internal/domain/repositories"
	"github.com/uduxpass/backend/internal/infrastructure/database/postgres"
)

type DatabaseManager struct {
	db *sqlx.DB
	
	// Repository instances
	userRepo           repositories.UserRepository
	adminUserRepo      repositories.AdminUserRepository
	orderRepo          repositories.OrderRepository
	orderLineRepo      repositories.OrderLineRepository
	organizerRepo      repositories.OrganizerRepository
	eventRepo          repositories.EventRepository
	ticketTierRepo     repositories.TicketTierRepository
	tourRepo           repositories.TourRepository
	ticketRepo         repositories.TicketRepository
	paymentRepo        repositories.PaymentRepository
	inventoryHoldRepo  repositories.InventoryHoldRepository
	otpTokenRepo       repositories.OTPTokenRepository
	scannerUserRepo    repositories.ScannerUserRepository
}

func NewDatabaseManager(databaseURL string) (*DatabaseManager, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}
	
	// Configure connection pool
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	
	// Test connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return &DatabaseManager{
		db:                db,
		userRepo:          postgres.NewUserRepository(db),
		adminUserRepo:     postgres.NewAdminUserRepository(db),
		orderRepo:         postgres.NewOrderRepository(db),
		orderLineRepo:     postgres.NewOrderLineRepository(db),
		organizerRepo:     postgres.NewOrganizerRepository(db),
		eventRepo:         postgres.NewEventRepository(db),
		ticketTierRepo:    postgres.NewTicketTierRepository(db),
		tourRepo:          postgres.NewTourRepository(db),
		ticketRepo:        postgres.NewTicketRepository(db),
		paymentRepo:       postgres.NewPaymentRepository(db),
		inventoryHoldRepo: postgres.NewInventoryHoldRepository(db),
		otpTokenRepo:      postgres.NewOTPTokenRepository(db),
		scannerUserRepo:   postgres.NewScannerUserRepository(db),
	}, nil
}

func (dm *DatabaseManager) Close() error {
	return dm.db.Close()
}

func (dm *DatabaseManager) Health(ctx context.Context) error {
	return dm.db.PingContext(ctx)
}

// Repository accessors
func (dm *DatabaseManager) Users() repositories.UserRepository {
	return dm.userRepo
}

func (dm *DatabaseManager) AdminUsers() repositories.AdminUserRepository {
	return dm.adminUserRepo
}

func (dm *DatabaseManager) Orders() repositories.OrderRepository {
	return dm.orderRepo
}

func (dm *DatabaseManager) OrderLines() repositories.OrderLineRepository {
	return dm.orderLineRepo
}

func (dm *DatabaseManager) Organizers() repositories.OrganizerRepository {
	return dm.organizerRepo
}

func (dm *DatabaseManager) Events() repositories.EventRepository {
	return dm.eventRepo
}

func (dm *DatabaseManager) TicketTiers() repositories.TicketTierRepository {
	return dm.ticketTierRepo
}

func (dm *DatabaseManager) Tours() repositories.TourRepository {
	return dm.tourRepo
}

func (dm *DatabaseManager) Tickets() repositories.TicketRepository {
	return dm.ticketRepo
}

func (dm *DatabaseManager) Payments() repositories.PaymentRepository {
	return dm.paymentRepo
}

func (dm *DatabaseManager) InventoryHolds() repositories.InventoryHoldRepository {
	return dm.inventoryHoldRepo
}

func (dm *DatabaseManager) OTPTokens() repositories.OTPTokenRepository {
	return dm.otpTokenRepo
}

// Transaction support
func (dm *DatabaseManager) BeginTx(ctx context.Context) (*sqlx.Tx, error) {
	return dm.db.BeginTxx(ctx, nil)
}

// Note: Transaction support would require repository interface changes
// For now, repositories work with the main DB connection


func (dm *DatabaseManager) ScannerUsers() repositories.ScannerUserRepository {
	return dm.scannerUserRepo
}


// UnitOfWork returns the unit of work for transaction management
func (dm *DatabaseManager) UnitOfWork() repositories.UnitOfWork {
	return postgres.NewUnitOfWork(dm.db)
}

// GetDB returns the database connection for direct queries
func (dm *DatabaseManager) GetDB() *sqlx.DB {
	return dm.db
}


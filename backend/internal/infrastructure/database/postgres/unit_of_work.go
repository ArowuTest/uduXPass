package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/uduxpass/backend/internal/domain/repositories"
)

// postgresTransaction implements the Transaction interface
type postgresTransaction struct {
	tx  *sqlx.Tx
	ctx context.Context
	
	// Repository instances
	organizers      repositories.OrganizerRepository
	tours           repositories.TourRepository
	events          repositories.EventRepository
	ticketTiers     repositories.TicketTierRepository
	users           repositories.UserRepository
	orders          repositories.OrderRepository
	orderLines      repositories.OrderLineRepository
	tickets         repositories.TicketRepository
	payments        repositories.PaymentRepository
	inventoryHolds  repositories.InventoryHoldRepository
	adminUsers      repositories.AdminUserRepository
	scannerUsers    repositories.ScannerUserRepository
	otpTokens       repositories.OTPTokenRepository
}

// Commit commits the transaction
func (t *postgresTransaction) Commit() error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	return t.tx.Commit()
}

// Rollback rolls back the transaction
func (t *postgresTransaction) Rollback() error {
	if t.tx == nil {
		return fmt.Errorf("transaction is nil")
	}
	return t.tx.Rollback()
}

// Context returns the transaction context
func (t *postgresTransaction) Context() context.Context {
	return t.ctx
}

// Organizers returns the organizer repository within this transaction
func (t *postgresTransaction) Organizers() repositories.OrganizerRepository {
	if t.organizers == nil {
		t.organizers = &organizerRepository{db: t.tx}
	}
	return t.organizers
}

// Tours returns the tour repository within this transaction
func (t *postgresTransaction) Tours() repositories.TourRepository {
	if t.tours == nil {
		t.tours = &tourRepository{db: t.tx}
	}
	return t.tours
}

// Events returns the event repository within this transaction
func (t *postgresTransaction) Events() repositories.EventRepository {
	if t.events == nil {
		t.events = NewEventRepositoryWithTx(t.tx)
	}
	return t.events
}

// TicketTiers returns the ticket tier repository within this transaction
func (t *postgresTransaction) TicketTiers() repositories.TicketTierRepository {
	if t.ticketTiers == nil {
		t.ticketTiers = NewTicketTierRepositoryWithTx(t.tx)
	}
	return t.ticketTiers
}

// Users returns the user repository within this transaction
func (t *postgresTransaction) Users() repositories.UserRepository {
	if t.users == nil {
		t.users = NewUserRepositoryWithTx(t.tx)
	}
	return t.users
}

// Orders returns the order repository within this transaction
func (t *postgresTransaction) Orders() repositories.OrderRepository {
	if t.orders == nil {
		t.orders = NewOrderRepositoryWithTx(t.tx)
	}
	return t.orders
}

// OrderLines returns the order line repository within this transaction
func (t *postgresTransaction) OrderLines() repositories.OrderLineRepository {
	if t.orderLines == nil {
		t.orderLines = &orderLineRepository{db: t.tx}
	}
	return t.orderLines
}

// Tickets returns the ticket repository within this transaction
func (t *postgresTransaction) Tickets() repositories.TicketRepository {
	if t.tickets == nil {
		t.tickets = &ticketRepository{db: t.tx}
	}
	return t.tickets
}

// Payments returns the payment repository within this transaction
func (t *postgresTransaction) Payments() repositories.PaymentRepository {
	if t.payments == nil {
		t.payments = &paymentRepository{db: t.tx}
	}
	return t.payments
}

// InventoryHolds returns the inventory hold repository within this transaction
func (t *postgresTransaction) InventoryHolds() repositories.InventoryHoldRepository {
	if t.inventoryHolds == nil {
		t.inventoryHolds = &inventoryHoldRepository{db: t.tx}
	}
	return t.inventoryHolds
}

// AdminUsers returns the admin user repository within this transaction
func (t *postgresTransaction) AdminUsers() repositories.AdminUserRepository {
	if t.adminUsers == nil {
		t.adminUsers = NewAdminUserRepositoryWithTx(t.tx)
	}
	return t.adminUsers
}

// ScannerUsers returns the scanner user repository within this transaction
func (t *postgresTransaction) ScannerUsers() repositories.ScannerUserRepository {
	if t.scannerUsers == nil {
		t.scannerUsers = NewScannerUserRepositoryWithTx(t.tx)
	}
	return t.scannerUsers
}

// OTPTokens returns the OTP token repository within this transaction
func (t *postgresTransaction) OTPTokens() repositories.OTPTokenRepository {
	if t.otpTokens == nil {
		t.otpTokens = &otpTokenRepository{db: t.tx}
	}
	return t.otpTokens
}

// postgresUnitOfWork implements the UnitOfWork interface
type postgresUnitOfWork struct {
	db *sqlx.DB
}

// NewUnitOfWork creates a new PostgreSQL UnitOfWork
func NewUnitOfWork(db *sqlx.DB) repositories.UnitOfWork {
	return &postgresUnitOfWork{db: db}
}

// Begin starts a new transaction
func (uow *postgresUnitOfWork) Begin(ctx context.Context) (repositories.Transaction, error) {
	tx, err := uow.db.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelReadCommitted,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	
	return &postgresTransaction{
		tx:  tx,
		ctx: ctx,
	}, nil
}

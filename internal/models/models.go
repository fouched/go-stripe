package models

import (
	"context"
	"database/sql"
	"time"
)

const dbTimeout = 3 * time.Second

// DBModel is the type for database connection values
type DBModel struct {
	DB *sql.DB
}

// Models is the wrapper for all models
type Models struct {
	DB DBModel
}

// NewModels return a model type with database connection pool
func NewModels(db *sql.DB) Models {
	return Models{
		DB: DBModel{DB: db},
	}
}

// Widget is the type for widgets
type Widget struct {
	ID             int       `json:"ID"`
	Name           string    `json:"Name"`
	Description    string    `json:"Description"`
	InventoryLevel int       `json:"InventoryLevel"`
	Price          int       `json:"Price"`
	Image          string    `json:"Image"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// Order is the type for orders
type Order struct {
	ID            int       `json:"ID"`
	WidgetID      int       `json:"WidgetID"`
	TransactionID int       `json:"TransactionID"`
	CustomerID    int       `json:"CustomerID"`
	StatusID      int       `json:"StatusID"`
	Quantity      int       `json:"Quantity"`
	Amount        int       `json:"Amount"`
	CreatedAt     time.Time `json:"-"`
	UpdatedAt     time.Time `json:"-"`
}

// Status is the type for order statuses
type Status struct {
	ID        int       `json:"ID"`
	Name      string    `json:"Name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// TransactionStatus is the type for transaction statuses
type TransactionStatus struct {
	ID        int       `json:"ID"`
	Name      string    `json:"Name"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Transaction is the type for transactions
type Transaction struct {
	ID                  int       `json:"ID"`
	Amount              int       `json:"Amount"`
	Currency            string    `json:"Currency"`
	LastFour            string    `json:"LastFour"`
	PaymentMethod       string    `json:"PaymentMethod"`
	PaymentIntent       string    `json:"PaymentIntent"`
	ExpiryMonth         int       `json:"ExpiryMonth"`
	ExpiryYear          int       `json:"ExpiryYear"`
	BankReturnCode      string    `json:"BankReturnCode"`
	TransactionStatusID int       `json:"TransactionStatusID"`
	CreatedAt           time.Time `json:"-"`
	UpdatedAt           time.Time `json:"-"`
}

// User is the type for users
type User struct {
	ID        int       `json:"ID"`
	FirstName string    `json:"FirstName"`
	LastName  string    `json:"LastName"`
	Email     string    `json:"Email"`
	Password  string    `json:"Password"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

// Customer is the type for users
type Customer struct {
	ID        int       `json:"ID"`
	FirstName string    `json:"FirstName"`
	LastName  string    `json:"LastName"`
	Email     string    `json:"Email"`
	CreatedAt time.Time `json:"-"`
	UpdatedAt time.Time `json:"-"`
}

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var widget Widget

	row := m.DB.QueryRowContext(ctx, `
			select 
				id, name, description, inventory_level, price, image, 
				created_at, updated_at 
			from 
				widgets 
			where id = ?`, id)
	err := row.Scan(
		&widget.ID,
		&widget.Name,
		&widget.Description,
		&widget.InventoryLevel,
		&widget.Price,
		&widget.Image,
		&widget.CreatedAt,
		&widget.UpdatedAt,
	)
	if err != nil {
		return widget, err
	}

	return widget, nil
}

// InsertTransaction insert a new transaction and returns its id
func (m *DBModel) InsertTransaction(txn Transaction) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		insert into transactions 
			(amount, currency, last_four, expiry_month, expiry_year, payment_method, 
			 payment_intent, bank_return_code, transaction_status_id, created_at, updated_at) 
		values(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.ExecContext(ctx, stmt,
		txn.Amount,
		txn.Currency,
		txn.LastFour,
		txn.ExpiryMonth,
		txn.ExpiryYear,
		txn.PaymentMethod,
		txn.PaymentIntent,
		txn.BankReturnCode,
		txn.TransactionStatusID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertOrder insert a new order and returns its id
func (m *DBModel) InsertOrder(order Order) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		insert into orders 
			(widget_id, transaction_id, status_id, customer_id, quantity, amount, created_at, updated_at) 
		values(?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := m.DB.ExecContext(ctx, stmt,
		order.WidgetID,
		order.TransactionID,
		order.StatusID,
		order.CustomerID,
		order.Quantity,
		order.Amount,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// InsertCustomer insert a new order and returns its id
func (m *DBModel) InsertCustomer(c Customer) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	stmt := `
		insert into customers 
			(first_name, last_name, email, created_at, updated_at) 
		values(?, ?, ?, ?, ?)`

	result, err := m.DB.ExecContext(ctx, stmt,
		c.FirstName,
		c.LastName,
		c.Email,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

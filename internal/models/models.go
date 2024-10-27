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
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Description    string    `json:"description"`
	InventoryLevel int       `json:"inventory_level"`
	Price          int       `json:"price"`
	Image          string    `json:"image"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// Order is the type for orders
type Order struct {
	ID            int       `json:"ID"`
	WidgetID      int       `json:"WidgetID"`
	TransactionID int       `json:"TransactionID"`
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

func (m *DBModel) GetWidget(id int) (Widget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var widget Widget

	row := m.DB.QueryRowContext(ctx, "select id, name from widgets where id = ?", id)
	err := row.Scan(&widget.ID, &widget.Name)
	if err != nil {
		return widget, err
	}

	return widget, nil
}

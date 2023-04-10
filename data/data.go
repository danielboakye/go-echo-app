package data

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"golang.org/x/crypto/bcrypt"
)

const dbTimeout = time.Second * 3

var psql sq.StatementBuilderType

type Repository struct {
	db *sql.DB
}

type IRepository interface {
	GetAll() ([]*User, error)
	GetOne(string) (*User, error)
	GetByEmail(string) (*User, error)
	Update(User) error
	DeleteByID(string) error
	Insert(User) (string, error)
}

func NewRepository(pool *sql.DB) IRepository {
	psql = sq.StatementBuilder.PlaceholderFormat(sq.Dollar)
	return &Repository{db: pool}
}

type User struct {
	ID        string    `json:"user_id"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Password  string    `json:"password,omitempty"`
	Active    int       `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// GetAll returns a slice of all users, sorted by last name
func (r *Repository) GetAll() ([]*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	uq := psql.Select("user_id, email, first_name, last_name, user_active, created_at, updated_at").
		From("users").
		OrderBy("last_name ASC")
	rows, err := uq.RunWith(r.db).QueryContext(ctx)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		var user User
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Active,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		users = append(users, &user)
	}

	return users, nil
}

// GetOne returns one user by id
func (r *Repository) GetOne(id string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user User
	uq := psql.Select("user_id, email, first_name, last_name, password, user_active, created_at, updated_at").
		From("users").
		Where(sq.Eq{"user_id": id})
	row := uq.RunWith(r.db).QueryRowContext(ctx)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail returns one user by email
func (r *Repository) GetByEmail(email string) (*User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var user User
	uq := psql.Select("user_id, email, first_name, last_name, password, user_active, created_at, updated_at").
		From("users").
		Where(sq.Eq{"email": email})
	row := uq.RunWith(r.db).QueryRowContext(ctx)
	err := row.Scan(
		&user.ID,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Password,
		&user.Active,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Update updates one user in the database, using the information
// stored in the receiver u
func (r *Repository) Update(u User) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := psql.Update("users").
		SetMap(
			sq.Eq{
				"email": u.Email, "first_name": u.FirstName,
				"last_name": u.LastName, "user_active": u.Active,
				"updated_at": time.Now(),
			}).
		Where(sq.Eq{"user_id": u.ID}).
		RunWith(r.db).ExecContext(ctx)

	return err
}

// DeleteByID deletes one user from the database, by ID
func (r *Repository) DeleteByID(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	_, err := psql.Delete("users").Where(sq.Eq{"user_id": id}).
		RunWith(r.db).ExecContext(ctx)

	return err
}

// Insert inserts a new user into the database, and returns the ID of the newly inserted row
func (r *Repository) Insert(u User) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var newID string
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), 12)
	if err != nil {
		return newID, err
	}

	uq := psql.Insert("users").
		Columns("email", "first_name", "last_name", "password", "user_active", "created_at", "updated_at").
		Values(u.Email, u.FirstName, u.LastName, hashedPassword, u.Active, time.Now(), time.Now()).
		Suffix("RETURNING user_id").
		RunWith(r.db).QueryRowContext(ctx)

	err = uq.Scan(&newID)
	if err != nil {
		return newID, err
	}

	return newID, nil
}

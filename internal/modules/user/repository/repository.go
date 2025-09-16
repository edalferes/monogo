package repository

import (
	"database/sql"

	"github.com/edalferes/monogo/internal/modules/user/domain"
)

type UserRepository interface {
	CreateUser(name string, email string) (int64, error)
	FindByID(id int64) (*domain.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(name string, email string) (int64, error) {
	result, err := r.db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", name, email)
	if err != nil {
		return 0, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *userRepository) FindByID(id int64) (*domain.User, error) {
	row := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id)
	var user domain.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		return nil, err
	}
	return &user, nil
}

var _ UserRepository = (*userRepository)(nil)

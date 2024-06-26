package repository

import (
	"database/sql"
	"fmt"
	"sync"

	"github.com/sraynitjsr/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (repo *UserRepository) GetAll() ([]model.User, error) {
	rows, err := repo.db.Query("SELECT id, name, age FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Channel to collect users from goroutines
	userChan := make(chan model.User)
	defer close(userChan)

	var wg sync.WaitGroup
	var mu sync.Mutex
	users := make([]model.User, 0)

	// Launch goroutine for each row
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Age)
		if err != nil {
			return nil, err
		}

		wg.Add(1)
		go func(u model.User) {
			defer wg.Done()
			userChan <- u
		}(user)
	}

	// Close channel after all goroutines finish
	go func() {
		wg.Wait()
		close(userChan)
	}()

	// Collect users from goroutines
	for user := range userChan {
		mu.Lock()
		users = append(users, user)
		mu.Unlock()
	}

	return users, nil
}

func (repo *UserRepository) GetByID(id int) (model.User, error) {
	var user model.User
	err := repo.db.QueryRow("SELECT id, name, age FROM users WHERE id=?", id).Scan(&user.ID, &user.Name, &user.Age)
	if err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (repo *UserRepository) Create(user model.User) (model.User, error) {
	result, err := repo.db.Exec("INSERT INTO users (name, age) VALUES (?, ?)", user.Name, user.Age)
	if err != nil {
		return model.User{}, err
	}
	lastInsertID, err := result.LastInsertId()
	if err != nil {
		return model.User{}, err
	}
	user.ID = int(lastInsertID)
	return user, nil
}

func (repo *UserRepository) Update(id int, user model.User) (model.User, error) {
	_, err := repo.db.Exec("UPDATE users SET name = ?, age = ? WHERE id = ?", user.Name, user.Age, id)
	if err != nil {
		return model.User{}, err
	}
	user.ID = id
	return user, nil
}

func (repo *UserRepository) Delete(id int) error {
	_, err := repo.db.Exec("DELETE FROM users WHERE id = ?", id)
	if err != nil {
		return err
	}
	return nil
}

func ConnectToDB(dataSourceName string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging database: %v", err)
	}
	return db, nil
}

package user

import (
	"github.com/google_oauth/db"
	"github.com/google_oauth/entity"
)

// Service procides user's behavior
type Service struct{}

// User is alias of entity.User struct
type User entity.User

// GetAll is get all User
func (s Service) GetAll() ([]User, error) {
	// DB connect
	db := db.GetDB()
	var u []User

	if err := db.Find(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (s Service) MyFunc(user *User) (User, error) {
	db := db.GetDB()
	var u User
	u.Email = user.Email
	u.Password = user.Password

	if err := db.Create(&u).Error; err != nil {
		return u, err
	}

	return u, nil
}

// CreateModel is create User model
// func (s Service) CreateModel(c *gin.Context) (User, error) {
// 	db := db.GetDB()
// 	var u User

// 	if err := c.BindJSON(&u); err != nil {
// 		return u, err
// 	}

// 	if err := db.Create(&u).Error; err != nil {
// 		return u, err
// 	}

// 	return u, nil
// }

// CreateUser is create User
func (s Service) CreateUser(email string, password string) (User, error) {
	// func (s Service) CreateUser(u *User) (User, error) {
	db := db.GetDB()
	// var user User
	var user User

	user.Email = email
	user.Password = password
	//user.Description = "It`s Twitter username"
	//db.Save(&user)
	if err := db.Create(&user).Error; err != nil {
		return user, err
	}
	return user, nil
}

// GetByID is get a User
func (s Service) GetByID(id string) (User, error) {
	db := db.GetDB()
	var u User

	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		return u, err
	}

	return u, nil
}

// GetByEmail&Password
func (s Service) GetByEmailAndPassword(email string, password string) (User, error) {
	db := db.GetDB()
	var u User

	if err := db.Where(" email = ? AND password = ?", email, password).First(&u).Error; err != nil {
		return u, err
	}

	return u, nil
}

// GetByPassword
func (s Service) GetByPassword(email string) (User, error) {
	db := db.GetDB()
	var u User

	if err := db.Where(" email = ?", email).First(&u).Error; err != nil {
		return u, err
	}

	return u, nil
}

// UpdateByID is update a User
// func (s Service) UpdateByID(id string, c *gin.Context) (User, error) {
// 	db := db.GetDB()
// 	var u User

// 	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
// 		return u, err
// 	}

// 	if err := c.BindJSON(&u); err != nil {
// 		return u, err
// 	}

// 	db.Save(&u)

// 	return u, nil
// }

// DeleteByID is delete a User
func (s Service) DeleteByID(id string) error {
	db := db.GetDB()
	var u User

	if err := db.Where("id = ?", id).Delete(&u).Error; err != nil {
		return err
	}

	return nil
}

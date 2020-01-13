package user

import (
	"github.com/google_oauth/db"
	"github.com/google_oauth/entity"
)

// Service procides user's behavior
type Service struct{}

// GetAll is get all User
func (s Service) GetAll() ([]entity.User, error) {
	// DB connect
	db := db.GetDB()
	var u []entity.User

	if err := db.Find(&u).Error; err != nil {
		return nil, err
	}

	return u, nil
}

func (s Service) MyFunc(user *entity.User) (*entity.User, error) {
	db := db.GetDB()

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
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
func (s Service) CreateUser(user *entity.User) (*entity.User, error) {
	// func (s Service) CreateUser(u *User) (User, error) {
	db := db.GetDB()
	//user.Description = "It`s Twitter username"
	//db.Save(&user)
	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

// GetByID is get a User
func (s Service) GetByID(id string) (*entity.User, error) {
	db := db.GetDB()
	var u entity.User

	if err := db.Where("id = ?", id).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

// GetByEmail&Password
func (s Service) GetByEmailAndPassword(email string, password string) (*entity.User, error) {
	db := db.GetDB()
	var u entity.User

	if err := db.Where(" email = ? AND password = ?", email, password).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
}

// GetByPassword
func (s Service) GetByPassword(email string) (*entity.User, error) {
	db := db.GetDB()
	var u entity.User

	if err := db.Where(" email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}

	return &u, nil
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
	var u entity.User

	if err := db.Where("id = ?", id).Delete(&u).Error; err != nil {
		return err
	}

	return nil
}

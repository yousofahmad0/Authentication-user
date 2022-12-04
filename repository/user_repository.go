package repository

import (
	"authentication-user/entity"
	"fmt"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

// GetAll students
func (s UserRepository) GetAll() []*entity.User {
	var users []*entity.User
	//_ = s.DB.Find(&students)
	if err := s.DB.Find(&users).Error; err != nil {
		return nil
	}
	return users
}

// GetOne user
func (s UserRepository) GetStudent(id int) *entity.User {
	user := new(entity.User)
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil
	}
	return user
}

// Create user
func (s UserRepository) Create(u *entity.User) *entity.User {

	result := s.DB.Create(&u)
	fmt.Println(result.Row())
	s.DB.Save(&u)
	return u
}

// Delete user
func (s UserRepository) Delete(id int) *entity.User {

	user := new(entity.User)
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil
	}
	deleted := s.DB.Delete(&entity.User{}, id)
	if deleted.RowsAffected < 1 {
		return nil
	}

	return user
}

// Update user
func (s UserRepository) Update(id int, u *entity.User) *entity.User {

	user := new(entity.User)
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil
	}
	if u.FirstName == "" || u.LastName == "" {
		return nil
	}

	user.FirstName = u.FirstName
	user.LastName = u.LastName
	s.DB.Save(user)
	return user
}

// Patch Update user
func (s UserRepository) Patch(id int, u *entity.User) *entity.User {

	user := new(entity.User)
	if err := s.DB.First(&user, id).Error; err != nil {
		return nil
	}

	if u.FirstName != "" {
		user.FirstName = u.FirstName
	}
	if u.LastName != "" {
		user.LastName = u.LastName
	}
	s.DB.Save(user)
	return user
}

package models

import (
	"fintech/utils/token"

	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID        uuid.UUID `gorm:"primary_key; unique" json:"id"`
	Email     string    `gorm:"type:varchar(255);unique;not null" json:"email" validate:"omitempty,email,min=3,max=255"`
	Password  string    `gorm:"type:varchar(255);not null" json:"password" validate:"omitempty,min=6,max=255"`
	Role      string    `gorm:"type:varchar(255);not null" json:"role"` // 3 roles : admin, investor, investee
	Status    string    `gorm:"type:varchar(255);not null;default:aktif" json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relation
	UserProfile  UserProfile `gorm:"foreignKey:ID"`
	Transactions []Transaction
	Projects     []Projects
	ShoppingCart []ShoppingCart
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func LoginCheck(email, password string, db *gorm.DB) (string, string, string, string, error) {
	var err error

	u := User{}

	err = db.Model(User{}).Where("email = ?", email).Take(&u).Error

	if err != nil {
		return "", "", "", "", err
	}
	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return "", "", "", "", err
	}

	token, err := token.GenerateToken(u.ID, u.Role, u.Status)

	if err != nil {
		return "", "", "", "", err
	}
	return token, u.Role, u.Status, u.ID.String(), nil
}

func PasswordCheck(UUID, password string, db *gorm.DB) error {
	var err error

	u := User{}

	err = db.Model(User{}).Where("id = ?", UUID).Take(&u).Error

	if err != nil {
		return err
	}

	err = VerifyPassword(password, u.Password)

	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		return err
	}

	return nil
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	// Turn password into hash
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if errPassword != nil {
		return &User{}, errPassword
	}
	u.Password = string(hashedPassword)
	// remove space in email
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.ID = uuid.New()
	// perlu user Role?
	var err error = db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	profil := UserProfile{
		ID: u.ID,
	}

	err = db.Create(&profil).Error
	if err != nil {
		return &User{}, err
	}

	// ketika register otomatis membuat shopping cart dengan status "order"
	shoppingCart := ShoppingCart{
		ID:            uuid.New(),
		PaymentStatus: "order",
		UserID:        u.ID,
	}

	err = db.Create(&shoppingCart).Error
	if err != nil {
		return &User{}, err
	}

	return u, nil
}

func (u *User) SavePassword() error {
	hashedPassword, errPassword := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if errPassword != nil {
		return errPassword
	}
	u.Password = string(hashedPassword)
	return nil
}

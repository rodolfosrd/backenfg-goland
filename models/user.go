package models

import (
	"fmt"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/twinj/uuid" //jwt-best-practices
	"golang.org/x/crypto/bcrypt"

	//"github.com/gofrs/uuid" //offersapp
	"time"

	"gorm.io/gorm"
)

var (
	tokenSecret = []byte(os.Getenv("TOKEN_SECRET"))
)

//https://gorm.io/docs/conventions.html
//type Tabler interface {
//TableName() string
//}

// TableName overrides the table name used by Empleado to `employee`
func (User) TableName() string {
	return "user_account"
}

// BeforeCreate will set a UUID rather than numeric ID. https://gorm.io/docs/create.html
func (tab *User) BeforeCreate(*gorm.DB) error {
	//uuidx := uuid.NewV4()
	tab.ID = uuid.NewV4().String()
	return nil
}

type User struct {
	ID string `gorm:"primary_key;column:id" json:"id"` //json:"id,omitempty"
	//ID              uuid.UUID `json:"id"`
	CreatedAt       time.Time `json:"-"`
	UpdatedAt       time.Time `json:"-"`
	Email           string    `gorm:"column:email" json:"email"`
	PasswordHash    string    `json:"-"`
	Password        string    `json:"password"`
	PasswordConfirm string    `json:"password_confirm"`
}

func (u *User) Register(conn *gorm.DB) error {

	if len(u.Password) < 4 || len(u.PasswordConfirm) < 4 {
		return fmt.Errorf("Password must be at least 4 characters long.")
	}

	if u.Password != u.PasswordConfirm {
		return fmt.Errorf("Passwords do not match.")
	}

	if len(u.Email) < 4 {
		return fmt.Errorf("Email must be at least 4 characters long.")
	}

	u.Email = strings.ToLower(u.Email)
	var userLookup User
	var err error
	if err = conn.First(&userLookup, "email = ?", u.Email).Error; err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		//return fmt.Errorf("Error p" + err.Error())
	}

	//row := conn.QueryRow(context.Background(), "SELECT id from user_account WHERE email = $1", u.Email)
	//userLookup := User{}
	//err := row.Scan(&userLookup)
	if u.Email == strings.ToLower(userLookup.Email) {
		fmt.Println("found user")
		fmt.Println(userLookup.Email)
		return fmt.Errorf("A user with that email already exists")
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("There was an error creating your account.")
	}
	u.PasswordHash = string(pwdHash)

	//now := time.Now()
	//_, err = conn.Exec(context.Background(), "INSERT INTO user_account (created_at, updated_at, email, password_hash) VALUES($1, $2, $3, $4)", now, now, u.Email, u.PasswordHash)

	//row2 := conn.QueryRow(context.Background(), "SELECT id, password_hash from user_account WHERE email = $1", u.Email)
	//row2.Scan(&u.ID, &u.PasswordHash)
	//if err := c.BindJSON(&userLookup); err != nil {

	//	return fmt.Errorf(err.Error())
	//}

	u.Password = ""
	u.PasswordConfirm = ""
	conn.Create(&u)
	return err // ya te asigna el ID del user
}

// GetAuthToken returns the auth token to be used
func (u *User) GetAuthToken() (string, error) { //CreateToken
	claims := jwt.MapClaims{}
	claims["authorized"] = true
	claims["mio"] = "hola"
	claims["user_id"] = u.ID
	claims["exp"] = time.Now().Add(time.Minute * 10).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)
	authToken, err := token.SignedString(tokenSecret)
	return authToken, err
}

func DelTokenValid(tokenString string) bool {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// fmt.Printf("Parsing: %v \n", token)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			return nil, fmt.Errorf("Token signing method is not valid: %v", token.Header["alg"])
		}

		return tokenSecret, nil
	})

	if err != nil {
		fmt.Printf("Err %v \n", err)
		return false
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		claims["authorized"] = false
		claims["user_id"] = nil
		claims["exp"] = time.Now().Add(time.Hour * 0).Unix()
		fmt.Println(claims)
		//userID := claims["user_id"]
		return true
	} else {
		fmt.Printf("The alg header %v \n", claims["alg"])
		fmt.Println(err)
		return false
	}
}

func IsTokenValid(tokenString string) (bool, string) { //VerifyToken and TokenValid
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// fmt.Printf("Parsing: %v \n", token)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok == false {
			return nil, fmt.Errorf("Token signing method is not valid: %v", token.Header["alg"])
		}

		return tokenSecret, nil
	})

	if err != nil {
		fmt.Printf("Err %v \n", err)
		return false, ""
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
		userID := claims["user_id"]
		return true, userID.(string)
	} else {
		fmt.Printf("The alg header %v \n", claims["alg"])
		fmt.Println(err)
		return false, "uuid.UUID{}"
	}
}

// IsAuthenticated checks to make sure password is correct and user is active
func (u *User) IsAuthenticated(conn *gorm.DB) error {
	//row := conn.QueryRow(context.Background(), "SELECT id, password_hash from user_account WHERE email = $1", u.Email)
	//err := row.Scan(&u.ID, &u.PasswordHash)
	//if err == pgx.ErrNoRows {
	//	fmt.Println("User with email not found")
	//	return fmt.Errorf("Invalid login credentials")
	//}

	u.Email = strings.ToLower(u.Email)
	var userLookup User
	var err error
	if err = conn.First(&userLookup, "email = ?", u.Email).Error; err != nil {
		//return fmt.Errorf("Error p" + err.Error())
		fmt.Println("User with email not found")
		return fmt.Errorf("Invalid login credentials email")
	}
	//if u.Email == strings.ToLower(userLookup.Email) {
	//	fmt.Println("found user")
	//	fmt.Println(userLookup.Email)
	//	return fmt.Errorf("A user with that email already exists")
	//}

	//pwdHash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	//if err != nil {
	//	return fmt.Errorf("There was an error creating your account.")
	//}
	u.PasswordHash = string(userLookup.PasswordHash)
	u.ID = userLookup.ID

	err = bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(u.Password))
	if err != nil {
		//fmt.Println("pss: " + err.Error())
		return fmt.Errorf("Invalid login credentials pass")
	}

	return nil
}

package server

import (
	"context"
	"fmt"
	"go_ant_work/internal/database"
	"go_ant_work/internal/structs"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

const expired = 60 * 60 * 1000
const secretKey = "iloveu"

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"pasword"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	AcccountId string `json:"accountId"`
	Token      string `json:"token"`
	Name       string `json:"Name"`
	Avatar     string `json:"Avatar"`
}

func (s *Server) Register(ctx context.Context, registerRequest RegisterRequest) (string, error) {
	account, _ := s.db.FindAccountByUsername(registerRequest.Username)
	fmt.Println(account)
	if account != nil {
		return "", fmt.Errorf("USERNAME_EXITS")
	}

	accountId := uuid.NewString()
	password, err := HashPassword(registerRequest.Password)

	if err != nil {
		return "", err
	}

	accountSave := &database.Account{
		Id:        accountId,
		Username:  registerRequest.Username,
		Password:  password,
		Role:      "USER",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = s.db.CreateAccount(accountSave)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	userId := uuid.NewString()
	userSave := &database.User{
		Id:        userId,
		Name:      registerRequest.Name,
		Email:     registerRequest.Email,
		Avatar:    registerRequest.Avatar,
		AccountId: accountId,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = s.db.CreateUser(userSave)
	if err != nil {
		return "", err
	}

	return accountId, nil
}

func (s *Server) Login(ctx context.Context, loginRequest *LoginRequest) (*LoginResponse, error) {
	account, err := s.db.FindAccountByUsername(loginRequest.Username)

	if err != nil {
		return nil, err
	}

	result, err := s.redis.HashGet("retry_password", account.Id)

	if err != nil {
		return nil, err
	}

	if result != nil {
		retryCount, err := strconv.ParseInt(result.(string), 10, 64)
		if err != nil {
			return nil, err
		}

		if retryCount >= 5 {
			return nil, err
		}
	}

	check := CheckPasswordHash(loginRequest.Password, account.Password)
	fmt.Println("check password", check)
	if check {
		user, err := s.db.FindUserByAccountId(account.Id)
		if err != nil {
			return nil, err
		}

		token, err := GennerateJwt(account.Id, user.Id)
		if err != nil {
			return nil, err
		}

		return &LoginResponse{
			Token:      token,
			AcccountId: account.Id,
			Name:       user.Name,
			Avatar:     user.Avatar,
		}, nil
	} else {
		result, err := s.redis.HashExists("retry_password", account.Id)
		if err != nil {
			return nil, err
		}
		if result {
			s.redis.HashIncrement("retry_password", account.Id)
		} else {
			s.redis.HashSet("retry_password", account.Id, 1, time.Minute*30)
		}
		return nil, fmt.Errorf("Password wrong")
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func CheckPasswordHash(password string, passwordHash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(passwordHash), []byte(password))
	fmt.Println(err)
	return err == nil
}

func GennerateJwt(accountId string, userID string) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"accountId": accountId,                   // Subject (user identifier)
		"exp":       time.Now().Unix() + expired, // Expiration time
		"iat":       time.Now().Unix(),           // Issued at
	})

	token, err := claims.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %v", err)
	}

	return token, nil
}

func DecryptToken(tokenString string) (structs.TokenData, error) {
	var tokenData structs.TokenData
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretKey), nil
	})

	if err != nil {
		return tokenData, fmt.Errorf("failed to parse token: %v", err)
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if accountID, ok := claims["accountId"].(string); ok {
			tokenData.AcccountID = accountID
			fmt.Println(tokenData.AcccountID)
		} else {
			return tokenData, fmt.Errorf("accountId is not a string")
		}

		if userID, ok := claims["userID"].(string); ok {
			tokenData.UserID = userID
			fmt.Println(tokenData.UserID)
		} else if uidFloat, ok := claims["userID"].(float64); ok {
			tokenData.UserID = fmt.Sprintf("%.0f", uidFloat) // convert float64 to string
			fmt.Println(tokenData.UserID)
		} else {
			return tokenData, fmt.Errorf("userID is not a string or float64")
		}

		if expFloat, ok := claims["exp"].(float64); ok {
			tokenData.Exp = fmt.Sprintf("%.0f", expFloat) // convert float64 to string
			fmt.Println(tokenData.Exp)
		} else {
			return tokenData, fmt.Errorf("exp is not a float64")
		}

		if iatFloat, ok := claims["iat"].(float64); ok {
			tokenData.Iat = fmt.Sprintf("%.0f", iatFloat) // convert float64 to string
			fmt.Println(tokenData.Iat)
		} else {
			return tokenData, fmt.Errorf("iat is not a float64")
		}

		fmt.Println(tokenData)
		return tokenData, nil
	} else {
		return tokenData, fmt.Errorf("invalid token")
	}
}

func VerifyToken(tokenString string) (*jwt.MapClaims, error) {
	secretKey := []byte(secretKey)

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Kiểm tra phương thức ký có đúng không
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Kiểm tra xem token có hợp lệ và claims là jwt.MapClaims không
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"exchange-tutorials/db"
	"exchange-tutorials/types"
	"exchange-tutorials/utils"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"math/big"
	"time"
)

type AuthService struct {
	secret []byte
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{secret: []byte(secret)}
}

func (as *AuthService) Register(email, password string) (string, error) {
	hashed := sha256.Sum256([]byte(password))
	user := &types.User{
		ID:        utils.GenerateID(),
		Email:     email,
		Password:  hex.EncodeToString(hashed[:]),
		Balance:   make(map[string]*big.Float),
		Margin:    make(map[string]*big.Float),
		Positions: make(map[string]*types.Position),
		Mode:      "isolated",
	}
	// 存储到数据库 (省略)
	return user.ID, nil
}

func (as *AuthService) Login(email, password string) (string, error) {
	// 从数据库查询用户 (假设已实现)
	hashed := sha256.Sum256([]byte(password))
	user := db.GetUserFromDB(email) // 伪代码
	if user.Password != hex.EncodeToString(hashed[:]) {
		return "", fmt.Errorf("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})
	return token.SignedString(as.secret)
}

func (as *AuthService) ValidateToken(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return as.secret, nil
	})
	if err != nil || !token.Valid {
		return "", fmt.Errorf("invalid token")
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims["user_id"].(string), nil
}

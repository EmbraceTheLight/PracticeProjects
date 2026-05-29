package auth

import (
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"mall/config"
	"mall/service/do"
	"time"
)

type JWTAuth struct {
	Secret            string
	AccessExpireTime  time.Duration
	RefreshExpireTime time.Duration
}

func NewJWTAuth(cfg *config.Jwt) *JWTAuth {
	return &JWTAuth{
		Secret:            cfg.Secret,
		AccessExpireTime:  time.Duration(cfg.AccessTokenExpiration) * time.Hour,
		RefreshExpireTime: time.Duration(cfg.RefreshTokenExpiration) * 24 * time.Hour,
	}
}

const (
	BearerTokenPrefix  = "Bearer"
	AccessTokenHeader  = "Authorization"
	RefreshTokenHeader = "Refresh-Token"

	// TokenIssuer is the default token issuer
	TokenIssuer = "ZEY_HUNTER_ETL"
)

// AccessTokenClaims access token claims, which is used to verify access token
type AccessTokenClaims struct {
	UserId     int64     `json:"user_id"`
	Name       string    `json:"name"`
	NickName   string    `json:"nick_name"`
	Sex        int64     `json:"sex"`
	Status     int64     `json:"status"`
	Mobile     string    `json:"mobile"`
	LarkOpenID string    `json:"lark_open_id"`
	CreateAt   time.Time `json:"create_at"`
	UpdateAt   time.Time `json:"update_at"`
	jwt.RegisteredClaims
}

func NewAccessTokenClaims() *AccessTokenClaims {
	return &AccessTokenClaims{}
}
func (t *AccessTokenClaims) isExpired() bool {
	// if token.ExpiresAt before now, the token is expired, return true
	//else, return false
	return t.ExpiresAt.Time.Before(time.Now())
}

// padding pads the access token claims
func (t *AccessTokenClaims) padding(adminUser *do.GenerateTokenReq, duration time.Duration) {
	t.UserId = adminUser.UserID
	t.Name = adminUser.Name
	t.NickName = adminUser.NickName
	t.Status = adminUser.Status
	t.Sex = adminUser.Sex
	t.LarkOpenID = adminUser.LarkOpenID
	t.Mobile = adminUser.Mobile
	t.CreateAt = adminUser.CreateAt
	t.UpdateAt = adminUser.UpdateAt
	t.Issuer = TokenIssuer
	t.ExpiresAt = jwt.NewNumericDate(time.Now().Add(duration))

}

// getTokenString returns the token string of access token claims, which is padded
func (t *AccessTokenClaims) getTokenString(secret string) (tokenString string, err error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	tokenString, err = token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return
}

// RefreshTokenClaims is used to refresh access token
type RefreshTokenClaims struct {
	jwt.RegisteredClaims
}

func NewRefreshTokenClaims() *RefreshTokenClaims {
	return &RefreshTokenClaims{}
}
func (t *RefreshTokenClaims) isExpired() bool {
	return t.ExpiresAt.Time.Before(time.Now())
}

// padding pads the refresh token claims
func (t *RefreshTokenClaims) padding(duration time.Duration) {
	t.Issuer = TokenIssuer
	t.ExpiresAt = jwt.NewNumericDate(time.Now().Add(duration))
}

// getTokenString returns the token string of refresh token claims,which is padded
func (t *RefreshTokenClaims) getTokenString(secret string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, t)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// GenerateToken create access token and refresh token
func (jwtAuthorizer *JWTAuth) GenerateToken(ctx context.Context, adminUser *do.GenerateTokenReq) (accessToken, refreshToken string, err error) {
	// Create access token
	atokenClaims := NewAccessTokenClaims()
	atokenClaims.padding(adminUser, jwtAuthorizer.AccessExpireTime)
	accessToken, err = atokenClaims.getTokenString(jwtAuthorizer.Secret)
	if err != nil {
		return "", "", err
	}

	// Create refresh token
	rtokenClaims := NewRefreshTokenClaims()
	rtokenClaims.padding(jwtAuthorizer.RefreshExpireTime)
	refreshToken, err = rtokenClaims.getTokenString(jwtAuthorizer.Secret)
	if err != nil {
		return "", "", err
	}
	return
}

// GetToken 函数用于解析和验证JWT令牌
// 参数:
//
//	tokenString: 包含Bearer认证信息的字符串
//	secret: 用于验证JWT签名的密钥
//
// 返回值:
//
//	tokenClaims: 解析后的JWT令牌指针
//	err: 错误信息，如果解析或验证失败则返回错误
func (jwtAuthorizer *JWTAuth) GetToken(ctx context.Context, tokenString string) (*jwt.Token, error) {
	// 检查认证格式是否正确，必须包含两部分且第一部分为"Bearer"
	if len(tokenString) == 0 {
		return nil, errors.New("token is empty")
	}
	// 创建AccessTokenClaims结构体实例，用于存储解析后的令牌声明
	atokenClaims := new(AccessTokenClaims)
	// 使用ParseWithClaims函数解析令牌，验证签名并提取声明
	// 参数分别为: 令牌字符串、声明结构体指针和签名验证函数
	tokenClaims, err := jwt.ParseWithClaims(tokenString, atokenClaims, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtAuthorizer.Secret), nil
	})

	// 处理解析过程中可能出现的错误
	if err != nil {
		return nil, err
	}
	return tokenClaims, nil
}

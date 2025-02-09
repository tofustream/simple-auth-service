package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/tofustream/simple-auth-service/config"
	"github.com/tofustream/simple-auth-service/models"
	"github.com/tofustream/simple-auth-service/security"
)

const (
	accessTokenCookie  = "access_token"
	refreshTokenCookie = "refresh_token"
)

func Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindBodyWithJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashPassword, err := security.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user := models.User{Email: req.Email, Password: hashPassword}

	if err := config.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Email already exists"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User registered successfully"})
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	// JSON のバインド
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// ユーザー取得
	var user models.User
	if err := config.DB.Where("email = ?", req.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// パスワードの検証
	if err := security.CheckPasswordHash(req.Password, user.Password); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// JWT の生成
	accessToken, err := security.GenerateJWT(user.ID, security.DefaultAccessTokenExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	refreshToken, err := security.GenerateJWT(user.ID, security.DefaultRefreshTokenExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate refresh token"})
		return
	}

	// `secure` フラグ（HTTPS なら true, HTTP なら false）
	secure := c.Request.TLS != nil

	// Cookie にセット
	c.SetCookie(accessTokenCookie, accessToken, security.DefaultAccessTokenExpiry, "/", "", secure, true)
	c.SetCookie(refreshTokenCookie, refreshToken, security.DefaultRefreshTokenExpiry, "/", "", secure, true)

	c.JSON(http.StatusOK, gin.H{"message": "Login successful"})
}

func RefreshToken(c *gin.Context) {
	// Cookie からリフレッシュトークンを取得
	refreshToken, err := c.Cookie(refreshTokenCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token required"})
		return
	}

	// トークンを検証
	token, err := security.ValidateJWT(refreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// `token.Claims` を `jwt.MapClaims` にキャストして使う
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	userID := uint(claims["sub"].(float64))

	// 新しいアクセストークンのみ発行
	newAccessToken, err := security.GenerateJWT(userID, security.DefaultAccessTokenExpiry)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate access token"})
		return
	}

	secure := c.Request.TLS != nil

	// 新しいアクセストークンを Cookie にセット
	c.SetCookie(accessTokenCookie, newAccessToken, security.DefaultAccessTokenExpiry, "/", "", secure, true)

	c.JSON(http.StatusOK, gin.H{"message": "Access token refreshed"})
}

func Logout(c *gin.Context) {
	// Cookie を空にして即時無効化
	c.SetCookie(accessTokenCookie, "", -1, "/", "", false, true)
	c.SetCookie(refreshTokenCookie, "", -1, "/", "", false, true)

	c.JSON(http.StatusOK, gin.H{"message": "Logged out successfully"})
}

func GetMe(c *gin.Context) {
	// Cookie から access_token を取得
	token, err := c.Cookie(accessTokenCookie)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token required"})
		return
	}

	// JWT を検証
	parsedToken, err := security.ValidateJWT(token)
	if err != nil || !parsedToken.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	// クレームから user_id を取得
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token claims"})
		return
	}

	userID := uint(claims["sub"].(float64))

	// DB からユーザーの email を取得
	var user models.User
	if err := config.DB.Where("id = ?", userID).First(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// ユーザー ID とメールアドレスを返す
	c.JSON(http.StatusOK, gin.H{
		"user_id": userID,
		"email":   user.Email,
	})
}

package services

import (
	"golang-app/dto"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

const (
	// TODO: Jwt Secret key is Demo only
	AccessTokenCookieName  = "access-token"
	JwtSecretKey           = "access-secret-key"
	RefreshTokenCookieName = "refresh-token"
	JwtRefreshSecretKey    = "refresh-secret-key"
)

type Claims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

func GenerateTokensAndSetCookies(user *dto.UserDto, etx echo.Context) error {
	accessToken, exp, err := generateAccessToken(user)
	if err != nil {
		return err
	}
	setTokenCookie(AccessTokenCookieName, accessToken, exp, etx)

	refreshToken, exp, err := generateRefreshToken(user)
	if err != nil {
		return err
	}
	setTokenCookie(RefreshTokenCookieName, refreshToken, exp, etx)

	return nil
}

func generateAccessToken(user *dto.UserDto) (string, time.Time, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	return generateToken(user, expirationTime, []byte(JwtSecretKey))
}

func generateRefreshToken(user *dto.UserDto) (string, time.Time, error) {
	// Declare the expiration time of the token - 24 hours.
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(JwtRefreshSecretKey))
}

func generateToken(user *dto.UserDto, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	// Create the JWT claims, which includes the username and expiry time.
	claims := &Claims{
		ID:       user.ID,
		Username: user.Username,
		StandardClaims: jwt.StandardClaims{
			// In JWT, the expiry time is expressed as unix seconds.
			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the HS256 algorithm used for signing, and the claims.
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Create the JWT string
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", time.Now(), err
	}

	return tokenString, expirationTime, nil
}

// Creating a new cookie, which will store the valid JWT token.
func setTokenCookie(name, token string, expiration time.Time, etx echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"

	// Http-only helps mitigate the risk of client side script accessing the protected cookie.
	cookie.HttpOnly = true

	etx.SetCookie(cookie)
}

// JWTErrorChecker will be executed when user try to access a protected path.
func JWTErrorChecker(etc echo.Context, err error) error {
	// Redirects to the main page.
	return etc.Redirect(http.StatusMovedPermanently, "/")
}

// TokenRefresherMiddleware middleware, which refreshes JWT tokens if the access token is about to expire.
func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(etx echo.Context) error {
		// If the user is not authenticated (no user token data in the context), don't do anything.
		if etx.Get("user") == nil {
			return next(etx)
		}
		// Gets user token from the context.
		u := etx.Get("user").(*jwt.Token)

		claims := u.Claims.(*Claims)

		// We ensure that a new token is not issued until enough time has elapsed.
		// In this case, a new token will only be issued if the old token is within
		// 15 mins of expiry.
		if time.Until(time.Unix(claims.ExpiresAt, 0)) < 15*time.Minute {
			// Gets the refresh token from the cookie.
			rc, err := etx.Cookie(RefreshTokenCookieName)
			if err == nil && rc != nil {
				// Parses token and checks if it valid.
				tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
					return []byte(JwtRefreshSecretKey), nil
				})
				if err != nil {
					if err == jwt.ErrSignatureInvalid {
						etx.Response().Writer.WriteHeader(http.StatusUnauthorized)
					}
				}

				if tkn != nil && tkn.Valid {
					// If everything is good, update tokens.
					_ = GenerateTokensAndSetCookies(&dto.UserDto{
						ID:       claims.ID,
						Username: claims.Username,
					}, etx)
				}
			}
		}

		return next(etx)
	}
}

// GuestMiddleware middleware, which blocks user from accessing guest routes.
func GuestMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		accessToken, err := c.Cookie(AccessTokenCookieName)
		if err != nil {
			return next(c)
		}
		if accessToken.Value != "" {
			// TODO: Fix the redirect
			return next(c)
			// return c.Redirect(http.StatusMovedPermanently, "/chat")
		}
		return next(c)
	}
}

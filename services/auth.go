package services

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"

	"github.com/alextilot/golang-htmx-chatapp/config"
)

const (
	AccessTokenCookieName  = "access-token"
	RefreshTokenCookieName = "refresh-token"
	oneHour                = 1 * time.Hour
	twentyFourHours        = 24 * time.Hour
	fifteenMinutes         = 15 * time.Minute
)

type Claims struct {
	UserID string `json:"userId"`
	jwt.StandardClaims
}

func GenerateTokensAndSetCookies(userID string, ctx echo.Context) error {
	accessToken, exp, err := generateAccessToken(userID)
	if err != nil {
		return err
	}
	setTokenCookie(AccessTokenCookieName, accessToken, exp, ctx)

	refreshToken, exp, err := generateRefreshToken(userID)
	if err != nil {
		return err
	}
	setTokenCookie(RefreshTokenCookieName, refreshToken, exp, ctx)

	return nil
}

func generateAccessToken(userID string) (string, time.Time, error) {
	expirationTime := time.Now().Add(oneHour)
	return generateToken(userID, expirationTime, []byte(config.Cfg.JwtSecretKey))
}

func generateRefreshToken(userID string) (string, time.Time, error) {
	expirationTime := time.Now().Add(twentyFourHours)
	return generateToken(userID, expirationTime, []byte(config.Cfg.JwtResfeshSecretKey))
}

func generateToken(userID string, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	// Create the JWT claims, which includes the userID and expiry time.
	claims := &Claims{
		UserID: userID,
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
func setTokenCookie(name string, token string, expiration time.Time, ctx echo.Context) {
	cookie := new(http.Cookie)
	cookie.Name = name
	cookie.Value = token
	cookie.Expires = expiration
	cookie.Path = "/"

	// Http-only helps mitigate the risk of client side script accessing the protected cookie.
	cookie.HttpOnly = true

	ctx.SetCookie(cookie)
}

// JWTErrorChecker will be executed when user try to access a protected path.
func JWTErrorChecker(ctx echo.Context, err error) error {
	// Redirects to the main page.
	return ctx.Redirect(http.StatusMovedPermanently, "/")
}

// TokenRefresherMiddleware middleware, which refreshes JWT tokens if the access token is about to expire.
func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		// If the user is not authenticated (no user token data in the context), don't do anything.
		if ctx.Get("user") == nil {
			return next(ctx)
		}
		// Gets user token from the context.
		u := ctx.Get("user").(*jwt.Token)

		claims := u.Claims.(*Claims)

		// We ensure that a new token is not issued until enough time has elapsed.
		// In this case, a new token will only be issued if the old token is within
		// 15 mins of expiry.
		if time.Until(time.Unix(claims.ExpiresAt, 0)) >= fifteenMinutes {
			return next(ctx)
		}

		// Gets the refresh token from the cookie.
		rc, err := ctx.Cookie(RefreshTokenCookieName)
		if err != nil || rc == nil {
			return next(ctx)
		}

		// Parses token and checks if it valid.
		tkn, err := jwt.ParseWithClaims(rc.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.Cfg.JwtResfeshSecretKey), nil
		})
		if err != nil {
			if err == jwt.ErrSignatureInvalid {
				ctx.Response().Writer.WriteHeader(http.StatusUnauthorized)
			}
			return next(ctx)
		}

		if tkn != nil && tkn.Valid {
			// If everything is good, update tokens.
			_ = GenerateTokensAndSetCookies(claims.UserID, ctx)
		}
		return next(ctx)
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

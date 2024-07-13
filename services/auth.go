package services

import (
	"errors"
	"net/http"
	"time"

	"github.com/alextilot/golang-htmx-chatapp/model"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

const (
	// TODO: Jwt Secret key is for Demo only
	AccessTokenCookieName  = "access-token"
	JwtSecretKey           = "access-secret-key"
	RefreshTokenCookieName = "refresh-token"
	JwtRefreshSecretKey    = "refresh-secret-key"
)

var (
	ErrMissingAuthCookie = errors.New("Request header is missing auth cookie.")
	ErrNullAuthToken     = errors.New("Null authentication token.")
	ErrInvalidAuthToken  = errors.New("Invalid authentication token.")
)

type jwtCustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func jwtErrorChecker(etx echo.Context, err error) error {
	if err != nil {
		return etx.Redirect(http.StatusMovedPermanently, "/login")
	}
	return nil
}

func EchoMiddlewareJWTConfig() echo.MiddlewareFunc {
	return echojwt.WithConfig(echojwt.Config{
		SigningKey:   []byte(JwtSecretKey),
		TokenLookup:  "cookie:access-token",
		ErrorHandler: jwtErrorChecker,
		NewClaimsFunc: func(etx echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
	})
}

type UserContext struct {
	Name       string
	IsLoggedIn bool
}

func GetUserContext(etx echo.Context) (UserContext, error) {
	userContext := UserContext{Name: "", IsLoggedIn: false}

	cookie, err := etx.Cookie(AccessTokenCookieName)
	if cookie == nil || err != nil {
		return userContext, ErrMissingAuthCookie
	}

	claims := &jwtCustomClaims{}
	token, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(JwtSecretKey), nil
	})

	if err != nil {
		return userContext, err
	}
	if token == nil {
		return userContext, ErrNullAuthToken
	}
	if !token.Valid {
		return userContext, ErrInvalidAuthToken
	}

	userContext.Name = claims.Username
	userContext.IsLoggedIn = token.Valid

	return userContext, nil
}

func RemoveTokensAndCookies(etx echo.Context) {
	setTokenCookie(AccessTokenCookieName, "", time.UnixMicro(0), etx)
	setTokenCookie(RefreshTokenCookieName, "", time.UnixMicro(0), etx)
}

func GenerateTokensAndSetCookies(user *model.User, etx echo.Context) error {
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

func generateAccessToken(user *model.User) (string, time.Time, error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	return generateToken(user, expirationTime, []byte(JwtSecretKey))
}

func generateRefreshToken(user *model.User) (string, time.Time, error) {
	// Declare the expiration time of the token - 24 hours.
	expirationTime := time.Now().Add(24 * time.Hour)

	return generateToken(user, expirationTime, []byte(JwtRefreshSecretKey))
}

func generateToken(user *model.User, expirationTime time.Time, secret []byte) (string, time.Time, error) {
	// Create the JWT claims, which includes the username and expiry time.
	claims := &jwtCustomClaims{
		user.Username,
		jwt.RegisteredClaims{
			// In JWT, the expiry time is expressed as unix seconds.
			ExpiresAt: jwt.NewNumericDate(expirationTime),
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

// TokenRefresherMiddleware middleware, which refreshes JWT tokens if the access token is about to expire.
func TokenRefresherMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(etx echo.Context) error {
		// If the user is not authenticated (no user token data in the context), don't do anything.
		if etx.Get("user") == nil {
			return next(etx)
		}
		// Gets user token from the context.
		user := etx.Get("user").(*jwt.Token)

		claims := user.Claims.(*jwtCustomClaims)

		// We ensure that a new token is not issued until enough time has elapsed.
		// In this case, a new token will only be issued if the old token is within
		// 15 mins of expiry.
		if time.Until(time.Unix(claims.ExpiresAt.Unix(), 0)) < 15*time.Minute {
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
					_ = GenerateTokensAndSetCookies(&model.User{
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
	return func(etx echo.Context) error {
		accessToken, err := etx.Cookie(AccessTokenCookieName)
		if err != nil {
			return next(etx)
		}
		if accessToken.Value != "" {
			// TODO: Fix the redirect
			// return next(etx)
			// return c.Redirect(http.StatusMovedPermanently, "/")
			etx.Response().Header().Set("HX-Redirect", "/login")
			return etx.String(http.StatusMovedPermanently, "Success")
		}
		return next(etx)
	}
}

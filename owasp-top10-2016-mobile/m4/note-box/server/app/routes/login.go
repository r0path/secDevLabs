package routes

import (
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app/db"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app/types"
	"github.com/globocom/secDevLabs/owasp-top10-2016-mobile/m4/note-box/server/app/util"
	"github.com/labstack/echo"
)

// Simple in-memory rate limiting to mitigate brute-force attacks on the login endpoint.
// This is intentionally conservative and per-IP. For production, consider a distributed
// rate limiter (Redis, memcached) to handle multiple instances and persistent state.

type attemptInfo struct {
	Count       int
	LastAttempt time.Time
	BlockExpires time.Time
}

var (
	loginAttempts      = make(map[string]*attemptInfo)
	loginAttemptsMutex sync.Mutex
	maxAttempts        = 5
	attemptWindow      = time.Minute
	blockDuration      = time.Minute * 15
)

// Login attempts to loggin an user.
func Login(c echo.Context) error {
	jwtSecret := os.Getenv("M4_SECRET")

	u := new(types.RequestUser)
	if err := c.Bind(u); err != nil {
		return err
	}

	attemptUsername := strings.TrimSpace(u.Username)
	attemptPassword := strings.TrimSpace(u.Password)

	// Rate limit based on client IP to mitigate brute-force attempts.
	ip := c.RealIP()
	now := time.Now()

	loginAttemptsMutex.Lock()
	ai, ok := loginAttempts[ip]
	if !ok {
		ai = &attemptInfo{Count: 0, LastAttempt: now}
		loginAttempts[ip] = ai
	}
	// If currently blocked, deny immediately.
	if ai.BlockExpires.After(now) {
		loginAttemptsMutex.Unlock()
		return c.JSON(http.StatusTooManyRequests, "Too many login attempts. Try again later.")
	}
	// Reset the counter if outside the attempt window.
	if now.Sub(ai.LastAttempt) > attemptWindow {
		ai.Count = 0
	}
	ai.LastAttempt = now
	loginAttemptsMutex.Unlock()

	user, err := db.FindOneUser(attemptUsername)
	if err != nil {
		// Increment attempt count for failed attempt (username not found).
		loginAttemptsMutex.Lock()
		ai.Count++
		if ai.Count >= maxAttempts {
			ai.BlockExpires = now.Add(blockDuration)
		}
		loginAttemptsMutex.Unlock()

		return c.JSON(http.StatusNotFound, "Username or password is wrong or the user doesn't exist")
	}

	if !util.VerifyHash(attemptPassword, user.HashedPassword, user.Salt) {
		// Wrong password: increment attempt count and possibly block.
		loginAttemptsMutex.Lock()
		ai.Count++
		if ai.Count >= maxAttempts {
			ai.BlockExpires = now.Add(blockDuration)
		}
		loginAttemptsMutex.Unlock()

		return c.JSON(http.StatusNotFound, "Username or password is wrong or the user doesn't exist")
	}

	// Successful login: clear attempt record for this IP.
	loginAttemptsMutex.Lock()
	delete(loginAttempts, ip)
	loginAttemptsMutex.Unlock()

	if user.IsLoggedIn {
		return c.JSON(http.StatusConflict, "User is already logged in!")
	}

	// Create token
	token := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = attemptUsername
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	// Generate encoded token and send it as response.
	t, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Error generating user token")
	}

	db.UpdateUserLoggedIn(attemptUsername, true)

	return c.JSON(http.StatusOK, echo.Map{
		"sessionToken": t,
	})
}

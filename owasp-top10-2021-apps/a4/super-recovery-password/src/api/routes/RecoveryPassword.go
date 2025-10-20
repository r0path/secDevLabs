package routes

import (
	"api/database"
	"api/services"
	"api/types"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/labstack/echo"
)

// Simple in-memory brute-force protection for recovery answers.
// This is intentionally conservative: it adds a per-login counter and lockout
// period. It's not distributed and resets on process restart, but it is a
// safe, low-risk mitigation that avoids extensive refactors.

var (
	recoveryAttempts = make(map[string]*attemptInfo)
	attemptsMu        sync.Mutex
	maxRecoveryAttempts = 5
	lockoutDuration     = 15 * time.Minute
)

type attemptInfo struct {
	Count      int
	LastFailed time.Time
	LockedUntil time.Time
}

func incrementFailed(login string) {
	attemptsMu.Lock()
	defer attemptsMu.Unlock()
	ai, ok := recoveryAttempts[login]
	if !ok {
		ai = &attemptInfo{}
		recoveryAttempts[login] = ai
	}
	ai.Count++
	ai.LastFailed = time.Now()
	if ai.Count >= maxRecoveryAttempts {
		ai.LockedUntil = time.Now().Add(lockoutDuration)
		ai.Count = 0
	}
}

func resetAttempts(login string) {
	attemptsMu.Lock()
	defer attemptsMu.Unlock()
	delete(recoveryAttempts, login)
}

func RecoveryPassword(c echo.Context) (err error) {
	u := new(types.RecoveryPasswordAnswers)
	if err = c.Bind(u); err != nil {
		return
	}
	u.Login = strings.ToLower(u.Login)

	// Check for temporary lockout before attempting verification
	attemptsMu.Lock()
	if ai, ok := recoveryAttempts[u.Login]; ok {
		if ai.LockedUntil.After(time.Now()) {
			attemptsMu.Unlock()
			return c.JSON(http.StatusTooManyRequests, echo.Map{"message": "too many attempts, try again later"})
		}
	}
	attemptsMu.Unlock()

	recoveryPasswordAnswers := types.RecoveryPasswordAnswers{
		Login:        u.Login,
		FirstAnswer:  u.FirstAnswer,
		SecondAnswer: u.SecondAnswer,
	}

	answers, err := database.RecoveryPassword(recoveryPasswordAnswers.Login, recoveryPasswordAnswers.FirstAnswer, recoveryPasswordAnswers.SecondAnswer)
	if err != nil {
		// Increment failed attempt and potentially lockout
		incrementFailed(recoveryPasswordAnswers.Login)
		return c.JSON(http.StatusConflict, echo.Map{"message": "incorrect answers!"})
	}

	if answers.FirstAnswer != recoveryPasswordAnswers.FirstAnswer || answers.SecondAnswer != recoveryPasswordAnswers.SecondAnswer {
		incrementFailed(recoveryPasswordAnswers.Login)
		return c.JSON(http.StatusConflict, echo.Map{"message": "incorrect answers!"})
	}

	// Successful verification: reset any counters for this login
	resetAttempts(recoveryPasswordAnswers.Login)

	token, err := services.GenerateJwt(recoveryPasswordAnswers.Login, true)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"token": "Error to generate token.",
		})
	}
	fmt.Println(token)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

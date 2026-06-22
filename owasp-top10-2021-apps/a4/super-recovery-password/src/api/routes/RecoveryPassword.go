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

// recoveryAttempts tracks failed recovery attempts per login for rate limiting.
var (
	recoveryMu       sync.Mutex
	recoveryAttempts = make(map[string][]time.Time)
)

const (
	maxRecoveryAttempts = 5
	recoveryWindow      = time.Hour
)

func checkRecoveryRateLimit(login string) bool {
	recoveryMu.Lock()
	defer recoveryMu.Unlock()

	now := time.Now()
	cutoff := now.Add(-recoveryWindow)

	// Purge attempts outside the window.
	attempts := recoveryAttempts[login]
	valid := attempts[:0]
	for _, t := range attempts {
		if t.After(cutoff) {
			valid = append(valid, t)
		}
	}
	recoveryAttempts[login] = valid

	if len(valid) >= maxRecoveryAttempts {
		return false
	}
	recoveryAttempts[login] = append(recoveryAttempts[login], now)
	return true
}

func RecoveryPassword(c echo.Context) (err error) {
	u := new(types.RecoveryPasswordAnswers)
	if err = c.Bind(u); err != nil {
		return
	}
	u.Login = strings.ToLower(u.Login)

	// Rate limit: max 5 attempts per account per hour to prevent brute-force.
	if !checkRecoveryRateLimit(u.Login) {
		return c.JSON(http.StatusTooManyRequests, echo.Map{
			"message": "too many recovery attempts, please try again later",
		})
	}

	recoveryPasswordAnswers := types.RecoveryPasswordAnswers{
		Login:        u.Login,
		FirstAnswer:  u.FirstAnswer,
		SecondAnswer: u.SecondAnswer,
	}

	answers, err := database.RecoveryPassword(recoveryPasswordAnswers.Login, recoveryPasswordAnswers.FirstAnswer, recoveryPasswordAnswers.SecondAnswer)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{"message": "incorrect answers!"})
	}

	if answers.FirstAnswer != recoveryPasswordAnswers.FirstAnswer || answers.SecondAnswer != recoveryPasswordAnswers.SecondAnswer {
		return c.JSON(http.StatusConflict, echo.Map{"message": "incorrect answers!"})
	}

	token, err := services.GenerateJwt(recoveryPasswordAnswers.Login, true)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"token": "Error to generate token.",
		})
	}

	// Store the generated recovery token server-side bound to the user to make it single-use.
	if err := database.StoreRecoveryToken(recoveryPasswordAnswers.Login, token); err != nil {
		fmt.Println("failed to store recovery token:", err)
		// Continue and return token to user even if storing fails; storing is best-effort in this example.
	}

	fmt.Println(token)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

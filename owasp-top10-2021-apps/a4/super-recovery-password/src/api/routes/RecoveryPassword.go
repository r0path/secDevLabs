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

	// reject empty submitted answers to avoid accepting blank answers
	if strings.TrimSpace(recoveryPasswordAnswers.FirstAnswer) == "" || strings.TrimSpace(recoveryPasswordAnswers.SecondAnswer) == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": "answers cannot be empty"})
	}

	answers, err := database.RecoveryPassword(recoveryPasswordAnswers.Login, recoveryPasswordAnswers.FirstAnswer, recoveryPasswordAnswers.SecondAnswer)
	if err != nil {
		return c.JSON(http.StatusConflict, echo.Map{"message": "incorrect answers!"})
	}

	// reject cases where stored answers are empty to prevent empty-match bypass
	if strings.TrimSpace(answers.FirstAnswer) == "" || strings.TrimSpace(answers.SecondAnswer) == "" {
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
	fmt.Println(token)
	return c.JSON(http.StatusOK, echo.Map{
		"token": token,
	})
}

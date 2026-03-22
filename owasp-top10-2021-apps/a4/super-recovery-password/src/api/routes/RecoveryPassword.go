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

var (
	recoveryAttempts = make(map[string]*recoveryAttempt)
	attemptsMu       sync.Mutex
)

type recoveryAttempt struct {
	count    int
	windowStart time.Time
}

const (
	maxRecoveryAttempts  = 5
	recoveryWindowPeriod = 1 * time.Hour
)

func checkRateLimit(login string) bool {
	attemptsMu.Lock()
	defer attemptsMu.Unlock()

	now := time.Now()
	attempt, exists := recoveryAttempts[login]
	if !exists || now.Sub(attempt.windowStart) > recoveryWindowPeriod {
		recoveryAttempts[login] = &recoveryAttempt{count: 0, windowStart: now}
		return true
	}
	return attempt.count < maxRecoveryAttempts
}

func recordFailedAttempt(login string) {
	attemptsMu.Lock()
	defer attemptsMu.Unlock()

	if attempt, exists := recoveryAttempts[login]; exists {
		attempt.count++
	}
}

func RecoveryPassword(c echo.Context) (err error) {
	u := new(types.RecoveryPasswordAnswers)
	if err = c.Bind(u); err != nil {
		return
	}
	u.Login = strings.ToLower(u.Login)

	if !checkRateLimit(u.Login) {
		return c.JSON(http.StatusTooManyRequests, echo.Map{
			"message": "Too many recovery attempts. Please try again later.",
		})
	}

	recoveryPasswordAnswers := types.RecoveryPasswordAnswers{
		Login:        u.Login,
		FirstAnswer:  u.FirstAnswer,
		SecondAnswer: u.SecondAnswer,
	}

	answers, err := database.RecoveryPassword(recoveryPasswordAnswers.Login, recoveryPasswordAnswers.FirstAnswer, recoveryPasswordAnswers.SecondAnswer)
	if err != nil {
		recordFailedAttempt(u.Login)
		return c.JSON(http.StatusConflict, echo.Map{"message": "incorrect answers!"})
	}

	if answers.FirstAnswer != recoveryPasswordAnswers.FirstAnswer || answers.SecondAnswer != recoveryPasswordAnswers.SecondAnswer {
		recordFailedAttempt(u.Login)
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

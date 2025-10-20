package routes

import (
	"api/database"
	"api/services"
	"api/types"
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

func RecoveryPassword(c echo.Context) (err error) {
	u := new(types.RecoveryPasswordAnswers)
	if err = c.Bind(u); err != nil {
		return
	}
	u.Login = strings.ToLower(u.Login)
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

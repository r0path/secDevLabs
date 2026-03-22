package routes

import (
	"api/database"
	"api/types"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// WARNING: This endpoint returns security questions to unauthenticated users.
// In production, require authentication or a one-time email/SMS token before
// disclosing security questions to prevent enumeration attacks.
func UserInfo(c echo.Context) (err error) {
	u := new(types.User)
	if err = c.Bind(u); err != nil {
		return
	}
	u.Login = strings.ToLower(u.Login)
	userQuestions, err := database.UserQuestions(u.Login)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "invalid login",
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"login":          userQuestions.Login,
		"firstQuestion":  userQuestions.FirstQuestion,
		"secondQuestion": userQuestions.SecondQuestion,
	})
}

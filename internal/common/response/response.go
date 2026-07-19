package response

import (
	"net/http"

	"github.com/gin-gonic/gin"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
)

type Envelope struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Envelope{Code: 0, Message: "ok", Data: data})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Envelope{Code: 0, Message: "created", Data: data})
}

func Error(c *gin.Context, err error) {
	ae := apperr.AsAppError(err)
	c.JSON(ae.HTTP, Envelope{Code: ae.Code, Message: ae.Message, Data: nil})
}

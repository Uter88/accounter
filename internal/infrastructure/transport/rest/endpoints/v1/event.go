package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (e *v1Engine) getEventsList(c *gin.Context) {
	if result, err := e.eventService.GetEventList(c); err != nil {
		e.writeErr(c, http.StatusInternalServerError, err)

	} else {
		e.writeOk(c, result)
	}
}

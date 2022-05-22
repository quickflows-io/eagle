package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"

	"github.com/go-eagle/eagle/internal/web"
	"github.com/go-eagle/eagle/pkg/log"
)

// Logout user logout
func Logout(c *gin.Context) {
	// delete cookie information
	session := web.GetCookieSession(c)
	session.Options = &sessions.Options{
		Domain: "",
		Path:   "/",
		MaxAge: -1,
	}
	err := session.Save(web.Request(c), web.ResponseWriter(c))
	if err != nil {
		log.Warnf("[user] logout save session err: %v", err)
		c.Abort()
		return
	}

	// redirect to the original page
	c.Redirect(http.StatusSeeOther, c.Request.Referer())
}

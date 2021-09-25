package main

import (
	// "fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"

	// "github.com/tpkeeper/gin-dump"
	"github.com/uptrace/bun"
)

var DB *bun.DB

func main() {
	loadDBCursor()
	r := engine()
	r.Use(gin.Logger())
	// r.Use(gindump.Dump())
	// r.Use(gindump.DumpWithOptions(true, true, false, true, false, func(dumpStr string) {
	// 	fmt.Println(dumpStr)
	// }))
	gin.SetMode(gin.DebugMode)
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Unable to start:", err)
	}
}

func engine() *gin.Engine {
	r := gin.New()
	r.Use(sessions.Sessions("mysession", sessions.NewCookieStore([]byte("secret"))))
	// store := cookie.NewStore([]byte("secret"))
	// r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLGlob(("templates/*"))

	r.POST("/login", loginUser)
	r.GET("/logout", logoutUser)
	r.GET("/session", showSession)

	r.GET("/list", listNotes).Use(AuthRequired)
	r.GET("/list/:id", getNoteById).Use(AuthRequired)
	r.POST("/create", postNote).Use(AuthRequired)
	r.POST("/update", updateNote).Use(AuthRequired)
	r.POST("/delete", deleteNote).Use(AuthRequired)

	return r
}

// AuthRequired is a simple middleware to check the session
func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("username")
	if user == nil {
		// Abort the request with the appropriate error code
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}
	// Continue down the chain to handler etc
	c.Next()
}

func showSession(c *gin.Context) {
	session := sessions.Default(c)
	// session := sessions.Default(c)
	// userid := session.Get("userid")
	// username := session.Get("username")
	// fmt.Println(hashString("crove"))
	// message := fmt.Sprintf("You are: %s with userid: %d", username, userid)
	// c.IndentedJSON(http.StatusOK, gin.H{"message": message})
	c.HTML(http.StatusOK, "session.tmpl", gin.H{"isLoggedIn": checkLoggedIn(c), "username": session.Get("username"), "userid": session.Get("userid")})
	// c.SetCookie("authenticated", "ahunt", 3600, "/", "localhost", false, true)
	// cookie, _ := c.Cookie("authenticated")
	// fmt.Println(cookie)
}

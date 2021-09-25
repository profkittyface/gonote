package main

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func loginUser(c *gin.Context) {
	session := sessions.Default(c)
	username := c.PostForm("username")
	password := c.PostForm("password")
	if checkCredentials(username, password) {
		userid := getUserIdFromUsername(username)
		session.Set("username", username)
		session.Set("userid", userid)
		session.Save()
		c.IndentedJSON(http.StatusOK, gin.H{"message": "logged_in"})
	}

	// ------ JSON based login
	// var (
	// 	jsonContainer map[string]interface{}
	// 	username      string
	// 	password      string
	// )
	// c.BindJSON(&jsonContainer)
	// if val, b := jsonContainer["username"]; b {
	// 	username = val.(string)
	// }
	// if val, b := jsonContainer["password"]; b {
	// 	password = val.(string)
	// }
	// if checkCredentials(username, password) {
	// 	c.SetCookie("authenticated", username, 3600, "/", "localhost", false, true)
	// 	c.IndentedJSON(http.StatusOK, gin.H{"message": "logged_in"})
	// } else {
	// 	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Incorrect username or password"})
	// }
}

func logoutUser(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.IndentedJSON(http.StatusOK, gin.H{"message": "successfully logged out"})
}

func checkCredentials(username string, password string) bool {
	user := User{}
	ctx := context.Background()
	DB.NewSelect().Table("user").Where("username = ?", username).Scan(ctx, &user)
	if user.Id != 0 {
		if username == user.Username && hashString(password) == user.Password {
			return true
		}
	} else {
		return false
	}
	return false
}

func hashString(s string) string {
	h := sha256.Sum256([]byte(s))
	res := fmt.Sprintf("%x", h)
	return res
}

func checkLoggedIn(c *gin.Context) bool {
	session := sessions.Default(c)
	res := session.Get("username")
	if res == nil {
		return false
	}
	return true
}

package main

import (
	"context"
	"net/http"
	"strconv"


	"github.com/gin-gonic/contrib/sessions"

	"github.com/gin-gonic/gin"
)

func listNotes(c *gin.Context) {
  session := sessions.Default(c)
	m := make([]map[string]interface{}, 0)
	var limit int
	// Return map with query string
	req_query := c.Request.URL.Query()
	if val, ok := req_query["limit"]; ok {
		limit, _ = strconv.Atoi(val[0])
	} else {
		limit = 10
	}
	ctx := context.Background()
  userid := session.Get("userid")
	DB.NewSelect().Table("note").Where("userid = ?", userid).Limit(limit).Scan(ctx, &m)
	c.IndentedJSON(http.StatusOK, m)
}

func postNote(c *gin.Context) {
	var jsonContainer map[string]interface{}
	c.BindJSON(&jsonContainer)
	newNote := Note{
		Title:   jsonContainer["title"].(string),
		Content: jsonContainer["content"].(string),
		UserID:  int(jsonContainer["userid"].(float64)),
	}
	ctx := context.Background()
	DB.NewInsert().Model(&newNote).Exec(ctx)
	c.IndentedJSON(http.StatusCreated, newNote)
}

func updateNote(c *gin.Context) {
	ctx := context.Background()
	noteToUpdate := Note{}
	var jsonContainer map[string]interface{}
	c.BindJSON(&jsonContainer)
	noteid := jsonContainer["id"]
	DB.NewSelect().Table("note").Where("id = ?", noteid).Scan(ctx, &noteToUpdate)

	noteToUpdate.Title = jsonContainer["title"].(string)
	noteToUpdate.Content = jsonContainer["content"].(string)

	DB.NewUpdate().Model(&noteToUpdate).Where("id = ?", noteid).Exec(ctx)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "updated"})
}

func deleteNote(c *gin.Context) {
	ctx := context.Background()
	var jsonContainer map[string]interface{}
	c.BindJSON(&jsonContainer)
	noteid := jsonContainer["id"]
	DB.NewDelete().Table("note").Where("id = ?", noteid).Exec(ctx)
	c.IndentedJSON(http.StatusOK, gin.H{"message": "deleted"})
}

func getNoteById(c *gin.Context) {
	note := Note{}
	stringId := c.Param("id")
	// id, _ := strconv.Atoi(stringId)
	ctx := context.Background()
	DB.NewSelect().Table("note").Where("id = ?", stringId).Scan(ctx, &note)
	if note.Id != 0 {
		c.IndentedJSON(http.StatusFound, note)
	} else {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "note not found"})
	}
}

func getUserIdFromUsername(username string) int {
	user := User{}
	ctx := context.Background()
	DB.NewSelect().Table("user").Column("id").Where("username = ?", username).Scan(ctx, &user)
	return user.Id
}

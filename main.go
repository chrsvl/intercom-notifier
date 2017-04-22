package main

import (
  "os"
  "log"
	"net/http"
	"bytes"
	"github.com/gin-gonic/gin"
)

type NotificationEvent struct {
  Data NotificationEventData `json:"data" binding:"required"`
}

type NotificationEventData struct {
	Item NotificationEventDataItem `json:"item" binding:"required"`
}

type NotificationEventDataItem struct {
	Assignee Assignee `json:"assignee" binding:"required"`
}

type Assignee struct {
	Name string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required"`
}

func main() {
	router := gin.Default()
	router.Use(gin.Logger())

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hi")
	})

	router.POST("/", func(c *gin.Context) {
		extractAssignee(c);
	})

	router.Run()
}

func extractAssignee(c *gin.Context) {
	supportTeamName := os.Getenv("SUPPORT_TEAM_NAME")
	var jsonNotificationEvent NotificationEvent
	if c.BindJSON(&jsonNotificationEvent) == nil {
		if jsonNotificationEvent.Data.Item.Assignee.Name == supportTeamName {
			log.Print("New ticket for Support team")
			notifyDevTeamOnSlack()
		}
		c.String(http.StatusOK, "thanks!")
	} else {
		log.Fatal("Could not parse JSON payload sent by Intercom.")
	}
}

func notifyDevTeamOnSlack() {
	webhookUrl := os.Getenv("SLACK_WEBHOOK_URL")
  payload := []byte(`{"text": "<!here>, A new conversation has been assigned to the Support Team. :rungun:"}`)
  req, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(payload))
  req.Header.Set("Content-Type", "application/json")
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      log.Fatal(err)
  }
  defer resp.Body.Close()
}

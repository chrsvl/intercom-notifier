package main

import (
  "os"
  "log"
	"net/http"
	"bytes"
	"github.com/gin-gonic/gin"
	"fmt"
)

type NotificationEvent struct {
  Data NotificationEventData `json:"data" binding:"required"`
}

type NotificationEventData struct {
	Item NotificationEventDataItem `json:"item" binding:"required"`
}

type NotificationEventDataItem struct {
	Id string `json:"id" binding:"required"`
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
	teamName := os.Getenv("TEAM_NAME")
	var jsonNotificationEvent NotificationEvent
	if c.BindJSON(&jsonNotificationEvent) == nil {
		if jsonNotificationEvent.Data.Item.Assignee.Name == teamName {
			log.Print("New ticket for the team")
			notifyDevTeamOnSlack(jsonNotificationEvent.Data.Item.Id)
		}
		c.String(http.StatusOK, "thanks!")
	} else {
		log.Fatal("Could not parse JSON payload sent by Intercom.")
	}
}

func notifyDevTeamOnSlack(conversationId string) {
	webhookUrl := os.Getenv("SLACK_WEBHOOK_URL")
	intercomBaseUrl := os.Getenv("CONVERSATION_BASE_URL")
	conversationUrl := fmt.Sprintf("%s%s", intercomBaseUrl, conversationId)
  payload := []byte(fmt.Sprintf(`{"text": "<!here>, A <%s|new conversation> has been assigned to your team. :rungun:"}`, conversationUrl))
  req, err := http.NewRequest("POST", webhookUrl, bytes.NewBuffer(payload))
  req.Header.Set("Content-Type", "application/json")
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
      log.Fatal(err)
  }
  defer resp.Body.Close()
}

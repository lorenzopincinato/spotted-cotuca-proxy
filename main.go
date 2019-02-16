package main

import (
	"strconv"

	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "running",
		})
	})
	r.POST("/tweet", func(c *gin.Context) {
		request := struct {
			Message        string `json:"message"`
			AccessToken    string `json:"accessToken"`
			AccessSecret   string `json:"accessSecret"`
			ConsumerKey    string `json:"consumerKey"`
			ConsumerSecret string `json:"consumerSecret"`
		}{}

		err := c.BindJSON(&request)

		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error,
			})
		}

		api := anaconda.NewTwitterApiWithCredentials(request.AccessToken, request.AccessSecret, request.ConsumerKey, request.ConsumerSecret)

		tweet, err := api.PostTweet(request.Message, nil)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error,
			})
		}

		c.JSON(200, gin.H{
			"tweetId": tweet.Id,
		})
	})
	r.DELETE("/tweet/:id", func(c *gin.Context) {
		request := struct {
			AccessToken    string `json:"accessToken"`
			AccessSecret   string `json:"accessSecret"`
			ConsumerKey    string `json:"consumerKey"`
			ConsumerSecret string `json:"consumerSecret"`
		}{}

		tweetIDStr := c.Param("id")

		tweetID, err := strconv.ParseInt(tweetIDStr, 10, 64)

		err = c.BindJSON(&request)

		if err != nil {
			c.JSON(400, gin.H{
				"error": err.Error,
			})
		}

		api := anaconda.NewTwitterApiWithCredentials(request.AccessToken, request.AccessSecret, request.ConsumerKey, request.ConsumerSecret)

		_, err = api.DeleteTweet(tweetID, false)

		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error,
			})
		}

		c.JSON(200, nil)
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}

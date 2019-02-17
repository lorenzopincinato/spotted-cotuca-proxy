package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ChimeraCoder/anaconda"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func between(value string, a string, b string) string {
	posFirst := strings.Index(value, a)
	if posFirst == -1 {
		return ""
	}
	posLast := strings.Index(value, b)
	if posLast == -1 {
		return ""
	}
	posFirstAdjusted := posFirst + len(a)
	if posFirstAdjusted >= posLast {
		return ""
	}
	return value[posFirstAdjusted:posLast]
}

func after(value string, a string) string {
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}

func handleError(errorString string) (errorCode int, errorBody []byte, err error) {
	errorCodeStr := between(errorString, "returned status ", ", {")
	errorCode, err = strconv.Atoi(errorCodeStr)

	fmt.Print(errorString)

	errorBody = []byte("{" + after(errorString, ", {"))

	return errorCode, errorBody, err
}

func main() {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://newspottedctc.appspot.com"}

	r.Use(cors.New(config))

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
				"error": err.Error(),
			})
		}

		api := anaconda.NewTwitterApiWithCredentials(request.AccessToken, request.AccessSecret, request.ConsumerKey, request.ConsumerSecret)

		tweet, err := api.PostTweet(request.Message, nil)

		if err != nil {
			errorCode, errorBody, err := handleError(err.Error())

			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})

			} else {
				c.Data(errorCode, "application/json; charset=utf-8", errorBody)
			}

		} else {
			c.JSON(200, gin.H{
				"tweetId": tweet.Id,
			})
		}
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
				"error": err.Error(),
			})
		}

		api := anaconda.NewTwitterApiWithCredentials(request.AccessToken, request.AccessSecret, request.ConsumerKey, request.ConsumerSecret)

		_, err = api.DeleteTweet(tweetID, false)

		if err != nil {
			errorCode, errorBody, err := handleError(err.Error())

			if err != nil {
				c.JSON(500, gin.H{
					"error": err.Error(),
				})

			} else {
				c.Data(errorCode, "application/json; charset=utf-8", errorBody)
			}

		} else {
			c.JSON(200, nil)
		}
	})
	r.Run()
}

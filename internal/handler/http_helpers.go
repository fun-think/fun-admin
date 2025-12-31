package handler

import (
	"errors"
	"io"

	"github.com/gin-gonic/gin"
)

const defaultLanguage = "zh-CN"

func getLanguage(c *gin.Context) string {
	language := c.Query("language")
	if language == "" {
		return defaultLanguage
	}
	return language
}

func messageWithDebugError(message string, err error) string {
	if err == nil {
		return message
	}
	if gin.Mode() != gin.DebugMode {
		return message
	}
	return message + ": " + err.Error()
}

func isEmptyBodyJSONError(err error) bool {
	return errors.Is(err, io.EOF)
}

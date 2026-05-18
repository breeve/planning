package dto

import "github.com/gin-gonic/gin"

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
	Meta    *Meta       `json:"meta,omitempty"`
}

type Meta struct {
	Total int64 `json:"total"`
	Page  int   `json:"page"`
	Limit int   `json:"limit"`
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(200, Response{Success: true, Data: data})
}

func OKList(c *gin.Context, data interface{}, total int64, page, limit int) {
	c.JSON(200, Response{
		Success: true,
		Data:    data,
		Meta:    &Meta{Total: total, Page: page, Limit: limit},
	})
}

func Created(c *gin.Context, data interface{}) {
	c.JSON(201, Response{Success: true, Data: data})
}

func Err(c *gin.Context, status int, msg string) {
	c.JSON(status, Response{Success: false, Error: msg})
}

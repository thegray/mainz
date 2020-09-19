package server

import (
	"github.com/labstack/echo"
)

type CustomContext struct {
	echo.Context
	customInfo map[string]interface{}
}

func (c *CustomContext) SetCustomInfo(key string, i interface{}) {
	c.customInfo[key] = i
}

func (c *CustomContext) GetCustomInfo(key string) interface{} {
	info, ok := c.customInfo[key]
	if !ok {
		return nil
	}
	return info
}

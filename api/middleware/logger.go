package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)


func Logger() fiber.Handler{
    return logger.New(logger.Config{
        Format: "[${ip}]:${port} ${status} - ${method} - ${path} ${latency} - response: ${rsBody} request: ${body}\n",
    }) 
}


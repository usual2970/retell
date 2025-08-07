package middleware

import "github.com/labstack/echo/v4"

// CORS will handle the CORS middleware
func CORS(next echo.HandlerFunc) echo.HandlerFunc {
    return func(c echo.Context) error {
        c.Response().Header().Set("Access-Control-Allow-Origin", "*")
        c.Response().Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
        c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type,Authorization,X-Requested-With")
        c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
        
        // Handle preflight requests
        if c.Request().Method == "OPTIONS" {
            return c.NoContent(204)
        }
        
        return next(c)
    }
}

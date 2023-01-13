package api

import (
	"encoding/json"
	"file/infrastructure"
	"file/internal/controller"
	"file/internal/domain"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slices"
)

type errorresponse struct {
	Message string
}

func Routes(
	router *echo.Echo,
	FileController controller.FileController,
) {

	router.Use(authMiddleware)
	router.POST("/upload", FileController.Upload)

}

func authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		whitelistUrl := []string{"/auth/login", "/auth/register"}
		redisConn := infrastructure.OpenRedis()

		defer redisConn.Close()

		ctx := c.Request().Context()

		// Register public route
		if slices.Contains(whitelistUrl, c.Request().RequestURI) {
			return next(c)
		}

		// Check header Authorization is empty
		if c.Request().Header.Get("Authorization") == "" {
			response := errorresponse{
				Message: "Missing JWT",
			}

			return c.JSON(http.StatusUnauthorized, response)
		}

		// Get value Header Authorization
		token := c.Request().Header.Get("Authorization")

		// Check value Header contain bearer
		if !strings.Contains(token, "Bearer ") {
			response := errorresponse{
				Message: "Missing JWT",
			}

			return c.JSON(http.StatusUnauthorized, response)
		}

		// Delete bearer and left only token
		tokenFix := strings.Replace(token, "Bearer ", "", 1)

		// Parse token & verify secret
		claims := jwt.MapClaims{}
		tokenParse, errFix := jwt.ParseWithClaims(tokenFix, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_KEY")), nil
		})

		// Return if cant parse and validation
		if errFix != nil {
			response := errorresponse{
				Message: "Missing JWT",
			}

			return c.JSON(http.StatusUnauthorized, response)
		}

		// Convert parse to claim
		claim := tokenParse.Claims.(jwt.MapClaims)
		// convert from interface to string
		uuid := fmt.Sprintf("%v", claim["uuid"])

		resUuid, _ := redisConn.Get(ctx, uuid).Result()

		// check uuid if uuid exist pass it
		if resUuid != "" {
			// Parse return to struct
			user := domain.User{}
			json.Unmarshal([]byte(resUuid), &user)
			// set to context
			c.Set("user", user)
			return next(c)
		}

		// Convert exp jwt to time
		var expTime time.Time
		now := time.Now()
		switch iat := claim["exp"].(type) {
		case float64:
			expTime = time.Unix(int64(iat), 0)
		case json.Number:
			v, _ := iat.Int64()
			expTime = time.Unix(v, 0)
		}

		// Find Second diff
		secondDiff := math.Floor(expTime.Sub(now).Seconds())

		client := &http.Client{}
		req, _ := http.NewRequest("GET", os.Getenv("AUTH_SERVICE_URL")+"/profile", nil)
		req.Header.Set("Authorization", token)

		resp, errResp := client.Do(req)
		if errResp != nil {
			log.Fatal(errResp)
		}
		defer resp.Body.Close()
		bodyText, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != http.StatusOK {
			response := errorresponse{
				Message: "Missing JWT",
			}

			return c.JSON(http.StatusUnauthorized, response)
		}

		// Parse return to struct
		userRespAuthService := domain.UserResponseAuthService{}
		json.Unmarshal([]byte(bodyText), &userRespAuthService)

		// Set struct to byte
		userModel, _ := json.Marshal(userRespAuthService.Profile)
		// Set to redis broker
		err := redisConn.Set(ctx, uuid, userModel, time.Duration(secondDiff)*time.Second).Err()

		if err != nil {
			response := errorresponse{
				Message: err.Error(),
			}

			return c.JSON(http.StatusUnauthorized, response)
		}

		c.Set("user", userRespAuthService.Profile)

		return next(c)
	}
}

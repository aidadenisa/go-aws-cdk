package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/golang-jwt/jwt/v5"
)

// Extract the request headers
// extracting our claims
// validating everything

func ValidateJWTMiddleware(next func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error)) func(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return func (request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

		// extract token from headers
		tokenString := extractTokenFromHeaders(request.Headers)
		if tokenString == "" {
			return events.APIGatewayProxyResponse{
				Body: "Missing auth tokens",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		// parse the token for the claims
		claims, err := parseToken(tokenString)
		if err != nil {
			return events.APIGatewayProxyResponse{
				Body: "User unauthorized",
				StatusCode: http.StatusUnauthorized,
			}, err
		}

		expires := int64(claims["expires"].(float64))

		if time.Now().Unix() > expires {
			// expired => token invalid
			return events.APIGatewayProxyResponse{
				Body: "Token expired",
				StatusCode: http.StatusUnauthorized,
			}, nil
		}

		return next(request)
	}
}

func extractTokenFromHeaders(headers map[string]string) string {
	authHeader, ok := headers["Authorization"]
	if !ok {
		return ""
	}

	splitToken := strings.Split(authHeader, "Bearer ")
	if len(splitToken) != 2 {
		return ""
	}

	return splitToken[1]
}

func parseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error){
		secret := "secret"

		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("unauthorized") 
	}

	if !token.Valid {
		return nil, fmt.Errorf("Token is not valid")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("Claims of unauthorized type")
	}

	return claims, nil
}
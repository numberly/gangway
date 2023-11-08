package jwt

import (
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/golang-jwt/jwt"
	"github.com/numberly/gangway/internal/config"
)

type ExtendedClaims struct {
	jwt.StandardClaims
	ClusterName     string `json:"cluster_name"`
	OAuth2State     string `json:"oauth2_state"`
	OAuth2TokenId   string `json:"oauth2_tokenid"`
	OAuth2RefreshId string `json:"oauth2_refreshid"`
}

// CreateToken creates a new JWT token with the given custom claims.
func CreateToken(clusterName, oauth2State string, clusterCfg config.MultiClusterConfig) (string, error) {
	// Define the expiration time of the token
	// Here, we have set it to 1 hours.
	expTime := time.Now().Add(1 * time.Hour).Unix()

	// Create the custom claims
	claims := ExtendedClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expTime,
			Issuer:    "Gangway",
		},
		ClusterName:     clusterName,
		OAuth2State:     oauth2State,
		OAuth2TokenId:   "",
		OAuth2RefreshId: "",
	}

	// Create a new token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	signedToken, err := token.SignedString([]byte(clusterCfg.SessionSecurityKey))
	if err != nil {
		log.Errorf("failed to sign the token: %v", err)
		return "", err
	}

	return signedToken, nil
}

// ValidateToken checks if the token is valid and returns the custom claims.
func ValidateToken(signedToken string, clusterCfg config.MultiClusterConfig) (*ExtendedClaims, error) {
	token, err := jwt.ParseWithClaims(signedToken, &ExtendedClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what you expect
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			log.Errorf("unexpected signing method: %v", token.Header["alg"])
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(clusterCfg.SessionSecurityKey), nil
	})

	if err != nil {
		log.Errorf("failed to parse token: %v", err)
		return nil, err
	}

	if claims, ok := token.Claims.(*ExtendedClaims); ok && token.Valid {
		return claims, nil
	} else {
		log.Errorf("invalid token or claims")
		return nil, err
	}
}

func UpdateToken(signedToken string, tokenId string, refreshTokenId string, clusterCfg config.MultiClusterConfig) (string, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(signedToken, &ExtendedClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("YourSecretKey"), nil
	})
	if err != nil {
		log.Errorf("failed to parse token: %v", err)
		return "", err
	}

	// Check if the token is valid and get the claims
	claims, ok := token.Claims.(*ExtendedClaims)
	if !ok || !token.Valid {
		log.Errorf("invalid token or claims")
		return "", err
	}

	claims.OAuth2TokenId = tokenId
	claims.OAuth2RefreshId = refreshTokenId

	// Create a new token with the updated claims
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedNewToken, err := newToken.SignedString([]byte(clusterCfg.SessionSecurityKey))
	if err != nil {
		log.Errorf("failed to sign the updated token: %v", err)
		return "", err
	}

	return signedNewToken, nil
}

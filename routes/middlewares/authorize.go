package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/StarWarsDev/legion-ops/internal/orm/models/user"

	"github.com/StarWarsDev/legion-ops/internal/data"

	"github.com/dgrijalva/jwt-go"
)

type Auth0OpenID struct {
	Issuer                            string   `json:"issuer"`
	AuthorizationEndpoint             string   `json:"authorization_endpoint"`
	TokenEndpoint                     string   `json:"token_endpoint"`
	UserinfoEndpoint                  string   `json:"userinfo_endpoint"`
	MfaChallengeEndpoint              string   `json:"mfa_challenge_endpoint"`
	JwksURI                           string   `json:"jwks_uri"`
	RegistrationEndpoint              string   `json:"registration_endpoint"`
	RevocationEndpoint                string   `json:"revocation_endpoint"`
	ScopesSupported                   []string `json:"scopes_supported"`
	ResponseTypesSupported            []string `json:"response_types_supported"`
	CodeChallengeMethodsSupported     []string `json:"code_challenge_methods_supported"`
	ResponseModesSupported            []string `json:"response_modes_supported"`
	SubjectTypesSupported             []string `json:"subject_types_supported"`
	IDTokenSigningAlgValuesSupported  []string `json:"id_token_signing_alg_values_supported"`
	TokenEndpointAuthMethodsSupported []string `json:"token_endpoint_auth_methods_supported"`
	ClaimsSupported                   []string `json:"claims_supported"`
	RequestURIParameterSupported      bool     `json:"request_uri_parameter_supported"`
	DeviceAuthorizationEndpoint       string   `json:"device_authorization_endpoint"`
}

type Auth0JWKS struct {
	Keys []struct {
		Alg string   `json:"alg"`
		Kty string   `json:"kty"`
		Use string   `json:"use"`
		N   string   `json:"n"`
		E   string   `json:"e"`
		Kid string   `json:"kid"`
		X5T string   `json:"x5t"`
		X5C []string `json:"x5c"`
	} `json:"keys"`
}

type ctxKey struct {
	name string
}

var userCtxKey = &ctxKey{name: "user"}

func (f *MiddlewareFuncs) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var dbUser user.User

		// get the auth token from the request headers
		authHeader := r.Header.Get("Authorization")
		authHeader = strings.Replace(authHeader, "Bearer ", "", 1)

		// if the token is present, decode it
		if authHeader != "" {
			token, err := jwt.Parse(authHeader, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}

				auth0Domain := os.Getenv("AUTH0_DOMAIN")
				auth0ClientID := os.Getenv("AUTH0_CLIENT_ID")

				checkAud := token.Claims.(jwt.MapClaims).VerifyAudience(auth0ClientID, false)
				if !checkAud {
					return nil, errors.New("invalid audience")
				}

				checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(fmt.Sprintf("https://%s/", auth0Domain), false)
				if !checkIss {
					return nil, errors.New("invalid issuer")
				}

				key, err := getAuth0JWKS(auth0Domain, token.Header["kid"].(string))
				if err != nil {
					return nil, err
				}

				cert := "-----BEGIN CERTIFICATE-----\n" + key + "\n-----END CERTIFICATE-----"
				result, _ := jwt.ParseRSAPublicKeyFromPEM([]byte(cert))

				return result, nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				// if valid, look up user in DB and set them in context
				dbUser, err = data.FindUserWithUsername(fmt.Sprintf("%v", claims["nickname"]), data.NewDB(f.dbORM))
				if err != nil {
					log.Println(err)
					next.ServeHTTP(w, r)
					return
				}

				// if the user doesn't have a picture in the database use the one in the claims
				if dbUser.Picture == "" {
					dbUser.Picture = claims["picture"].(string)
				}
			} else {
				log.Println(err)
				next.ServeHTTP(w, r)
				return
			}
		}
		ctx := context.WithValue(r.Context(), userCtxKey, &dbUser)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func getAuth0JWKS(domain, kid string) (string, error) {
	openID, err := getAuth0OpenID(domain)
	if err != nil {
		return "", err
	}

	jwksURI := openID.JwksURI
	resp, err := http.Get(jwksURI)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var jwks Auth0JWKS
	err = json.Unmarshal(body, &jwks)
	if err != nil {
		return "", err
	}

	var x5c string
	for _, key := range jwks.Keys {
		if key.Kid == kid && len(key.X5C) == 1 {
			x5c = key.X5C[0]
		}
	}

	if x5c == "" {
		return "", fmt.Errorf("no key found for kid [%s]", kid)
	}

	return x5c, nil
}

func getAuth0OpenID(domain string) (Auth0OpenID, error) {
	var auth0OpenID Auth0OpenID

	openIDEndpoint := fmt.Sprintf("https://%s/.well-known/openid-configuration", domain)
	resp, err := http.Get(openIDEndpoint)
	if err != nil {
		return auth0OpenID, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return auth0OpenID, err
	}

	err = json.Unmarshal(body, &auth0OpenID)
	if err != nil {
		return auth0OpenID, err
	}

	return auth0OpenID, nil
}

// UserInContext finds the user from the context. REQUIRES Middleware to have run.
func UserInContext(ctx context.Context) *user.User {
	raw, _ := ctx.Value(userCtxKey).(*user.User)
	return raw
}

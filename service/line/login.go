package line

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/golog"

	"github.com/jinil-ha/blind-msg/service"
	"github.com/jinil-ha/blind-msg/utils/config"
	"github.com/jinil-ha/blind-msg/utils/token"
)

var authorizeURL string
var tokenURL string

//var profileURL string

var loginChannelID string
var loginChannelSecret string
var redirectURI string

const stateLength = 32
const nonceLength = 32

func init() {
	loginChannelID = config.GetString("line.login.channel_id")
	loginChannelSecret = config.GetString("line.login.channel_secret")
	redirectURI = config.GetString("line.login.redirect_uri")

	authorizeURL = "https://access.line.me/oauth2/v2.1/authorize"
	tokenURL = "https://api.line.me/oauth2/v2.1/token"
	//profileURL = "https://api.line.me/v2/profile"
}

// AuthorizeURL return url for LINE login
func AuthorizeURL() string {
	authURL, _ := url.Parse(authorizeURL)

	state := token.Generate(stateLength)
	nonce := token.Generate(nonceLength)

	params := url.Values{}
	params.Add("response_type", "code")
	params.Add("client_id", loginChannelID)
	params.Add("redirect_uri", redirectURI)
	params.Add("state", state)
	params.Add("scope", "openid profile")
	params.Add("nonce", nonce)
	params.Add("bot_prompt", "aggressive")
	authURL.RawQuery = params.Encode()

	return authURL.String()
}

// tokenType is struct of token from LINE Login
type lineTokenType struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	IDToken      string `json:"id_token"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

func getToken(code string, token *lineTokenType) error {
	v := url.Values{}
	v.Add("grant_type", "authorization_code")
	v.Add("code", code)
	v.Add("redirect_uri", redirectURI)
	v.Add("client_id", loginChannelID)
	v.Add("client_secret", loginChannelSecret)
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(v.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		golog.Debugf("response code : %d", resp.StatusCode)
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	golog.Debugf("token response : %s", body)

	err = json.Unmarshal(body, token)
	if err != nil {
		return err
	}

	return nil
}

// Auth implements LINE auth.
// param is param info redirected from LINE.
// and use ID Token as token
func Auth(param map[string]string, token *string) error {
	code := param["code"]
	state := param["state"]
	changed := param["friendship_status_changed"]
	error := param["error"]
	desc := param["error_description"]

	if error != "" {
		golog.Warnf("LINE login error(state:%s): %s [%s]", state, error, desc)
		return fmt.Errorf("%s %s", error, desc)
	}
	golog.Infof("LINE login succeeded: code[%s] state[%s] friendship_status_changed[%s]",
		code, state, changed)

	// get line token
	var tk lineTokenType
	err := getToken(code, &tk)
	if err != nil {
		return err
	}

	// post process
	*token = tk.IDToken
	return nil
}

// Claims is ID Token(JWT)'s payload type
type Claims struct {
	AuthTime int    `json:"auth_time"`
	Nonce    string `json:"nonce"`
	Name     string `json:"name"`
	Picture  string `json:"picture"`

	jwt.StandardClaims
}

// GetProfile get user's profile info from token(ID Token)
func GetProfile(token string, profile *service.ProfileType) error {
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(loginChannelSecret), nil
	})
	if err != nil {
		return fmt.Errorf("parsing error: %s", err)
	}
	if !tkn.Valid {
		return fmt.Errorf("id token invalid: %s", err)
	}

	profile.UserID = claims.Subject
	profile.DisplayName = claims.Name
	profile.PictureURL = claims.Picture

	return nil
}

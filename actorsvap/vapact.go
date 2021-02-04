package actorsvap

// https://yourbasic.org/golang/iota/
// https://stackoverflow.com/questions/27236827/idiomatic-way-to-make-a-request-response-communication-using-channels

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/va-voice-gateway/appconfig"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	instance *vapActor
	once sync.Once
)

const DURATION_23_HOURS = 23 * 60 * 60

type VapToken struct {
	AccessToken    string `json:"accessToken"`
	Authentication VapTokenAuthentication
	User           VapTokenUser
}

type VapTokenAuthentication struct {
	strategy string
}

type VapTokenUser struct {
	AccessToken string `json:"accessToken"`
	Email string `json:"email"`
	Description string `json:"description"`
	AllowedServices []string `json:"allowedServices"`
}

type VapTokenCacheEntry struct {
	CurrentToken VapToken
	CreatedTime  time.Time
}

type VapTokenCache struct {
	SvcAccUsr string
	SvcAccPwd string
	VapBaseUrl string
	Token *VapTokenCacheEntry
}

func NewVapTokenCache(SvcAccUsr string, SvcAccPwd string, VapBaseUrl string) VapTokenCache {
	return VapTokenCache{
		SvcAccUsr: SvcAccUsr,
		SvcAccPwd: SvcAccPwd,
		VapBaseUrl: VapBaseUrl,
		Token: nil,
	}
}

func (c *VapTokenCache) GetNewToken() (*VapToken, error) {
	reqBody := fmt.Sprintf(`
		{
    		"strategy": "local",
    		"email": "%s",
    		"password": "%s"
    	}
	`, c.SvcAccUsr, c.SvcAccPwd)
	log.Println("GetNewToken.reqBody ",reqBody)

	url := fmt.Sprintf("%s%s", c.VapBaseUrl, "/vapapi/authentication/v1")
	resp, err := http.Post(url, "application/json", bytes.NewBuffer([]byte(reqBody)))
	defer resp.Body.Close()

	if err != nil {
		log.Printf("GetNewToken: error when calling  /vapapi/authentication/v1: %v\n", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GetNewToken: error when reading http response: %v\n", err)
		return nil, err
	}

	vapToken := &VapToken{}
	err = json.Unmarshal(body, vapToken)
	if err != nil {
		log.Printf("GetNewToken: error when parsing json: %v\n", err)
		return nil, err
	}

	return vapToken, nil
}

func (c *VapTokenCache) GetToken() (*string, error) {

	if c.Token == nil {
		token, err := c.GetNewToken()
		if err != nil {
			log.Printf("GetToken: error when calling GetNewToken: %v\n", err)
			return nil, err
		}
		c.Token = &VapTokenCacheEntry{
			CurrentToken: *token,
			CreatedTime: time.Now(),
		}
		return &c.Token.CurrentToken.AccessToken, nil
	} else {
		token := c.Token
		if IsExpired(time.Now(), token.CreatedTime, DURATION_23_HOURS) {
			token, err := c.GetNewToken()
			if err != nil {
				log.Printf("GetToken: error when calling GetNewToken: %v\n", err)
				return nil, err
			}
			c.Token = &VapTokenCacheEntry{
				CurrentToken: *token,
				CreatedTime: time.Now(),
			}
			return &c.Token.CurrentToken.AccessToken, nil
		} else {
			return &c.Token.CurrentToken.AccessToken, nil
		}
	}

}

func IsExpired(Now time.Time, CreatedTime time.Time, AllowedTokenAge int64) bool {
	if Now.Unix() - CreatedTime.Unix() < 0 {
		// token creation specified in the future? -> return true (expired)
		return true
	}
	return Now.Unix() - CreatedTime.Unix() > AllowedTokenAge
}

type VapTokenRequest struct {
	Responder chan string
}

type vapActor struct {
	CommandsChannel chan VapTokenRequest
	VapTokenCache   VapTokenCache
}
func VapActor() *vapActor {
	once.Do(func() {
		SvcAccUsr := appconfig.AppConfig(nil).NlpVap.Username
		SvcAccPwd := appconfig.AppConfig(nil).NlpVap.Password
		VapBaseUrl := appconfig.AppConfig(nil).NlpVap.VapBaseUrl

		cache := NewVapTokenCache(SvcAccUsr, SvcAccPwd, VapBaseUrl)
		chnl := make(chan VapTokenRequest)
		instance = &vapActor{
			CommandsChannel: chnl,
			VapTokenCache: cache,
		}
	})
	return instance
}

func (va *vapActor) VapActorProcessingLoop() {
	for command := range va.CommandsChannel {
		token, err:= va.VapTokenCache.GetToken()
		if err != nil {
			log.Printf("VapTokenRequest processing error %v\n", err)
		} else {
			command.Responder <- *token
		}
	}
}

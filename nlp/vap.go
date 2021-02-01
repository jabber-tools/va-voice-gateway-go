package nlp

// https://yourbasic.org/golang/iota/
// https://stackoverflow.com/questions/27236827/idiomatic-way-to-make-a-request-response-communication-using-channels

import (
	"encoding/json"
	"fmt"
	"github.com/va-voice-gateway/gateway/config"
	"io/ioutil"
	"log"
	"time"
	"net/http"
	"bytes"
)

const DURATION_23_HOURS = 23 * 60 * 60

type VapToken struct {
	AccessToken string `json:"accessToken"`
	Authentication VapTokenAuthentication
	User VapTokenUser
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
	CreatedTime time.Time
}

type VapTokenCache struct {
	SvcAccUsr string
	SvcAccPwd string
	VapBaseUrl string
	Token *VapTokenCacheEntry
}

func NewVapTokenCache(SvcAccUsr string, SvcAccPwd string, VapBaseUrl string) VapTokenCache {
	return VapTokenCache {
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
		c.Token = &VapTokenCacheEntry {
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
			c.Token = &VapTokenCacheEntry {
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

type VapActor struct {
	CommandsChannel chan VapTokenRequest
	VapTokenCache VapTokenCache
}
func NewVapActor(SvcAccUsr string, SvcAccPwd string, VapBaseUrl string) VapActor {
	cache :=NewVapTokenCache(SvcAccUsr, SvcAccPwd, VapBaseUrl)
	chnl := make(chan VapTokenRequest)
	return VapActor {
		CommandsChannel: chnl,
		VapTokenCache: cache,
	}
}

func (va *VapActor) VapActorProcessingLoop() {
	for command := range va.CommandsChannel {
		token, err:= va.VapTokenCache.GetToken()
		if err != nil {
			log.Printf("VapTokenRequest processing error %v\n", err)
		} else {
			command.Responder <- *token
		}
	}
}


// for now just taking from file and parsing
// in the future VAP API will be called:
// /vapapi/vap-mgmt/config-mgmt/v1?voiceEnabled=1
func GetBotConfigs() ([]config.BotConfig, error) {

	var botConfigs []config.BotConfig

	content, err := ioutil.ReadFile("c:/tmp/botconfigs.json")
	if err != nil {
		log.Printf("GetBotConfigs ReadFile error: %v", err)
		return nil, err
	}

	data := []byte(content)

	err = json.Unmarshal(data, &botConfigs)

	if err != nil {
		log.Printf("GetBotConfigs Unmarshal error: %v", err)
		return nil, err
	}

	return botConfigs, nil
}

// this is target replacement of GetBotConfigs above
func (va *VapActor) GetBotConfigsFromVap() ([]config.BotConfig, error) {
	c := make(chan string)
	request := VapTokenRequest {Responder: c}
	va.CommandsChannel <- request
	token := <- c

	url := fmt.Sprintf("%s%s", va.VapTokenCache.VapBaseUrl, "/vapapi/vap-mgmt/config-mgmt/v1?voiceEnabled=1")

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", token)
	resp, err := client.Do(req)
	defer resp.Body.Close()

	if err != nil {
		log.Printf("GetBotConfigsFromVap: error when calling  /vapapi/vap-mgmt/config-mgmt/v1: %v\n", err)
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("GetBotConfigsFromVap: error when reading http response: %v\n", err)
		return nil, err
	}

	botConfigs := make([]config.BotConfig, 0)
	err = json.Unmarshal(body, &botConfigs)
	if err != nil {
		log.Printf("GetBotConfigsFromVap: error when parsing json: %v\n", err)
		return nil, err
	}
	return botConfigs, nil
}

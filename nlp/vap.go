package nlp

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
	`, c.SvcAccUsr, c.VapBaseUrl)
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

	vapToken := VapToken{}
	err = json.Unmarshal(body, vapToken)
	if err != nil {
		log.Printf("GetNewToken: error when parsing json: %v\n", err)
		return nil, err
	}

	return &vapToken, nil
}

func (c *VapTokenCache) GetToken() (*VapToken, error) {
	// TBD: implement proper retrieval based on expiration TS
	return c.GetNewToken()
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

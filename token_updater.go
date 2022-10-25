package selling_partner_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fond-of-vertigo/logger"
	"io"
	"net/http"
	"sync/atomic"
	"time"
)

type TokenUpdaterConfig struct {
	RefreshToken string
	ClientID     string
	ClientSecret string
	Logger       logger.Logger
	QuitSignal   chan bool
}

type TokenUpdaterInterface interface {
	GetAccessToken() string
}

type TokenUpdater struct {
	accessToken     *atomic.Value
	ExpireTimestamp *atomic.Int64
	config          TokenUpdaterConfig
	log             logger.Logger
	quitSignal      chan bool
}
type AccessTokenResponse struct {
	AccessToken      string `json:"access_token"`
	RefreshToken     string `json:"refresh_token"`
	ExpiresIn        int    `json:"expires_in"`
	TokenType        string `json:"token_type"`
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func NewTokenUpdater(config TokenUpdaterConfig) (*TokenUpdater, error) {
	t := TokenUpdater{
		config:     config,
		log:        config.Logger,
		quitSignal: config.QuitSignal,
	}
	if err := t.fetchNewToken(); err != nil {
		return nil, fmt.Errorf("accesstoken could not be fetched: %w", err)
	}
	return &t, nil
}

func (t *TokenUpdater) RunInBackground() {
	go t.checkAccessToken()
}

func (t *TokenUpdater) checkAccessToken() {
	for {
		select {
		case <-t.quitSignal:
			t.log.Infof("Received signal to stop token updates.")
			return
		default:
			secondsToWait := secondsUntilExpired(t.ExpireTimestamp.Load())
			if secondsToWait <= int64(ExpiryDelta.Seconds()) {
				if err := t.fetchNewToken(); err != nil {
					t.log.Errorf(err.Error())
				}
			} else {
				time.Sleep(time.Duration(secondsToWait-int64(ExpiryDelta.Seconds())) * time.Second)
			}
		}
	}
}

func (t *TokenUpdater) GetAccessToken() string {
	return fmt.Sprintf("%v", t.accessToken.Load())
}

func (t *TokenUpdater) fetchNewToken() error {
	reqBody, _ := json.Marshal(map[string]string{
		"grant_type":    "refresh_token",
		"refresh_token": t.config.RefreshToken,
		"client_id":     t.config.ClientID,
		"client_secret": t.config.ClientSecret,
	})

	resp, err := http.Post(
		"https://api.amazon.com/auth/o2/token",
		"application/json",
		bytes.NewBuffer(reqBody))

	if err != nil {
		return err
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			t.log.Errorf(err.Error())
		}
	}(resp.Body)

	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	parsedResp := &AccessTokenResponse{}

	if err = json.Unmarshal(respBodyBytes, parsedResp); err != nil {
		return fmt.Errorf("RefreshToken response parse failed. Body: %s", string(respBodyBytes))
	}
	if parsedResp.AccessToken != "" {
		t.accessToken.Swap(parsedResp.AccessToken)

		expireTimestamp := time.Now().UTC().Add(time.Duration(parsedResp.ExpiresIn) * time.Second)
		t.ExpireTimestamp.Swap(expireTimestamp.Unix())
	}
	return nil
}

func secondsUntilExpired(unixTimestamp int64) int64 {
	currentTimestamp := time.Now().Unix()
	return unixTimestamp - currentTimestamp
}

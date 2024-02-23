package model

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"galxe/utils"
	http "github.com/bogdanfinn/fhttp"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"strings"
	"time"
)

func (galxe *Galxe) Login() error {
	for i := 0; i < galxe.config.Info.MaxTasksRetries; i++ {
		err := galxe.loginProcess()
		if err != nil {
			continue
		}

		exists, err := galxe.checkIfAccountExists()
		if err != nil {
			continue
		}

		if exists == false {
			err = galxe.createNewAccount()
			if err != nil {
				continue
			} else {
				return nil
			}
		} else {
			return nil
		}
	}
	return errors.New("")
}

func (galxe *Galxe) loginProcess() error {
	messageToSign, hexSignature := galxe.generateLoginSignature()
	if messageToSign == "" {
		return errors.New("")
	}

	data := map[string]interface{}{
		"operationName": "SignIn",
		"variables": map[string]interface{}{
			"input": map[string]interface{}{
				"address":     galxe.walletAddress,
				"message":     messageToSign,
				"signature":   hexSignature,
				"addressType": "EVM",
			},
		},
		"query": "mutation SignIn($input: Auth) {\n  signin(input: $input)\n}\n",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		galxe.logger.Error("%d | Failed in marshaling login signature: %s", galxe.index, err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://graphigo.prd.galaxy.eco/query", bytes.NewBuffer(jsonData))
	if err != nil {
		galxe.logger.Error("%d | Error in creating SignIn request: %s", galxe.index, err)
		return err
	}
	req.Header.Set("accept", "*/*")
	req.Header.Set("cookie", galxe.cookieClient.CookiesToHeader())
	req.Header.Set("content-type", "application/json")
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("origin", "https://galxe.com")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")
	resp, err := galxe.client.Do(req)
	if err != nil {
		galxe.logger.Error("%d | Failed to SignIn: %s", galxe.index, err)
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		galxe.logger.Error("%d | Failed to read SignIn body: %s", galxe.index, err)
		return err
	}
	galxe.cookieClient.SetCookieFromResponse(resp)
	var response signInResponse
	if err := json.Unmarshal(body, &response); err != nil {
		galxe.logger.Error("%d | Failed to unmarshal SignIn response: %s", galxe.index, string(body))
		return err
	}

	galxe.jwtToken = response.Data.AuthKey

	galxe.cookieClient.AddCookies([]http.Cookie{
		{
			Name:     "authorization",
			Value:    galxe.jwtToken,
			Path:     "/",
			Domain:   "galxe.com",
			Secure:   false,
			HttpOnly: false,
			SameSite: 0,
		},
	})

	galxe.logger.Success("%d | Logged into Galxe.", galxe.index)
	return nil
}

// createNewAccount sends wallet address and username if success updates galxeAccountID
func (galxe *Galxe) createNewAccount() error {
	username, err := galxe.checkIfUsernameExists()
	if username == "" || err != nil {
		galxe.logger.Error("%d | Failed to create new username.", galxe.index)
		return errors.New("")
	}

	data := map[string]interface{}{
		"operationName": "CreateNewAccount",
		"variables": map[string]interface{}{
			"input": map[string]interface{}{
				"schema":         fmt.Sprintf("EVM:%s", galxe.walletAddress),
				"socialUsername": "",
				"username":       username,
			},
		},
		"query": "mutation CreateNewAccount($input: CreateNewAccount!) {\n  createNewAccount(input: $input)\n}\n",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		galxe.logger.Error("%d | Error in marshaling json createNewAccount: %s", galxe.index, err)
		return err
	}

	req, err := http.NewRequest(http.MethodPost, "https://graphigo.prd.galaxy.eco/query", bytes.NewBuffer(jsonData))
	if err != nil {
		galxe.logger.Error("%d | Error in creating createNewAccount request: %s", galxe.index, err)
		return err
	}

	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("authorization", galxe.jwtToken)
	req.Header.Set("cookie", galxe.cookieClient.CookiesToHeader())
	req.Header.Set("content-type", "application/json")
	req.Header.Set("accept", "*/*")
	req.Header.Set("sec-ch-ua-platform", `"Windows"`)
	req.Header.Set("origin", "https://galxe.com")
	req.Header.Set("sec-fetch-site", "cross-site")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-dest", "empty")

	resp, err := galxe.client.Do(req)
	if err != nil {
		galxe.logger.Error("%d | Failed to createNewAccount: %s", galxe.index, err)
		return err
	}

	defer resp.Body.Close()
	galxe.cookieClient.SetCookieFromResponse(resp)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		galxe.logger.Error("%d | Failed to read createNewAccount body: %s", galxe.index, err)
		return err
	}

	var response createNewAccountResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		galxe.logger.Error("Failed to unmarshal createNewAccount response: %s", err)
		return err
	}

	if response.Data.CreateNewAccount == "" {
		galxe.logger.Error("%d | Wrong CreateAccount response: %s", galxe.index, string(body))
		return errors.New("")
	}

	galxe.galxeAccountID = response.Data.CreateNewAccount
	galxe.logger.Success("%d | New Galxe account created. Username -> %s", galxe.index, username)
	return nil
}

func (galxe *Galxe) checkIfAccountExists() (bool, error) {
	data := map[string]interface{}{
		"operationName": "GalxeIDExist",
		"variables": map[string]interface{}{
			"schema": fmt.Sprintf("EVM:%s", galxe.walletAddress),
		},
		"query": "query GalxeIDExist($schema: String!) {\n  galxeIdExist(schema: $schema)\n}\n",
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		galxe.logger.Error("%d | Error in marshaling json ifAccExists: %s", galxe.index, err)
		return false, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://graphigo.prd.galaxy.eco/query", bytes.NewBuffer(jsonData))
	if err != nil {
		galxe.logger.Error("%d | Error in creating ifAccExists request: %s", galxe.index, err)
		return false, err
	}
	req.Header = http.Header{
		"accept":             {"*/*"},
		"content-type":       {"application/json"},
		"cookie":             {galxe.cookieClient.CookiesToHeader()},
		"sec-ch-ua-mobile":   {"?0"},
		"user-agent":         {"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
		"sec-ch-ua-platform": {`"Windows"`},
		"origin":             {"https://galxe.com"},
		"sec-fetch-site":     {"cross-site"},
		"sec-fetch-mode":     {"cors"},
		"sec-fetch-dest":     {"empty"},
	}
	resp, err := galxe.client.Do(req)
	if err != nil {
		galxe.logger.Error("%d | Failed to check if account exists: %s", galxe.index, err)
	}

	defer resp.Body.Close()
	galxe.cookieClient.SetCookieFromResponse(resp)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		galxe.logger.Error("%d | Failed to read ifAccExists body: %s", galxe.index, err)
		return false, err
	}

	var response ifExistsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		if err != nil {
			galxe.logger.Error("Failed to unmarshal ifAccountExists response: %s", err)
		}
		return false, err
	}

	return response.Data.GalxeIDExist, nil
}

func (galxe *Galxe) checkIfUsernameExists() (string, error) {
	for i := 0; i < 10; i++ {
		username := utils.GenerateRandomUsername()

		data := map[string]interface{}{
			"operationName": "IsUsernameExisting", "variables": map[string]interface{}{
				"username": username,
			},
			"query": "query IsUsernameExisting($username: String!) {\n  usernameExist(username: $username)\n}\n",
		}
		jsonData, err := json.Marshal(data)
		if err != nil {
			galxe.logger.Error("%d | Error in marshaling json checkIfUsernameExists: %s", galxe.index, err)
			return "", err
		}

		req, err := http.NewRequest(http.MethodPost, "https://graphigo.prd.galaxy.eco/query", bytes.NewBuffer(jsonData))
		if err != nil {
			galxe.logger.Error("%d | Error in creating checkIfUsernameExists request: %s", galxe.index, err)
			return "", err
		}
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("authorization", galxe.jwtToken)
		req.Header.Set("cookie", galxe.cookieClient.CookiesToHeader())
		req.Header.Set("content-type", "application/json")
		req.Header.Set("accept", "*/*")
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		req.Header.Set("origin", "https://galxe.com")
		req.Header.Set("sec-fetch-site", "cross-site")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-dest", "empty")
		resp, err := galxe.client.Do(req)
		if err != nil {
			galxe.logger.Error("%d | Failed to check if username exists: %s", galxe.index, err)
			return "", err
		}
		galxe.cookieClient.SetCookieFromResponse(resp)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			galxe.logger.Error("%d | Failed to read checkIfUsernameExists body: %s", galxe.index, err)
			return "", err
		}
		resp.Body.Close()

		var response ifUsernameExistsResponse
		if err := json.Unmarshal(body, &response); err != nil {
			galxe.logger.Error("%d | Failed to unmarshal checkIfUsernameExists response: %s", galxe.index, err)
			return "", err
		}

		if response.Data.UsernameExists {
			galxe.logger.Warning("%d | Username %s already exists. Trying next one...", galxe.index, username)
			continue
		} else {
			return username, nil
		}
	}
	return "", errors.New("")
}

func (galxe *Galxe) generateLoginSignature() (string, string) {
	now := time.Now().UTC()
	nowFormatted := now.Format("2006-01-02T15:04:05.000Z")
	sevenDaysLater := now.Add(time.Hour * 24 * 7)
	sevenDaysLaterFormatted := sevenDaysLater.Format("2006-01-02T15:04:05.000Z")

	nonce := utils.GenerateRandomString(17)

	messageToSign := fmt.Sprintf("galxe.com wants you to sign in with your Ethereum account:\n%s\n\nSign in with Ethereum to the app.\n\nURI: https://galxe.com\nVersion: 1\nChain ID: 1\nNonce: %s\nIssued At: %s\nExpiration Time: %s", strings.ToLower(galxe.walletAddress), nonce, nowFormatted, sevenDaysLaterFormatted)

	privateKey, err := crypto.HexToECDSA(galxe.privateKey)
	if err != nil {
		galxe.logger.Error("%d | Unable to convert private key to ECDSA: %s", galxe.index, err)
		return "", ""
	}

	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(messageToSign), messageToSign)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		galxe.logger.Error("%d | Failed to sign login message: %s", galxe.index, err)
		return "", ""
	}

	// Convert to Ethereum signature format
	signature[64] += 27
	hexSignature := fmt.Sprintf("0x%x", signature)

	return messageToSign, hexSignature
}

type createNewAccountResponse struct {
	Data struct {
		CreateNewAccount string `json:"createNewAccount"`
	} `json:"data"`
}

type ifExistsResponse struct {
	Data struct {
		GalxeIDExist bool `json:"galxeIdExist"`
	} `json:"data"`
}

type ifUsernameExistsResponse struct {
	Data struct {
		UsernameExists bool `json:"usernameExist"`
	} `json:"data"`
}

type signInResponse struct {
	Data struct {
		AuthKey string `json:"signin"`
	} `json:"data"`
}

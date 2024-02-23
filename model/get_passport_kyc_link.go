package model

import (
	"bytes"
	"encoding/json"
	"fmt"
	"galxe/model/additional_galxe_methods"
	http "github.com/bogdanfinn/fhttp"
	"github.com/ethereum/go-ethereum/crypto"
	"io"
	"strings"
)

func (galxe *Galxe) GetPassportKYCLink() string {
	passportKYCLink := ""
	for i := 0; i < galxe.config.Info.MaxTasksRetries; i++ {
		signature := galxe.generateKYCSignature()
		if signature == "" {
			continue
		}

		URL := "https://graphigo.prd.galaxy.eco/query"
		kycPassportData := map[string]interface{}{
			"operationName": "GetOrCreateInquiryByAddress",
			"variables": map[string]interface{}{
				"input": map[string]interface{}{
					"address":   galxe.walletAddress,
					"signature": signature,
				},
			},
			"query": additional_galxe_methods.PassportKYCQuery,
		}

		jsonData, err := json.Marshal(kycPassportData)
		if err != nil {
			galxe.logger.Error("%d | Error in marshaling json PassportKYC: %s", galxe.index, err)
			continue
		}

		req, err := http.NewRequest(http.MethodPost, URL, bytes.NewBuffer(jsonData))
		if err != nil {
			galxe.logger.Error("%d | Error in creating PassportKYC request: %s", galxe.index, err)
			continue
		}

		req.Header.Set("authority", "graphigo.prd.galaxy.eco")
		req.Header.Set("accept", "*/*")
		req.Header.Set("authorization", galxe.jwtToken)
		req.Header.Set("content-type", "application/json")
		req.Header.Set("origin", "https://galxe.com")
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"Windows"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "cross-site")

		resp, err := galxe.client.Do(req)
		if err != nil {
			galxe.logger.Error("%d | Failed to PassportKYC: %s", galxe.index, err)
			continue
		}

		defer resp.Body.Close()
		galxe.cookieClient.SetCookieFromResponse(resp)
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			galxe.logger.Error("%d | Failed to read PassportKYC body: %s", galxe.index, err)
			continue
		}

		var passportLink passportLinkResponse
		err = json.Unmarshal(body, &passportLink)
		if passportLink.Data.GetOrCreateInquiryByAddress.PersonaInquiry.InquiryID == "" {
			galxe.logger.Error("%d | Wrong PassportKycID response: %s", galxe.index, string(body))
			continue
		} else {
			galxe.logger.Success("%d | Got passport KYC link!", galxe.index)
			inquiryID := passportLink.Data.GetOrCreateInquiryByAddress.PersonaInquiry.InquiryID
			sessionToken := passportLink.Data.GetOrCreateInquiryByAddress.PersonaInquiry.SessionToken
			personaWidgetID := additional_galxe_methods.GeneratePersonaWidget()
			passportKYCLink = fmt.Sprintf("https://withpersona.com/widget?client-version=4.7.1&container-id=persona-widget-%s&flow-type=embedded&environment=production&iframe-origin=https://galxe.com&inquiry-id=%s&session-token=%s", personaWidgetID, inquiryID, sessionToken)
			return passportKYCLink
		}
	}
	return passportKYCLink
}

func (galxe *Galxe) generateKYCSignature() string {
	messageToSign := fmt.Sprintf("get_or_create_address_inquiry:%s", strings.ToLower(galxe.walletAddress))

	privateKey, err := crypto.HexToECDSA(galxe.privateKey)
	if err != nil {
		galxe.logger.Error("%d | Unable to convert private key to ECDSA: %s", galxe.index, err)
		return ""
	}

	prefixedMessage := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(messageToSign), messageToSign)
	hash := crypto.Keccak256Hash([]byte(prefixedMessage))

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		galxe.logger.Error("%d | Failed to sign KYC message: %s", galxe.index, err)
		return ""
	}

	hexSignature := fmt.Sprintf("0x%x", signature)

	return hexSignature
}

type passportLinkResponse struct {
	Data struct {
		GetOrCreateInquiryByAddress struct {
			Status         string `json:"status"`
			Vendor         string `json:"vendor"`
			PersonaInquiry struct {
				InquiryID      string `json:"inquiryID"`
				SessionToken   string `json:"sessionToken"`
				DeclinedReason string `json:"declinedReason"`
				Typename       string `json:"__typename"`
			} `json:"personaInquiry"`
			Typename string `json:"__typename"`
		} `json:"getOrCreateInquiryByAddress"`
	} `json:"data"`
}

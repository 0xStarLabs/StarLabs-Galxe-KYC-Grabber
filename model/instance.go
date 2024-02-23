package model

import (
	"crypto/ecdsa"
	"galxe/extra"
	"galxe/utils"
	tlsClient "github.com/bogdanfinn/tls-client"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
)

type Galxe struct {
	index         int
	proxy         string
	config        extra.Config
	privateKey    string
	walletAddress string

	jwtToken       string
	galxeAccountID string

	client       tlsClient.HttpClient
	cookieClient *utils.CookieClient
	logger       extra.Logger
}

func (galxe *Galxe) InitGalxe(index int, proxy, privateKey string, config extra.Config) bool {
	private := strings.Replace(privateKey, "0x", "", 1)
	privateKeyObj, err := crypto.HexToECDSA(private)
	if err != nil {
		galxe.logger.Error("%d | Wrong private key format: %s", index, err)
		return false
	}
	publicKey := privateKeyObj.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		galxe.logger.Error("%d | Error casting public key to ECDSA", index)
		return false
	}

	galxe.walletAddress = crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

	galxe.index = index
	galxe.proxy = proxy
	galxe.config = config
	galxe.privateKey = private

	return galxe.prepareClient()
}

func (galxe *Galxe) prepareClient() bool {
	var err error

	for i := 0; i < galxe.config.Info.MaxTasksRetries; i++ {
		galxe.cookieClient = utils.NewCookieClient()
		galxe.client, err = utils.CreateHttpClient(galxe.proxy)
		if err != nil {
			galxe.logger.Warning("%d | Failed to prepare client: %s", galxe.index, err)
			continue
		} else {
			return true
		}
	}

	galxe.logger.Error("%d | Failed to prepare client.", galxe.index)
	return false
}

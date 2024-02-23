package utils

import (
	"fmt"
	"galxe/extra"
	tlsClient "github.com/bogdanfinn/tls-client"
	"github.com/bogdanfinn/tls-client/profiles"
)

func CreateHttpClient(proxies string) (tlsClient.HttpClient, error) {

	options := []tlsClient.HttpClientOption{
		tlsClient.WithClientProfile(profiles.Chrome_120),
		tlsClient.WithNotFollowRedirects(),
		tlsClient.WithRandomTLSExtensionOrder(),
		tlsClient.WithInsecureSkipVerify(),
		tlsClient.WithTimeoutSeconds(30),
	}
	if proxies != "" {
		options = append(options, tlsClient.WithProxyUrl(fmt.Sprintf("http://%s", proxies)))
	}

	client, err := tlsClient.NewHttpClient(tlsClient.NewNoopLogger(), options...)
	if err != nil {
		extra.Logger{}.Error("Failed to create Http Client: %s", err)
		return nil, err
	}

	return client, nil
}

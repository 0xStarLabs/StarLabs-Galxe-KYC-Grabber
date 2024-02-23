package extra

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
	"log"
)

func ExtractPrivateKeyFromMnemo(mnemonic string) string {
	seed := bip39.NewSeed(mnemonic, "")

	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		Logger{}.Error("Failed to extract private key from mnemonic: %s", mnemonic)
		return ""
	}

	// Derive the path for the first account using BIP44
	// m/44'/60'/0'/0/0 is the derivation path for the first Ethereum address
	childKey, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		log.Fatal(err)
	}
	childKey, err = childKey.NewChildKey(bip32.FirstHardenedChild + 60)
	if err != nil {
		log.Fatal(err)
	}
	childKey, err = childKey.NewChildKey(bip32.FirstHardenedChild)
	if err != nil {
		log.Fatal(err)
	}
	childKey, err = childKey.NewChildKey(0)
	if err != nil {
		log.Fatal(err)
	}
	childKey, err = childKey.NewChildKey(0)
	if err != nil {
		log.Fatal(err)
	}

	privateKey, err := crypto.ToECDSA(childKey.Key)
	if err != nil {
		log.Fatal(err)
	}

	return "0x" + fmt.Sprintf("%x", crypto.FromECDSA(privateKey))
}

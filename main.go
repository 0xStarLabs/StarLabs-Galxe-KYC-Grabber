package main

import (
	"fmt"
	"galxe/extra"
	"galxe/model"
	"github.com/zenthangplus/goccm"
	"sync"
	"time"
)

func main() {
	extra.ShowLogo()
	extra.ShowDevInfo()
	options()
}

func options() {
	config := extra.ReadConfig()
	//fmt.Print("\n\033[H\033[2J")
	privateKeys := extra.ReadPrivateKeys()
	proxies := extra.ReadTxtFile("proxies", "data/proxies.txt")

	var ipChangeLinks []string
	var threads int
	var ok bool
	if config.Proxy.MobileProxy == "yes" {
		ipChangeLinks = extra.ReadTxtFile("ip change links", "data/ip_change_links.txt")
		threads = len(ipChangeLinks)
	} else {
		threads, ok = extra.UserInputInteger("How many concurrent threads do you want")
		if ok == false {
			return
		}
	}

	if len(proxies) == 0 {
		if extra.NoProxies() == false {
			return
		} else {
			for i := range privateKeys {
				proxies[i] = ""
			}
		}
	} else if config.Proxy.MobileProxy == "no" && len(proxies) < len(privateKeys) {
		newProxies := make([]string, len(privateKeys))
		for i := range privateKeys {
			newProxies[i] = proxies[i%len(proxies)]
		}
		proxies = newProxies
	}

	var mutex = &sync.Mutex{}
	var KYCLinks sync.Map
	for _, key := range privateKeys {
		KYCLinks.Store(key, "")
	}

	failedAccounts := make(chan int, len(privateKeys))
	indexes := make(chan int, len(privateKeys))
	for i := range privateKeys {
		indexes <- i
	}
	close(indexes)

	goroutines := goccm.New(threads)
	fmt.Println()
	// set up account range. use all accounts if accounts range set to 0-0 in config by default
	var start, end int
	if config.Info.AccountRange.Start == 0 && config.Info.AccountRange.End == 0 || config.Info.AccountRange.Start >= config.Info.AccountRange.End {
		start = 0
		end = len(privateKeys)
	} else {
		start = config.Info.AccountRange.Start - 1
		end = config.Info.AccountRange.End
	}

	if config.Proxy.MobileProxy == "yes" {
		for i, proxy := range proxies {
			data := MobileProxyData{
				Proxy:          proxy,
				IPChangeLink:   ipChangeLinks[i],
				Indexes:        indexes,
				PrivateKeys:    privateKeys,
				Config:         config,
				KYCLinks:       &KYCLinks,
				Mutex:          mutex,
				FailedAccounts: failedAccounts,
			}
			goroutines.Wait()
			go func(i int) {
				MobileProxyWrapper(data)
				goroutines.Done()
			}(i)
		}
	} else {
		for i := start; i < end && i < len(privateKeys); i++ {
			goroutines.Wait()
			go func(i int) {
				Process(i, proxies[i], privateKeys[i], config, &KYCLinks, mutex, failedAccounts)
				goroutines.Done()
				extra.RandomSleep(config.Info.PauseBetweenAccounts.Start, config.Info.PauseBetweenAccounts.End)
			}(i)
		}
	}
	goroutines.WaitAllDone()
	close(failedAccounts)

	extra.CreateXLSXFromSyncMap(&KYCLinks, privateKeys, "data/collected_links.xlsx")
	fmt.Println()
	var failedCounter int

	for _ = range failedAccounts {
		failedCounter++
	}
	extra.Logger{}.Info("STATISTICS: %d SUCCESS | %d FAILED", end-start-failedCounter, failedCounter)
	// exit with Enter
	fmt.Println("Press Enter to exit...")
	_, err := fmt.Scanln()
	if err != nil {
		return
	}
}

func Process(index int, proxy, privateKey string, config extra.Config, KYCLinks *sync.Map, mutex *sync.Mutex, failedAccounts chan<- int) {
	session := model.Galxe{}
	session.InitGalxe(index+1, proxy, privateKey, config)

	err := session.Login()
	if err != nil {
		failedAccounts <- index
		return
	}

	KYCLink := session.GetPassportKYCLink()

	if KYCLink == "" {
		failedAccounts <- index
		return
	} else {
		mutex.Lock()
		KYCLinks.Store(privateKey, KYCLink)
		mutex.Unlock()
	}
}

func MobileProxyWrapper(data MobileProxyData) {
	for i := range data.Indexes {
		Process(i, data.Proxy, data.PrivateKeys[i], data.Config, data.KYCLinks, data.Mutex, data.FailedAccounts)
		extra.RandomSleep(data.Config.Info.PauseBetweenAccounts.Start, data.Config.Info.PauseBetweenAccounts.End)

		extra.ChangeProxyURL(data.IPChangeLink)
		extra.Logger{}.Success("Successfully changed the IP address of mobile proxies")
		time.Sleep(time.Duration(data.Config.Proxy.ChangeIPPause) * time.Second)

	}
}

type MobileProxyData struct {
	Proxy          string
	IPChangeLink   string
	Indexes        chan int
	PrivateKeys    []string
	Config         extra.Config
	KYCLinks       *sync.Map
	Mutex          *sync.Mutex
	FailedAccounts chan int
}

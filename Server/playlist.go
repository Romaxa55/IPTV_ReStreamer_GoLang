package Server

import (
	"IPTV_ReStreamer_GoLang/Config"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
)

var config = Config.GetServerConfig()
var playlist *PlayList

type PlayList struct {
	Content  string
	Original string
}

func NewPlayList(content string, original string) *PlayList {
	return &PlayList{
		Content:  content,
		Original: original,
	}
}

func init() {
	err := LogApp.InitLogger("Parser")
	if err != nil {
		fmt.Print(err)
		return
	}

	resp, err := http.Get(config.IptvUrl)
	if err != nil {
		log.Error("Failed to download M3U8 playlist:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Error("Failed to download M3U8 playlist. Status code:", resp.StatusCode)
		return
	}
	// Читаем содержимое плейлиста
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("Failed to read M3U8 playlist body:", err)
		return
	}
	tmp := string(body)
	pls := replaceDomainWithIP(tmp)
	playlist = NewPlayList(pls, tmp)
}

func replaceDomainWithIP(playlist string) string {
	// Define a regular expression to match the domain name in the playlist
	domainRegex := regexp.MustCompile(`http://[^/]+`)
	if config.IP == "0.0.0.0" {
		externalIP, err := getExternalIP()
		if err != nil {
			log.Error("Error getting external IP:", err)
		}
		config.IP = externalIP
	}
	// Replace the matched domain name with the external IP address
	replacedPlaylist := domainRegex.ReplaceAllString(playlist, "http://"+config.IP+":"+config.Port)

	return replacedPlaylist
}

func getExternalIP() (string, error) {
	resp, err := http.Get("https://api64.ipify.org?format=json")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Extract the IP address from the response
	ipAddress := strings.Trim(string(body), `{"ip":"}"`)
	return ipAddress, nil
}

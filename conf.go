package main

import (
	"encoding/json"
	"os"
)

var config Conf

type MatchRule struct {
	Domain []string
	Ip     []string
}

type Conf struct {
	RoutePath     string               `json:"RoutePath"`
	OutTag        string               `json:"OutTag"`
	CheckGlobal   bool                 `json:"CheckGlobal"`
	CheckRegion   bool                 `json:"CheckRegion"`
	MediaList     map[string][]string  `json:"MediaList"`
	MatchRuleList map[string]MatchRule `json:"MatchRuleList"`
}

func LoadConfig(path string) error {
	// gen default config
	config = Conf{
		RoutePath:   "/etc/V2bX/route.json",
		OutTag:      "unlock",
		CheckGlobal: true,
		CheckRegion: false,
		MediaList: map[string][]string{
			"global": {
				"Hotstar",
				"Disney+",
				"Netflix",
				"Netflix CDN",
				"Youtube",
				"Youtube CDN",
				"Amazon Prime Video",
				"TVBAnywhere+",
				"iQyi",
				"Viu.com",
				"Spotify",
				"Steam",
				"ChatGPT",
			},
			"jp": {
				"DMM",
				"Abema",
				"Niconico",
				"music.jp",
				"Telasa",
				"Paravi",
				"U-NEXT",
				"Hulu Japan",
				"GYAO!",
				"VideoMarket",
				"FOD(Fuji TV)",
				"Radiko",
				"Karaoke@DAM",
				"J:COM On Demand",
				"Kancolle",
				"Pretty Derby Japan",
				"Konosuba Fantastic Days",
				"Princess Connect Re:Dive Japan",
				"World Flipper Japan",
				"Project Sekai: Colorful Stage",
			},
			"hk": {
				"Now E",
				"Viu.TV",
				"MyTVSuper",
				"HBO GO Aisa",
				"BiliBili Hongkong/Macau Only",
			},
			"us": {
				"FOX",
				"Hulu",
				"ESPN+",
				"Epix",
				"Starz",
				"Philo",
				"FXNOW",
				"TLC GO",
				"HBO Max",
				"Shudder",
				"BritBox",
				"CW TV",
				"NBA TV",
				"Tubi TV",
				"Sling TV",
				"Pluto TV",
				"Acorn TV",
				"SHOWTIME",
				"encoreTVB",
				"Funimation",
				"Discovery+",
				"Paramount+",
				"Peacock TV",
				"Popcornflix",
				"Crunchyroll",
				"Direct Stream",
				"CBC Gem",
			},
			"tw": {
				"KKTV",
				"LiTV",
				"MyVideo",
				"4GTV",
				"LineTV",
				"Hami Video",
				"CatchPlay+",
				"Bahamut Anime",
				"Bilibili Taiwan Only",
			},
		},
		MatchRuleList: map[string]MatchRule{
			"Netflix": {
				Domain: []string{
					"geosite:netflix",
				},
				Ip: []string{
					"geoip:netflix",
				},
			},
			"Niconico": {
				Domain: []string{
					"geosite:niconico",
				},
			},
			"DMM": {
				Domain: []string{
					"geosite:dmm",
				},
			},
			"Abema": {
				Domain: []string{
					"geosite:abema",
				},
			},
			"music.jp": {
				Domain: []string{
					"domain:music-book.jp",
				},
			},
			"Telasa": {
				Domain: []string{
					"domain:telasa.jp",
					"domain:kddi-video.com",
					"domain:videopass.jp",
					"domain:d2lmsumy47c8as.cloudfront.net",
				},
			},
			"Paravi": {
				Domain: []string{
					"domain:paravi.jp",
				},
			},
			"U-NEXT": {
				Domain: []string{
					"domain:unext.jp",
					"domain:nxtv.jp",
				}},
			"Hulu Japan": {
				Domain: []string{
					"geosite:hulu",
				},
			},
		},
	}
	// load config file
	f, err := os.Open(path)
	if err != nil {
		return err
	}
	defer f.Close()
	err = json.NewDecoder(f).Decode(&config)
	if err != nil {
		return err
	}
	return nil
}

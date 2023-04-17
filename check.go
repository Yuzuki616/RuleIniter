package main

import (
	"errors"
	"fmt"
	m "github.com/Yuzuki616/MediaUnlockTest"
	"github.com/hashicorp/go-multierror"
	"net/http"
	"sync"
)

var c = m.AutoHttpClient
var MediaList = map[string]func(http.Client) m.Result{
	// Multination
	"Hotstar":            m.Hotstar,
	"Disney+":            m.DisneyPlus,
	"Netflix":            m.NetflixRegion,
	"Netflix CDN":        m.NetflixCDN,
	"Youtube":            m.YoutubeRegion,
	"Youtube CDN":        m.YoutubeCDN,
	"Amazon Prime Video": m.PrimeVideo,
	"TVBAnywhere+":       m.TVBAnywhere,
	"iQyi":               m.IqRegion,
	"Viu.com":            m.ViuCom,
	"Spotify":            m.Spotify,
	"Steam":              m.Steam,
	"ChatGPT":            m.ChatGPT,
	// HongKong
	"Now E":                        m.NowE,
	"Viu.TV":                       m.ViuTV,
	"MyTVSuper":                    m.MyTvSuper,
	"HBO GO Aisa":                  m.HboGoAisa,
	"BiliBili Hongkong/Macau Only": m.BilibiliHKMC,
	// TaiWan
	"KKTV":                 m.KKTV,
	"LiTV":                 m.LiTV,
	"MyVideo":              m.MyVideo,
	"4GTV":                 m.TW4GTV,
	"LineTV":               m.LineTV,
	"Hami Video":           m.HamiVideo,
	"CatchPlay+":           m.Catchplay,
	"Bahamut Anime":        m.BahamutAnime,
	"Bilibili Taiwan Only": m.BilibiliTW,

	"DMM":                            m.DMM,
	"Abema":                          m.Abema,
	"Niconico":                       m.Niconico,
	"music.jp":                       m.MusicJP,
	"Telasa":                         m.Telasa,
	"Paravi":                         m.Paravi,
	"U-NEXT":                         m.U_NEXT,
	"Hulu Japan":                     m.HuluJP,
	"GYAO!":                          m.GYAO,
	"VideoMarket":                    m.VideoMarket,
	"FOD(Fuji TV)":                   m.FOD,
	"Radiko":                         m.Radiko,
	"Karaoke@DAM":                    m.Karaoke,
	"J:COM On Demand":                m.J_COM_ON_DEMAND,
	"Kancolle":                       m.Kancolle,
	"Pretty Derby Japan":             m.PrettyDerbyJP,
	"Konosuba Fantastic Days":        m.KonosubaFD,
	"Princess Connect Re:Dive Japan": m.PCRJP,
	"World Flipper Japan":            m.WFJP,
	"Project Sekai: Colorful Stage":  m.PJSK,
	"FOX":                            m.Fox,
	"Hulu":                           m.Hulu,
	"ESPN+":                          m.ESPNPlus,
	"Epix":                           m.Epix,
	"Starz":                          m.Starz,
	"Philo":                          m.Philo,
	"FXNOW":                          m.FXNOW,
	"TLC GO":                         m.TlcGo,
	"HBO Max":                        m.HBOMax,
	"Shudder":                        m.Shudder,
	"BritBox":                        m.BritBox,
	"CW TV":                          m.CW_TV,
	"NBA TV":                         m.NBA_TV,
	"Tubi TV":                        m.TubiTV,
	"Sling TV":                       m.SlingTV,
	"Pluto TV":                       m.PlutoTV,
	"Acorn TV":                       m.AcornTV,
	"SHOWTIME":                       m.SHOWTIME,
	"encoreTVB":                      m.EncoreTVB,
	"Funimation":                     m.Funimation,
	"Discovery+":                     m.DiscoveryPlus,
	"Paramount+":                     m.ParamountPlus,
	"Peacock TV":                     m.PeacockTV,
	"Popcornflix":                    m.Popcornflix,
	"Crunchyroll":                    m.Crunchyroll,
	"Direct Stream":                  m.DirectvStream,
	"CBC Gem":                        m.CBCGem,
}

type Result struct {
	Name string
	Err  error
}

func CheckMediaUnlock(region string) ([]string, error) {
	var medias []string
	if region == "all" {
		// add all
		for _, n := range config.MediaList {
			medias = append(medias, n...)
		}
		config.CheckRegion = false
	} else {
		// add region
		if ms, e := config.MediaList[region]; e {
			medias = ms
		} else if !config.CheckRegion {
			return nil, errors.New("not have the region")
		}
		// add global
		if region == "global" {
			config.CheckRegion = false
		} else {
			if config.CheckGlobal {
				medias = append(medias, config.MediaList["global"]...)
			}
		}
	}
	rs := make(chan *Result)
	done := make(chan struct{})
	var fails []string
	var errs error
	temp := make(map[string]struct{}, len(medias))
	go func() {
		wg := sync.WaitGroup{}
		for _, n := range medias {
			if _, e := temp[n]; e {
				continue
			}
			temp[n] = struct{}{}
			n := n
			wg.Add(1)
			go func() {
				defer wg.Done()
				r := MediaList[n](c)
				if r.Region != region && r.Region != "" {
					rs <- &Result{
						Name: n,
						Err:  r.Err,
					}
					return
				}
				if !r.Success {
					rs <- &Result{
						Name: n,
						Err:  r.Err,
					}
					return
				}
				rs <- nil
			}()
		}
		wg.Wait()
		close(done)
	}()
	for {
		select {
		case <-done:
			return fails, errs
		case n := <-rs:
			if n != nil {
				if n.Err != nil {
					errs = multierror.Append(errs, fmt.Errorf("%s error: %s", n.Name, n.Err))
				}
			}
		}
	}
}

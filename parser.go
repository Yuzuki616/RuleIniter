package main

import (
	"encoding/json"
	"fmt"
	"github.com/xtls/xray-core/infra/conf"
	"os"
)

type Rule struct {
	Type        string   `json:"type"`
	OutboundTag string   `json:"outboundTag"`
	Domain      []string `json:"domain"`
}

func RuleChanged(old, new []string) bool {
	tmp := map[string]struct{}{}
	tmp2 := map[string]struct{}{}
	for i := range old {
		tmp[old[i]] = struct{}{}
	}
	l := len(tmp)
	for i := range new {
		e := new[i]
		tmp[e] = struct{}{}
		tmp2[e] = struct{}{}
		if l != len(tmp) {
			return true
		}
	}
	tmp = nil
	l = len(tmp2)
	for i := range old {
		tmp2[old[i]] = struct{}{}
		if l != len(tmp2) {
			return true
		}
	}
	return false
}

func ParseRouteConf(medias []string) error {
	route := &conf.RouterConfig{}
	f, err := os.OpenFile(config.RoutePath, os.O_RDWR, 0744)
	if err != nil {
		return fmt.Errorf("open route file error: %s", err)
	}
	err = json.NewDecoder(f).Decode(route)
	if err != nil {
		return fmt.Errorf("decode route file error: %s", err)
	}
	var domains []string
	for _, m := range medias {
		domains = append(domains, config.MatchRuleList[m]...)
	}
	save := false
	for i := range route.RuleList {
		rule := Rule{}
		err := json.Unmarshal(route.RuleList[i], &rule)
		if err != nil {
			return fmt.Errorf("parse rule error: %s", err)
		}
		if rule.OutboundTag == config.OutTag {
			if !RuleChanged(rule.Domain, domains) {
				return nil
			}
			rule.Domain = make([]string, 0)
			rule.Domain = domains
			r, _ := json.Marshal(rule)
			route.RuleList[i] = r
			save = true
		}
	}
	if !save {
		rule := Rule{
			Type:        "field",
			OutboundTag: config.OutTag,
			Domain:      domains,
		}
		r, _ := json.Marshal(rule)
		route.RuleList = append(route.RuleList, r)
	}
	err = f.Truncate(0)
	if err != nil {
		return fmt.Errorf("truncate file error: %s", err)
	}
	f.Seek(0, 0)
	err = json.NewEncoder(f).Encode(route)
	if err != nil {
		return fmt.Errorf("encode route file error: %s", err)
	}
	return nil
}

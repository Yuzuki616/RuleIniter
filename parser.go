package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Routing struct {
	DomainStrategy string          `json:"domainStrategy,omitempty"`
	DomainMatcher  string          `json:"domainMatcher,omitempty"`
	Rules          []Rule          `json:"rules,omitempty"`
	Balancers      json.RawMessage `json:"balancers,omitempty"`
}
type Rule struct {
	DomainMatcher string   `json:"domainMatcher,omitempty"`
	Type          string   `json:"type,omitempty"`
	Domain        []string `json:"domain,omitempty"`
	Ip            []string `json:"ip,omitempty"`
	Port          string   `json:"port,omitempty"`
	SourcePort    string   `json:"sourcePort,omitempty"`
	Network       string   `json:"network,omitempty"`
	Source        []string `json:"source,omitempty"`
	User          []string `json:"user,omitempty"`
	InboundTag    []string `json:"inboundTag,omitempty"`
	Protocol      []string `json:"protocol,omitempty"`
	Attrs         string   `json:"attrs,omitempty"`
	OutboundTag   string   `json:"outboundTag,omitempty"`
	BalancerTag   string   `json:"balancerTag,omitempty"`
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
	route := &Routing{}
	f, err := os.OpenFile(config.RoutePath, os.O_RDWR, 0744)
	if err != nil {
		return fmt.Errorf("open route file error: %s", err)
	}
	err = json.NewDecoder(f).Decode(route)
	if err != nil {
		return fmt.Errorf("decode route file error: %s", err)
	}
	var domains []string
	var ips []string
	for _, m := range medias {
		domains = append(domains, config.MatchRuleList[m].Domain...)
		ips = append(ips, config.MatchRuleList[m].Ip...)
	}
	save := false
	for i := range route.Rules {
		if route.Rules[i].OutboundTag != config.OutTag {
			continue
		}
		changed := true
		if RuleChanged(route.Rules[i].Domain, domains) {
			route.Rules[i].Domain = domains
			save = true
		} else {
			changed = false
		}
		if RuleChanged(route.Rules[i].Ip, ips) {
			route.Rules[i].Ip = ips
			save = true
		} else if !changed {
			return nil
		}
		if changed {
			break
		}
	}
	if !save {
		rule := Rule{
			Type:        "field",
			OutboundTag: config.OutTag,
			Domain:      domains,
			Ip:          ips,
		}
		route.Rules = append(route.Rules, rule)
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

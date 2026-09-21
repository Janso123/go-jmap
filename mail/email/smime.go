package email

import (
	"strings"

	"github.com/Janso123/go-jmap"
	"github.com/Janso123/go-jmap/mail"
)

func propertiesNeedSMIME(props []string) bool {
	for _, p := range props {
		if strings.HasPrefix(p, "smime") {
			return true
		}
	}
	return false
}

func filterNeedsSMIME(f jmap.Filter) bool {
	if f == nil {
		return false
	}
	switch x := f.(type) {
	case *FilterCondition:
		return x.HasSMIME != nil || x.HasVerifiedSMIME != nil || x.HasVerifiedSMIMEAtDelivery != nil
	case FilterCondition:
		return x.HasSMIME != nil || x.HasVerifiedSMIME != nil || x.HasVerifiedSMIMEAtDelivery != nil
	case *jmap.FilterOperator:
		for _, c := range x.Conditions {
			if filterNeedsSMIME(c) {
				return true
			}
		}
	case jmap.FilterOperator:
		for _, c := range x.Conditions {
			if filterNeedsSMIME(c) {
				return true
			}
		}
	}
	return false
}

func withSMIME(uris []jmap.URI, need bool) []jmap.URI {
	if !need {
		return uris
	}
	return append(uris, SMIMEVerify)
}

func mailRequires(needSMIME bool) []jmap.URI {
	return withSMIME([]jmap.URI{mail.URI}, needSMIME)
}

const SMIMEVerify jmap.URI = "urn:ietf:params:jmap:smimeverify"

type smimeVerify struct{}

func (s *smimeVerify) URI() jmap.URI { return SMIMEVerify }

func (s *smimeVerify) New() jmap.Capability { return &smimeVerify{} }

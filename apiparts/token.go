// This file is part of Epixel CRM Software

package epicrm_apiparts

import (
	"regexp"
)

func extractBearerTokenFromHeader(authHdr string) []byte {
	// TODO keep globally for reuse
	btokre := regexp.MustCompile("^Bearer (.+)$")
	
	matchAndBtoken := btokre.FindSubmatch([]byte(authHdr))
	if len(matchAndBtoken) == 2 {
		return matchAndBtoken[1]
	}
	
	return nil
}

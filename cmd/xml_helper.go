package cmd

import "strings"

type loginRes struct {
	sID    string
	orgID  string
	orgURL string
}

func parseLogin(data string) loginRes {
	var lr loginRes
	lr.sID = GetStringInBetween(data, "<sessionId>", "</sessionId>")
	lr.orgID = GetStringInBetween(data, "<organizationId>", "</organizationId>")
	lr.orgURL = GetStringInBetween(GetStringInBetween(data, "<serverUrl>", "</serverUrl>"), "", "/services")
	return lr
}

// GetStringInBetween ...
func GetStringInBetween(str string, start string, end string) (result string) {
	s := strings.Index(str, start)
	if s == -1 {
		return
	}
	s += len(start)
	e := strings.Index(str, end)
	return str[s:e]
}

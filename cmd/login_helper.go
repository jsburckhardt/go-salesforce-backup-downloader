package cmd

import (
	"fmt"
	"strings"

	"github.com/go-playground/log"
)

const soap = `<se:Envelope xmlns:se="http://schemas.xmlsoap.org/soap/envelope/">
				<se:Header/>
				<se:Body>
					<login xmlns="urn:partner.soap.sforce.com">
						<username>%s</username>
						<password>%s</password>
					</login>
				</se:Body>
			</se:Envelope>`

func login() loginRes {

	//return parseLogin(xml)

	url := "https://login.salesforce.com/services/Soap/u/40.0"

	headers := map[string]string{
		"Content-Type": "text/xml",
		"SOAPAction":   "Required",
	}

	payload := fmt.Sprintf(soap, salesForceUserName, salesForceUserPassword)
	data := []byte(payload)

	status, res, err, _ := sendRequest("POST", url, headers, []byte(data))
	if err != nil || status >= 400 {
		log.Fatalf("Error upon login: %v", err)
	} else if strings.Contains(GetStringInBetween(string(res), "<faultcode>", "</faultcode>"), "INVALID_LOGIN") {
		log.Fatal("INVALID_LOGIN")
	}

	// fmt.Println("*** response ***")
	// fmt.Println(status)
	// fmt.Println(string(res))
	// fmt.Println(parseLogin(string(res)))

	return parseLogin(string(res))
}

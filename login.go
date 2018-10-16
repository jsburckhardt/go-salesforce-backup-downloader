package main

import (
	"fmt"

	"github.com/go-playground/log"
	"github.com/spf13/viper"
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

	payload := fmt.Sprintf(soap, viper.GetString("sf.username"), viper.GetString("sf.password"))
	data := []byte(payload)

	status, res, err, _ := sendRequest("POST", url, headers, []byte(data))
	if err != nil || status >= 400 {
		log.Fatalf("Error upon login: %v", err)
	}

	// fmt.Println("*** response ***")
	// fmt.Println(string(res))
	// fmt.Println(parseLogin(string(res)))

	return parseLogin(string(res))
}

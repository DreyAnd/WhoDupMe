package CSRF

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"time"

	"github.com/DreyAnd/WhoDupMe/pkg/args"
	"github.com/DreyAnd/WhoDupMe/pkg/httpclient"
)

func Get_Token(opts args.Options) string {
	client := httpclient.NewHttpClient(10 * time.Second)

	client.SetCookies([]*http.Cookie{
		{Name: "__Host-session", Value: opts.H1_Session},
	})

	resp, err := client.Get(fmt.Sprintf("https://hackerone.com/%s/hacktivity", opts.Program_Name))
	if err != nil {
		fmt.Errorf("[-] Error: %s\n", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Errorf("[-] Error: %s\n", err)
		os.Exit(1)
	}

	match := regexp.MustCompile(`<meta name="csrf-token" content="([^"]+)" />`).FindStringSubmatch(string(body))
	if len(match) < 2 {
		fmt.Errorf("[-] CSRF not found.")
		os.Exit(1)
	}

	CSRF_token := match[1]

	return CSRF_token
}

package hacktivityexplorer

import (
	"bytes"
	b64 "encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/DreyAnd/WhoDupMe/pkg/args"
	"github.com/DreyAnd/WhoDupMe/pkg/httpclient"
)

type GraphQLSchema struct {
	Data struct {
		Me struct {
			ID       string `json:"id"`
			Typename string `json:"__typename"`
		} `json:"me"`
		HacktivityItems struct {
			TotalCount int `json:"total_count"`
			PageInfo   struct {
				EndCursor   string `json:"endCursor"`
				HasNextPage bool   `json:"hasNextPage"`
				Typename    string `json:"__typename"`
			} `json:"pageInfo"`
			Edges []struct {
				Node struct {
					ID         string `json:"id"`
					DatabaseID string `json:"databaseId"`
					Typename   string `json:"__typename"`
					Type       string `json:"type"`
					Votes      struct {
						TotalCount int    `json:"total_count"`
						Typename   string `json:"__typename"`
					} `json:"votes"`
					Upvoted  bool `json:"upvoted"`
					Reporter struct {
						ID       string `json:"id"`
						Username string `json:"username"`
						Typename string `json:"__typename"`
					} `json:"reporter"`
					Team struct {
						Handle               string `json:"handle"`
						Name                 string `json:"name"`
						MediumProfilePicture string `json:"medium_profile_picture"`
						URL                  string `json:"url"`
						ID                   string `json:"id"`
						Typename             string `json:"__typename"`
					} `json:"team"`
					LatestDisclosableAction     string    `json:"latest_disclosable_action"`
					LatestDisclosableActivityAt time.Time `json:"latest_disclosable_activity_at"`
					RequiresViewPrivilege       bool      `json:"requires_view_privilege"`
					TotalAwardedAmount          any       `json:"total_awarded_amount"`
					Currency                    string    `json:"currency"`
				} `json:"node"`
				Typename string `json:"__typename"`
			} `json:"edges"`
			Typename string `json:"__typename"`
		} `json:"hacktivity_items"`
	} `json:"data"`
}

type Resolved_Report_Info struct {
	DatabaseID       string
	ReporterUsername string
}

var resolved_rep_info []Resolved_Report_Info
var resolvedRepInfoMutex sync.Mutex

func Load_Reports(wg *sync.WaitGroup, client *httpclient.HttpClient, data string, resultChan chan<- bool) {
	defer wg.Done()

	var end_cursor bool

	resp, err := client.Post("https://hackerone.com/graphql", bytes.NewBufferString(data))
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

	var response GraphQLSchema
	if err := json.Unmarshal([]byte(body), &response); err != nil {
		panic(err)
	}

	if string(response.Data.HacktivityItems.PageInfo.EndCursor) == "" {
		end_cursor = true
	} else {
		end_cursor = false
		for _, edge := range response.Data.HacktivityItems.Edges {
			edgeData := Resolved_Report_Info{
				DatabaseID:       edge.Node.DatabaseID,
				ReporterUsername: edge.Node.Reporter.Username,
			}

			resolvedRepInfoMutex.Lock() // prevent race conditions
			resolved_rep_info = append(resolved_rep_info, edgeData)
			resolvedRepInfoMutex.Unlock()

		}

	}

	resultChan <- end_cursor

}

func Get_All_Report_Info(opts args.Options, CSRF_token string) []Resolved_Report_Info {
	client := httpclient.NewHttpClient(5 * time.Second)

	client.SetCookies([]*http.Cookie{
		{Name: "__Host-session", Value: opts.H1_Session},
	})
	client.SetHeader("Content-Type", "application/json")
	client.SetHeader("X-Csrf-Token", CSRF_token)

	var wg sync.WaitGroup
	resultChan := make(chan bool)

	count := 0
	for {
		wg.Add(1)

		cursor := strings.ReplaceAll(b64.StdEncoding.EncodeToString([]byte(fmt.Sprintln(count))), "=", "")
		data := fmt.Sprintf(`{"operationName":"TeamHacktivityPageQuery","variables":{"where":{"team":{"handle":{"_eq":"%s"}}},"orderBy":{"field":"popular","direction":"DESC"},"secureOrderBy":null,"count":25, "cursor":"%s"},"query":"query TeamHacktivityPageQuery($orderBy: HacktivityItemOrderInput, $secureOrderBy: FiltersHacktivityItemFilterOrder, $where: FiltersHacktivityItemFilterInput, $count: Int, $cursor: String) {\n  me {\n    id\n    __typename\n  }\n  hacktivity_items(\n    first: $count\n    after: $cursor\n    order_by: $orderBy\n    secure_order_by: $secureOrderBy\n    where: $where\n  ) {\n    total_count\n    ...HacktivityList\n    __typename\n  }\n}\n\nfragment HacktivityList on HacktivityItemConnection {\n  pageInfo {\n    endCursor\n    hasNextPage\n    __typename\n  }\n  edges {\n    node {\n      ... on HacktivityItemInterface {\n        id\n        databaseId: _id\n        __typename\n      }\n      __typename\n    }\n    ...HacktivityItem\n    __typename\n  }\n  __typename\n}\n\nfragment HacktivityItem on HacktivityItemUnionEdge {\n  node {\n    ... on HacktivityItemInterface {\n      id\n      type: __typename\n    }\n    ... on Undisclosed {\n      id\n      ...HacktivityItemUndisclosed\n      __typename\n    }\n    ... on Disclosed {\n      id\n      ...HacktivityItemDisclosed\n      __typename\n    }\n    ... on HackerPublished {\n      id\n      ...HacktivityItemHackerPublished\n      __typename\n    }\n    __typename\n  }\n  __typename\n}\n\nfragment HacktivityItemUndisclosed on Undisclosed {\n  id\n  votes {\n    total_count\n    __typename\n  }\n  upvoted: upvoted_by_current_user\n  reporter {\n    id\n    username\n    ...UserLinkWithMiniProfile\n    __typename\n  }\n  team {\n    handle\n    name\n    medium_profile_picture: profile_picture(size: medium)\n    url\n    id\n    ...TeamLinkWithMiniProfile\n    __typename\n  }\n  latest_disclosable_action\n  latest_disclosable_activity_at\n  requires_view_privilege\n  total_awarded_amount\n  currency\n  __typename\n}\n\nfragment TeamLinkWithMiniProfile on Team {\n  id\n  handle\n  name\n  __typename\n}\n\nfragment UserLinkWithMiniProfile on User {\n  id\n  username\n  __typename\n}\n\nfragment HacktivityItemDisclosed on Disclosed {\n  id\n  reporter {\n    id\n    username\n    ...UserLinkWithMiniProfile\n    __typename\n  }\n  votes {\n    total_count\n    __typename\n  }\n  upvoted: upvoted_by_current_user\n  team {\n    handle\n    name\n    medium_profile_picture: profile_picture(size: medium)\n    url\n    id\n    ...TeamLinkWithMiniProfile\n    __typename\n  }\n  report {\n    id\n    databaseId: _id\n    title\n    substate\n    url\n    report_generated_content {\n      id\n      hacktivity_summary\n      __typename\n    }\n    __typename\n  }\n  latest_disclosable_action\n  latest_disclosable_activity_at\n  total_awarded_amount\n  severity_rating\n  currency\n  __typename\n}\n\nfragment HacktivityItemHackerPublished on HackerPublished {\n  id\n  reporter {\n    id\n    username\n    ...UserLinkWithMiniProfile\n    __typename\n  }\n  votes {\n    total_count\n    __typename\n  }\n  upvoted: upvoted_by_current_user\n  team {\n    id\n    handle\n    name\n    medium_profile_picture: profile_picture(size: medium)\n    url\n    ...TeamLinkWithMiniProfile\n    __typename\n  }\n  report {\n    id\n    url\n    title\n    substate\n    __typename\n  }\n  latest_disclosable_activity_at\n  severity_rating\n  __typename\n}\n"}`, opts.Program_Name, cursor)

		go Load_Reports(&wg, client, data, resultChan)
		end_cursor := <-resultChan

		if end_cursor == true {
			break
		}

		count++

	}
	wg.Wait()

	return resolved_rep_info
}

func Find_The_Duper(opts args.Options, resolved_report_info []Resolved_Report_Info) string {
	var duper string
	for _, info := range resolved_report_info {
		if info.DatabaseID == opts.Report_ID {
			duper = info.ReporterUsername
			break
		}
	}

	return duper
}

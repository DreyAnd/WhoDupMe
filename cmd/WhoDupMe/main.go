package main

import (
	"fmt"

	"github.com/DreyAnd/WhoDupMe/pkg/CSRF"

	hacktivityexplorer "github.com/DreyAnd/WhoDupMe/pkg/HacktivityExplorer"

	"github.com/DreyAnd/WhoDupMe/pkg/args"
)

func banner() {
	banner := `
		╭╮╭╮╭┳╮╱╱╱╱╭━━━╮╱╱╱╱╱╭━╮╭━╮
		┃┃┃┃┃┃┃╱╱╱╱╰╮╭╮┃╱╱╱╱╱┃┃╰╯┃┃
		┃┃┃┃┃┃╰━┳━━╮┃┃┃┣╮╭┳━━┫╭╮╭╮┣━━╮
		┃╰╯╰╯┃╭╮┃╭╮┃┃┃┃┃┃┃┃╭╮┃┃┃┃┃┃┃━┫
		╰╮╭╮╭┫┃┃┃╰╯┣╯╰╯┃╰╯┃╰╯┃┃┃┃┃┃┃━┫
		╱╰╯╰╯╰╯╰┻━━┻━━━┻━━┫╭━┻╯╰╯╰┻━━╯
		╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱┃┃
		╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╰╯
	`

	fmt.Println(banner)
}

func main() {
	banner()                                                                          // Print the tool banner
	opts := args.GetArgs()                                                            // Handle Command-line arguments
	CSRF_token := CSRF.Get_Token(opts)                                                // Get the Hackerone CSRF Token
	resolved_reports_info := hacktivityexplorer.Get_All_Report_Info(opts, CSRF_token) // Get Information about Resolved Reports
	duper := hacktivityexplorer.Find_The_Duper(opts, resolved_reports_info)           // Find the person who duped you

	if duper != "" {
		fmt.Println("[\033[1;32m+\033[0m] The person that gave you a duplicate found.")
		fmt.Printf("[\033[1;32m+\033[0m] The following person gave you a duplicate: \033[1;36m https://hackerone.com/%s \033[1;0m", duper)
	} else {
		fmt.Println("[\033[1;31m-\033[0m] Unfortunately, we were not able to find the person that gave you a duplicate.")
		fmt.Println("[\033[1;31m-\033[0m] Make sure you supplied the correct parameters/arguments.")
	}

}

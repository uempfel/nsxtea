/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"os"
	"net/http"
	"io/ioutil"
	"strconv"
    "strings"
)

var includedFields string
var sortBy string
var sortAscending bool
var cursor string
var pageSize string
var useManagerApi bool



// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search <query>",
	Short: "Interact with the Policy or Manager Search API",
	Long: `Interact with the Policy or Manager Search API

QUERY SYNTAX
A query is broken up into terms and operators. 
A term is case insensitive and can be a single word such as "Hello" or " World " or 
a phrase surrounded by double quotes such as "Hello World", which would search for the exact phrase.

FIELD NAMES
By default, all the fields will be searched for the search term specified. 
Specific fields to be searched on can be provided via the field name followed by a colon ":" and then the search term.
- To search for all the entities where display_name field is "App-VM-1", use
  display_name:App-VM-1
- Use the dot notation to search on nested fields
  tags.scope:prod

WILDCARDS
Wildcard searches can be run using ? to substitute a single character or * to substitute zero or more characters
  *vm* will match all the entities that contains the term "vm" in any of its fields
  display_name:App-VM-? will match App-VM-1, App-VM-2, App-VM-A etc..
  display_name:App* will match everything where display_name begins with App
Warning: Be aware that using wildcards especially at the beginning of a word i.e. *vm can use a large amount of memory and may perform badly.
	
BOOLEAN OPERATORS
Search terms can be combined using boolean operators AND, OR and NOT. (Note: Boolean operators must be ALL CAPS).
  AND
  The AND ( && ) operator matches entities where both terms exist in any of the fields of an entity.
  To search for Firewall rule with display_name containing block, use
    display_name:*block* AND resource_type:FirewallRule

  OR
  The OR ( || ) operator links two terms and finds matching entities if either of the terms exists in an entity.
  To search for Firewall rule with display_name containing either block or allow, use
    display_name:*block* OR display_name:*allow* AND resource_type:FirewallRule
    display_name:(*block* OR *allow*) AND resource_type:FirewallRule

  NOT
  The NOT ( ! ) operator excludes entities that contain the term after NOT.
  To search for Firewall rule with display_name does not contain the term block
    NOT display_name:*block* AND resource_type:FirewallRule
    !display_name:*block* AND resource_type:FirewallRule

RANGES
Ranges can be specified for numeric or string fields and use the following syntax
  vni:>50001
  vni:>=50001
  vni:<90000
  vni:<=90000    
To combine an upper and lower bound, you would need to join two clauses with AND operator:
  vni:(>=50001 AND <90000)

RESERVED CHARACTERS
If characters which function as operators are to be used in the query (not as operators), then they should be escaped with a leading backslash.
To search for (a+b)=c
  \(a\+b\)\=c.
The reserved characters are: + - = && || > < ! ( ) { } [ ] ^ " ~ * ? : \ /
Failing to escape these reserved characters correctly would lead to syntax errors and prevent the query from running.	
`,
    Run: handleSearch,
    Args: cobra.MinimumNArgs(1),
}

func handleSearch(cmd *cobra.Command, args []string) {

    endpoint := "/api/v1/search/query?"
    if !useManagerApi {
      endpoint = "/policy" + endpoint
    }
    req, err := http.NewRequest("GET", "https://"+ os.Getenv("NSXTEA_URL")+ endpoint, nil)

    query := concatArgs(args)
    q := req.URL.Query()
    q.Add("query", query)
    q.Add("included_fields", includedFields)
    q.Add("sort_ascending", strconv.FormatBool(sortAscending))
    q.Add("sort_by", sortBy)
    q.Add("cursor", cursor)
    q.Add("page_size", pageSize)
    req.URL.RawQuery = q.Encode()

    req.SetBasicAuth(
        os.Getenv("NSXTEA_USERNAME"),
        os.Getenv("NSXTEA_PASSWORD"),
    )    
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil{
        fmt.Println(err)
        os.Exit(1)
    }

    bodyText, err := ioutil.ReadAll(resp.Body)
    printErrIfNotNil(err)
    s := string(bodyText)
    fmt.Println(s)
}

func init() {

	rootCmd.AddCommand(searchCmd)
	searchCmd.PersistentFlags().StringVarP(&includedFields, "included_fields", "f", "", "Comma separated list of fields that should be included in query result")
	searchCmd.PersistentFlags().StringVarP(&sortBy, "sort_by", "s", "", "Field by which records are sorted")
  searchCmd.PersistentFlags().BoolVarP(&sortAscending, "sort_ascending","a", true, "Sorting order of the query results")
  searchCmd.PersistentFlags().StringVarP(&cursor, "cursor", "c", "", "Opaque cursor to be used for getting next page of records (supplied by current result page)")
  searchCmd.PersistentFlags().StringVarP(&pageSize, "page_size", "p", "1000", "Maximum number of results to return in this page \nMin: 0, Max: 1000")
  searchCmd.PersistentFlags().BoolVarP(&useManagerApi, "manager", "m", false, "Use the Manager API for the search request")
}

func concatArgs(args[] string) string {
	var concatArgs string
	for i := 0; i < len(args); i++ {
		concatArgs += args[i] + " "
	}
	return strings.TrimSpace(concatArgs)
}

func printErrIfNotNil(err error) {
	if err != nil {
		fmt.Println(err)
	}
}


![nsxtea logo](https://user-images.githubusercontent.com/28337775/106465522-593f6480-649a-11eb-9262-bf94e3add74a.png)

Nsxtea is a small CLI to interact with the VMware NSX-T Policy Search API. Sip a cup of tea while others search in the GUI!

## Installation
Assuming you have already [installed go](https://golang.org/doc/install):

```sh
export GO111MODULE=on
go get github.com/uempfel/nsxtea
go install github.com/uempfel/nsxtea
```

If you're getting  a `command not found` error, you might have to add your "go binary" directory to your PATH. To do so, run:
```sh 
export PATH=$PATH:$(go env GOPATH)/bin
```

## Configuration
The CLI is configured via environment variables. Just replace the following sample values with your own, and you're good to go:  

```bash
export NSXTEA_URL='nsxteamgr.mydomain.com'
export NSXTEA_USERNAME='admin'
export NSXTEA_PASSWORD='VMware1!VMware1!'
export NSXTEA_INSECURE='true'
```

| Name | Description | Example| Necessity |
|------|-------------|-----------|---------------|
| NSXTEA_URL | FQDN of your NSX-Manager (without the protocol part) | `nsxteamgr.mydomain.com` | required |
| NSXTEA_USERNAME | The username to authenticate with the NSX-Manager. The user must have at least read permissions for the API | `admin` | required |
| NSXTEA_PASSWORD | The password for the given user to authenticate with the NSX-Manager | `VMware1!VMware1!` | required |
| NSXTEA_INSECURE | Specifies whether to skip certificate validation or not. Set this variable to `true`, if you are using self-signed certificates. The default value is `false`  | `true` | optional |


## Usage
Simply type `nsxtea --help` to get help about `nsxtea`'s usage

```bash
Usage:
  nsxtea [command]

Available Commands:
  help        Help about any command
  search      Interact with the Policy Search API

Flags:
      --config string   config file (default is $HOME/.nsxtea.yaml)
  -h, --help            help for nsxtea

Use "nsxtea [command] --help" for more information about a command.
```

### Search command
This command lets you interact with the policy search API. It's a powerful tool to find all kinds of policy objects and filter results. The command returns a JSON string, which can be used for further processing.

```bash
Usage:
  nsxtea search <query> [flags]

Flags:
  -c, --cursor string            Opaque cursor to be used for getting next page of records (supplied by current result page)
  -h, --help                     help for search
  -f, --included_fields string   Comma separated list of fields that should be included in query result
  -p, --page_size string         Maximum number of results to return in this page
                                 Min: 0, Max: 1000 (default "1000")
  -a, --sort_ascending           Sorting order of the query results (default true)
  -s, --sort_by string           Field by which records are sorted

Global Flags:
      --config string   config file (default is $HOME/.nsxtea.yaml)
```

#### Query Syntax
The Query Syntax is exactly the same as described in the offical [VMware NSX-T API docs](https://code.vmware.com/apis/1083/nsx-t) and is copied below for reference.
For convenience, the Query Syntax is also returned by calling

```
nsxtea search --help
```
A query is broken up into **terms** and **operators**.
A **term** is case insensitive and can be a single word such as "Hello" or " World " or a phrase surrounded by double quotes such as "Hello World", which would search for the exact phrase.

##### FIELD NAMES
By default, all the fields will be searched for the search term specified.
Specific fields to be searched on can be provided via the field name followed by a colon ":" and then the search term.
- To search for all the entities where display_name field is "App-VM-1", use `display_name:App-VM-1`
- Use the dot notation to search on nested fields `tags.scope:prod`

##### WILDCARDS
Wildcard searches can be run using `?` to substitute a single character or `*` to substitute zero or more characters
* `*vm*` will match all the entities that contains the term "vm" in any of its fields
* `display_name:App-VM-?` will match App-VM-1, App-VM-2, App-VM-A etc..
* `display_name:App*` will match everything where display_name begins with App  

_Warning_: Be aware that using wildcards especially at the beginning of a word i.e. `*vm` can use a large amount of memory and may perform badly.

##### BOOLEAN OPERATORS
Search terms can be combined using boolean operators `AND`, `OR` and `NOT`. (Note: Boolean operators must be ALL CAPS).
  
###### AND
The `AND` ( && ) operator matches entities where both terms exist in any of the fields of an entity.
To search for Firewall rule with display_name containing block, use `display_name:*block* AND resource_type:FirewallRule`

###### OR
The `OR` ( || ) operator links two terms and finds matching entities if either of the terms exists in an entity.
To search for Firewall rule with `display_name` containing either _block_ or _allow_, use
* `display_name:*block* OR display_name:*allow* AND resource_type:FirewallRule`
* `display_name:(*block* OR *allow*) AND resource_type:FirewallRule`

###### NOT
The `NOT` ( ! ) operator excludes entities that contain the term after NOT. To search for Firewall rule with _display_name_ does not contain the term _block_
* `NOT display_name:*block* AND resource_type:FirewallRule`
* `!display_name:*block* AND resource_type:FirewallRule`

###### RANGES
Ranges can be specified for numeric or string fields and use the following syntax
* `vni:>50001`
* `vni:>=50001`
* `vni:<90000`
* `vni:<=90000` 

To combine an upper and lower bound, you would need to join two clauses with `AND` operator:
* `vni:(>=50001 AND <90000)`

#### RESERVED CHARACTERS
If characters which function as operators are to be used in the query (not as operators), then they should be escaped with a leading backslash.
To search for _(a+b)=c_
* `\(a\+b\)\=c.` 

The reserved characters are: `+ - = && || > < ! ( ) { } [ ] ^ " ~ * ? : \ /`  

Failing to escape these reserved characters correctly would lead to syntax errors and prevent the query from running.


### Image Credits
* Gopher: [Maria Letta - Free Gophers Pack](https://github.com/MariaLetta/free-gophers-pack)
* Teacup: [RROOK, NL](https://thenounproject.com/term/cup-of-tea/2870740/)

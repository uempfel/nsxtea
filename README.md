![nsxtea logo](https://user-images.githubusercontent.com/28337775/106465522-593f6480-649a-11eb-9262-bf94e3add74a.png)

Nsxtea is a small CLI to interact with the VMware NSX-T Search API. Sip a cup of tea while others search in the GUI!

## NSX-T Version compatibility
The CLI will only work with NSX-T versions `3.0` and above. Unfortunately, the API endpoints `nsxtea` relies on are not available in previous versions.  

## Installation
The CLI can be installed via two methods.

### Installation via binary
Binaries for Linux, macOS and Windows are available on the [repos's release page](https://github.com/uempfel/nsxtea/releases).  
To the install `nsxtea` this way, follow these steps:

```bash
# Set a variable to the release version you want to download
export NSXTEA_VERSION=0.3.0
# Download the release for your platform (macOS in this example)
curl -L https://github.com/uempfel/nsxtea/releases/download/v${NSXTEA_VERSION}/nsxtea_${NSXTEA_VERSION}_Darwin_x86_64.tar.gz -o nsxtea.tar.gz

# Unpack the compressed folder 
tar -xvzf nsxtea.tar.gz
x LICENSE
x README.md
x nsxtea
# Move the binary to your PATH
mv nsxtea /usr/local/bin
```

### Installation via go
Assuming you have already [installed go](https://golang.org/doc/install):

```sh
export GO111MODULE=on
go get github.com/uempfel/nsxtea
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
Currently, nsxtea supports three commands, which are documented below:  
* [search](#search-command)
* [apply](#apply-command)
* [curl](#curl-command)

Simply type `nsxtea --help` to get help about `nsxtea`'s usage

```bash
To configure the CLI, set the following environment variables:

NSXTEA_URL (required)
NSXTEA_USERNAME (required)
NSXTEA_PASSWORD (required)
NSXTEA_INSECURE (optional, default false)

Usage:
  nsxtea [command]

Available Commands:
  apply       Interact with the Hierarchical Policy API
  curl        Interact with any API endpoint
  help        Help about any command
  search      Interact with the Policy or Manager Search API

Flags:
  -h, --help   help for nsxtea

Use "nsxtea [command] --help" for more information about a command.
```

### Search command
This command lets you interact with the Policy or Manager search API. It's a powerful tool to find all kinds of objects and filter results. The command returns a JSON string, which can be used for further processing.  
By default, the command searches for Policy objects. If you would like to search for Manager objects instead, simply use the `--manager` flag.

```bash
Usage:
  nsxtea search <query> [flags]

Flags:
  -c, --cursor string            Opaque cursor to be used for getting next page of records (supplied by current result page)
  -h, --help                     help for search
  -f, --included_fields string   Comma separated list of fields that should be included in query result
  -m, --manager                  Use the Manager API for the search request
  -p, --page_size string         Maximum number of results to return in this page 
                                 Min: 0, Max: 1000 (default "1000")
  -a, --sort_ascending           Sorting order of the query results (default true)
  -s, --sort_by string           Field by which records are sorted
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


### Apply command
This command lets you interact with Hierarchical Policy API. It's enables declarative creation, updates and deletion of Policy Objects. To get an overview on how to create objects with the Hierarchical API, this [VMware Blog by Madhukar Krishnarao](https://blogs.vmware.com/networkvirtualization/2020/06/navigating-nsxt-policy-apis.html/) is highly recommended.  
Please refer to [the official API docs](https://vdc-download.vmware.com/vmwb-repository/dcr-public/d6de7a5e-636f-4677-8dbd-6f4ba91fa5e0/36b4881c-41cd-4c46-81d1-b2ca3a6c693b/api_includes/method_PatchInfra.html) for all objects that can be created.

```bash
Decalaratively apply configurations via yaml or json
files.

Examples:
nsxtea apply -f infra.yaml
nsxtea apply -f infra.json

Usage:
  nsxtea apply [flags]

Flags:
  -f, --filepath string   Path to the file that contains the configuration to apply
  -h, --help              help for apply
```

#### Example
The following example is taken from the Blog referenced above. It creates a Tier0 and a connected Tier1 Router in one call.

1) Create a file containing the Objects you want to create. The `apply` command accepts json and yaml input.
* Example as yaml: `infra.yaml`
```yaml
resource_type: Infra
display_name: infra
children:
  - resource_type: ChildTier1
    marked_for_delete: "false"
    Tier1:
      resource_type: Tier1
      display_name: my-Tier-1-GW-Prod
      id: my-Tier-1-GW-Prod
      tier0_path: /infra/tier-0s/Tier-0-GW-West-01
  - resource_type: ChildTier0
    marked_for_delete: "false"
    Tier0:
      resource_type: Tier0
      display_name: Tier-0-GW-West-01-Disconnected
      id: Tier-0-GW-West-01
```

* The same objects as json (taken from [Madhukar Krishnarao's Blog](https://blogs.vmware.com/networkvirtualization/2020/06/navigating-nsxt-policy-apis.html/)): `infra.json`
```json
{
  "resource_type": "Infra",
  "display_name": "infra",
  "children": [
    {
      "resource_type": "ChildTier1",
      "marked_for_delete": "false",
      "Tier1": {
        "resource_type": "Tier1",
        "display_name": "my-Tier-1-GW-Prod",
        "id": "my-Tier-1-GW-Prod",
        "tier0_path": "/infra/tier-0s/Tier-0-GW-West-01"
      }
    },
    {
      "resource_type": "ChildTier0",
      "marked_for_delete": "false",
      "Tier0": {
        "resource_type": "Tier0",
        "display_name": "Tier-0-GW-West-01-Disconnected",
        "id": "Tier-0-GW-West-01"
      }
    }
  ]
}
````

2) Run the apply command providing the path to the file you created
```bash
# Apply as yaml
nsxtea apply -f path/to/infra.yaml
# Apply as json
nsxtea apply -f path/to/infra.json
```

That's it! The objects should be created and be available after a short time.  

#### Updating and deleting objects
To update objects, simply adapt the file with the necessary configuration and re-run `nsxtea apply`.  

Deleting the objects created in the example above is as simple as changing the `marked_for_delete` properties from `false` to `true`. Once you've done that, simply re-run `nsxtea apply` and the objects should be deleted after a short time.

### Curl command
Sometimes you need full control over the API. This is where the `nsxtea curl` command comes in. It enables you to interact with any endpoint documented in the official docs.  

```bash
Interact with any API endpoint

Examples:
nsxtea curl -X PATCH /policy/api/v1/infra -d @path-to-body-file

Usage:
  nsxtea curl <endpoint> [flags]

Flags:
  -d, --data string     Body data. You can specifiy a path to a yaml or json file with the '@' prefix
  -h, --help            help for curl
  -X, --method string   HTTP Method (default "GET")
  -o, --override        Add the 'X-Allow-Overwrite: true' header to mutate protected objects
```

### Image Credits
* Gopher: [Maria Letta - Free Gophers Pack](https://github.com/MariaLetta/free-gophers-pack)
* Teacup: [RROOK, NL](https://thenounproject.com/term/cup-of-tea/2870740/)

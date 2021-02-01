![nsxtea logo](https://user-images.githubusercontent.com/28337775/106447575-f3df7980-6481-11eb-8cba-df7fd97d3c17.png)

Nsxtea is a small CLI to interact with the VMware NSX-T Policy Search API.

## Installation
Assuming you have already [installed go](https://golang.org/doc/install):

```sh
go get github.com/uempfel/nsxtea
go install github.com/uempfel/nsxtea
```

If you're getting  a `command not found` error, you might have to add your "go binary" directory to your PATH. To do so, run:
```sh 
export PATH=$PATH:$(go env GOPATH)/bin
```

### Image Credits
* Gopher: [Maria Letta - Free Gophers Pack](https://github.com/MariaLetta/free-gophers-pack)
* Teacup: [RROOK, NL](https://thenounproject.com/term/cup-of-tea/2870740/)
# Magic Admin Golang SDK

The Magic Admin Golang SDK provides convenient ways for developers to interact with Magic API endpoints and a variety of utilities to handle [DID Token](https://magic.link/docs/introduction/decentralized-id).

> **Note:** This repository is a fork maintained by community volunteers under [github.com/yukia3e/magic-admin-go]. The original project by magiclabs is no longer actively maintained.

## Table of Contents

* [Documentation](#documentation)
* [Quick Start](#quick-start)
* [Installation](#installation)
* [Command Line Utility](#command-line-utility)
* [Development](#development)
* [Changelog](#changelog)
* [License](#license)

## Documentation

See the [Magic doc](https://magic.link/docs/api-reference/server-side-sdks/go)!

## Installation

The SDK requires **Go 1.24.0** and Go Modules. To make sure your project is using Go Modules, check for a `go.mod` file in your project's root directory. If it exists, you're already set up; if not, please follow [this guide](https://blog.golang.org/migrating-to-go-modules) to migrate to Go Modules.

Simply reference the SDK in your Go program with an `import`:

```golang
import (
// ...
"github.com/yukia3e/magic-admin-go"
// ...
)
```

When you run any standard Go command (ex: `build` or `install`), the Go toolchain will automatically fetch the SDK for you.

Alternatively, you can explicitly add the package to your project with:

```sh
go get github.com/yukia3e/magic-admin-go
```

## Command Line Utility

A command line utility is provided for testing purposes and can be used to decode and validate DID tokens, as well as to retrieve user info.

You can install it using:

```bash
go install github.com/yukia3e/magic-admin-go/cmd/magic-cli
```

The currently available commands are:

```bash
$ magic-cli -h
NAME:
   magic-cli - command line utility to make requests to API and validate tokens

USAGE:
   magic-cli [global options] command [command options] [arguments...]

COMMANDS:
   token, t   magic-cli token [decode|validate] --did <DID token> [--clientId <Magic Client ID>]
   user, u    magic-cli -s <secret> user --did <DID token>
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --secret value, -s value  Secret token used for making requests to the backend API [$MAGIC_API_SECRET_KEY]
   --help, -h                show help (default: false)
```

## Quick Start

Before you start, you will need an API secret key which you can obtain from the [Magic Dashboard](https://dashboard.magic.link/). Once you have your API secret key, instantiate a Magic object as shown below.

Sample code to retrieve user info using a [DID token](https://docs.magic.link/decentralized-id):

```golang
package main

import (
	"fmt"
	"log"

	"github.com/yukia3e/magic-admin-go"
	"github.com/yukia3e/magic-admin-go/client"
)

func main() {
	m, err := client.New("<YOUR_API_SECRET_KEY>", magic.NewDefaultClient())
	if err != nil {
		log.Fatalf("Error initializing client: %s", err.Error())
	}
	userInfo, err := m.User.GetMetadataByToken("<DID_TOKEN>")
	if err != nil {
		log.Fatalf("Error retrieving user metadata: %s", err.Error())
	}

	fmt.Println(userInfo)
}
```

Sample code to validate a [DID token](https://docs.magic.link/decentralized-id) and retrieve the `claim` and `proof`:

```golang
package main

import (
	"fmt"
	"log"

	"github.com/yukia3e/magic-admin-go/client"
	"github.com/yukia3e/magic-admin-go/token"
	"github.com/yukia3e/magic-admin-go"
)

func main() {
	c, err := client.New("<YOUR_API_SECRET_KEY>", magic.NewDefaultClient())
	if err != nil {
		log.Fatalf("Unable to initialize client: %s", err.Error())
	}

	tk, err := token.NewToken("<DID_TOKEN>")
	if err != nil {
		log.Fatalf("DID token is malformed: %s", err.Error())
	}

	if err := tk.Validate(c.ClientInfo.ClientId); err != nil {
		log.Fatalf("DID token is invalid: %v", err)
	}

	fmt.Println(tk.GetClaim())
	fmt.Println(tk.GetProof())
}
```

### Configure Network Strategy

The `NewClientWithRetry` method creates a client with customizable `retries`, `retryWait`, and `timeout` options. It returns a `*resty.Client` instance that can be used with the Magic client:

```golang
cl := magic.NewClientWithRetry(5, time.Second, 10*time.Second)
m := client.New("<YOUR_API_SECRET_KEY>", cl)
```

## Development

We welcome contributions to the SDK. To get started:

1. **Clone the repository**  
   Forked from the original magiclabs repository, this project is now maintained by community volunteers.

2. **Install dependencies and run tests**  
   First, install the **mise** package by running:
   ```bash
   mise install
   ```
   Then, run the following command to execute existing tests:
   ```bash
   make test
   ```

3. **Build and install the magic-cli utility**  
   Use:
   ```bash
   make install
   ```
   or build the binary directly with:
   ```bash
   make build
   ```

For more details on contributing, please see our [CONTRIBUTING](CONTRIBUTING.md) guide.

## Changelog

See [Changelog](CHANGELOG.md) for a history of changes.

## License

See [License](LICENSE.txt) for license details.
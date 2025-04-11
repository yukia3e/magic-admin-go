package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/urfave/cli/v2"

	"github.com/yukia3e/magic-admin-go"
	"github.com/yukia3e/magic-admin-go/client"
	"github.com/yukia3e/magic-admin-go/token"
)

func main() {
	app := &cli.App{
		Name:     "magic-cli",
		Usage:    "command line utility to make requests to api and validate tokens",
		Compiled: time.Now(),
		Commands: []*cli.Command{
			{
				Name:    "token",
				Aliases: []string{"t"},
				Usage:   "magic-cli token [decode|validate] --did <DID token> [--clientId <Magic Client ID>]",
				Subcommands: []*cli.Command{
					{
						Name:  "decode",
						Usage: "magic-cli token decode --did <DID token>",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "did",
								Usage: "Did token which must be decoded",
							},
						},
						Action: decodeDIDToken,
					},
					{
						Name:  "validate",
						Usage: "magic-cli token validate --did <DID token> --clientId <Magic Client ID>",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:  "did",
								Usage: "Did token which must be validated",
							},
							&cli.StringFlag{
								Name:    "clientId",
								Usage:   "Magic Client ID to validate the aud field",
								EnvVars: []string{"MAGIC_CLIENT_ID"},
							},
						},
						Action: validateDIDToken,
					},
				},
			},
			{
				Name:    "user",
				Aliases: []string{"u"},
				Usage:   "magic-cli -s <secret> user --did <DID token>",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "did",
						Usage: "Did token used for user info receiving",
					},
				},
				Action: userMetadata,
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "secret",
				Usage:   "Secret token which will be used for making request to backend api",
				Aliases: []string{"s"},
				EnvVars: []string{"MAGIC_API_SECRET_KEY"},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func userMetadata(c *cli.Context) error {
	m, err := client.New(c.String("secret"), magic.NewDefaultClient())

	if err != nil {
		return err
	}

	userInfo, err := m.User.GetMetadataByToken(c.String("did"))
	if err != nil {
		return err
	}

	fmt.Println(userInfo.String())

	return nil
}

func decodeDIDToken(c *cli.Context) error {
	tk, err := token.NewToken(c.String("did"))
	if err != nil {
		return err
	}

	claim := tk.GetClaim()
	fmt.Println(claim.String())

	return nil
}

func validateDIDToken(c *cli.Context) error {

	tk, err := token.NewToken(c.String("did"))
	if err != nil {
		return err
	}

	if err := tk.Validate(c.String("clientId")); err != nil {
		return err
	}

	return nil
}

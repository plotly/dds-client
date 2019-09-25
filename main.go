package main

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/akamensky/argparse"
	"github.com/machinebox/graphql"
)

type AppsResponse struct {
	AppsWrapper AppsWrapper `json:"apps"`
}

type AppsWrapper struct {
	Apps     []App `json:"apps"`
	NextPage int   `json:"nextPage"`
}

type App struct {
	Name string `json:"name"`
}

type AddAppResponse struct {
	AddApp AddApp `json:"addApp"`
}

type DeleteAppResponse struct {
	DeleteApp DeleteApp `json:"deleteApp"`
}

type AddApp struct {
	App   App    `json:"app"`
	Error string `json:"error"`
}

type DeleteApp struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

var ClientURL = os.Getenv("DASH_API_URL")
var Username = os.Getenv("DASH_API_USER")
var ApiKey = os.Getenv("DASH_API_KEY")

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func appsList() {
	req := graphql.NewRequest(`
query apps($allApps: Boolean) {
  apps(allApps: $allApps) {
    nextPage
    apps {
      name
    }
  }
}
`)

	req.Var("allApps", true)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Basic "+basicAuth(Username, ApiKey))

	ctx := context.Background()

	var respData AppsResponse

	client := graphqlClient()
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	for _, app := range respData.AppsWrapper.Apps {
		fmt.Printf("%v\n", app.Name)
	}
}

func appsCreate(name string) {
	if name == "" {
		log.Fatal(errors.New("No name specified"))
	}

	req := graphql.NewRequest(`
mutation AddApp($name: String!) {
  addApp(name: $name) {
    app {
      name
    }
    error
  }
}
`)

	req.Var("name", name)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Basic "+basicAuth(Username, ApiKey))

	ctx := context.Background()

	var respData AddAppResponse

	client := graphqlClient()
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	if respData.AddApp.Error != "" {
		fmt.Printf(" !    %v!\n", respData.AddApp.Error)
	} else {
		fmt.Printf("====> %v created!\n", respData.AddApp.App.Name)
	}
}

func appsDelete(name string) {
	if name == "" {
		log.Fatal(errors.New("No name specified"))
	}

	req := graphql.NewRequest(`
mutation DeleteApp($name: String!) {
  deleteApp(name: $name) {
    ok
    error
  }
}
`)

	req.Var("name", name)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Basic "+basicAuth(Username, ApiKey))

	ctx := context.Background()

	var respData DeleteAppResponse

	client := graphqlClient()
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	if respData.DeleteApp.Error != "" {
		fmt.Printf(" !    %v!\n", respData.DeleteApp.Error)
	} else {
		fmt.Printf("====> %v deleted!\n", name)
	}
}

func graphqlClient() *graphql.Client {
	httpclient := &http.Client{}
	if true {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpclient = &http.Client{Transport: tr}
	}

	return graphql.NewClient(ClientURL, graphql.WithHTTPClient(httpclient))
}

func main() {
	parser := argparse.NewParser("dds-client", "A simple dds client")
	name := parser.String("", "name", &argparse.Options{Help: "Name of app"})

	appsListCmd := parser.NewCommand("apps:list", "List all apps")
	appsCreateCmd := parser.NewCommand("apps:create", "Create an app")
	appsDeleteCmd := parser.NewCommand("apps:delete", "Delete an app")
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if appsListCmd.Happened() {
		appsList()
	} else if appsCreateCmd.Happened() {
		appsCreate(*name)
	} else if appsDeleteCmd.Happened() {
		appsDelete(*name)
	} else {
		err := fmt.Errorf("bad arguments, please check usage")
		fmt.Print(parser.Usage(err))
	}
}

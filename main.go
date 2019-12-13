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
	"sort"

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

var ApiKey = os.Getenv("DASH_ENTERPRISE_API_KEY")
var DashEnterpriseURL = os.Getenv("DASH_ENTERPRISE_URL")
var Username = os.Getenv("DASH_ENTERPRISE_USERNAME")

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}

func appExists(name string) {
	req := graphql.NewRequest(`
query apps($name: String!, $allApps: Boolean!) {
  apps(name: $name, allApps: $allApps) {
    apps {
      name
    }
  }
}
`)

	req.Var("allApps", false)
	req.Var("name", name)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Basic "+basicAuth(Username, ApiKey))

	ctx := context.Background()

	var respData AppsResponse

	client, err := graphqlClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}
	for _, app := range respData.AppsWrapper.Apps {
		fmt.Printf("%v exists\n", app.Name)
		os.Exit(0)
	}
	fmt.Printf("%v does not exist\n", name)
	os.Exit(1)
}

func appsList() {
	req := graphql.NewRequest(`
query apps($page: Int!, $allApps: Boolean!) {
  apps(page: $page, allApps: $allApps) {
    apps {
      name
    }
    nextPage
  }
}
`)

	req.Var("allApps", true)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Basic "+basicAuth(Username, ApiKey))

	ctx := context.Background()

	var respData AppsResponse

	client, err := graphqlClient()
	if err != nil {
		log.Fatal(err)
	}

	var apps []string
	page := 1
	for true {
		if page == 0 {
			break
		}

		req.Var("page", page)
		if err := client.Run(ctx, req, &respData); err != nil {
			log.Fatal(err)
		}
		for _, app := range respData.AppsWrapper.Apps {
			apps = append(apps, app.Name)
		}

		if page == respData.AppsWrapper.NextPage {
			break
		}
		page = respData.AppsWrapper.NextPage
	}

	sort.Strings(apps)
	for _, app := range apps {
		fmt.Printf("%v\n", app)
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

	client, err := graphqlClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	if respData.AddApp.Error != "" {
		fmt.Printf(" !    %v\n", respData.AddApp.Error)
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

	client, err := graphqlClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	if respData.DeleteApp.Error != "" {
		fmt.Printf(" !    %v\n", respData.DeleteApp.Error)
	} else {
		fmt.Printf("====> %v deleted!\n", name)
	}
}

func graphqlClient() (client *graphql.Client, err error) {
	httpclient := &http.Client{}
	if true {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpclient = &http.Client{Transport: tr}
	}

	if DashEnterpriseURL == "" {
		return client, errors.New("DASH_ENTERPRISE_URL environment variable not defined")
	}

	if Username == "" {
		return client, errors.New("DASH_ENTERPRISE_USERNAME environment variable not defined")
	}

	if ApiKey == "" {
		return client, errors.New("DASH_ENTERPRISE_API_KEY environment variable not defined")
	}

	client = graphql.NewClient(DashEnterpriseURL+"/Manager/graphql", graphql.WithHTTPClient(httpclient))
	return client, err
}

func main() {
	parser := argparse.NewParser("dds-client", "A simple dds client")
	name := parser.String("", "name", &argparse.Options{Help: "Name of app"})

	appsListCmd := parser.NewCommand("apps:list", "List all apps")
	appsCreateCmd := parser.NewCommand("apps:create", "Create an app")
	appsDeleteCmd := parser.NewCommand("apps:delete", "Delete an app")
	appExistsCmd := parser.NewCommand("apps:exists", "Check if an app exists")
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
	} else if appExistsCmd.Happened() {
		appExists(*name)
	} else {
		err := fmt.Errorf("bad arguments, please check usage")
		fmt.Print(parser.Usage(err))
	}
}

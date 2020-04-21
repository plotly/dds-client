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

type ServicesResponse struct {
	Services []Service `json:"services"`
}

type Service struct {
	Name        string `json:"name"`
	ServiceType string `json:"serviceType"`
	Created     string `json:"created"`
}

type AddServiceResponse struct {
	AddService AddService `json:"addService"`
}

type DeleteServiceResponse struct {
	DeleteService DeleteService `json:"deleteService"`
}

type LinkServiceResponse struct {
	LinkService LinkService `json:"linkService"`
}

type UnlinkServiceResponse struct {
	UnlinkService UnlinkService `json:"linkService"`
}

type AddService struct {
	Service Service `json:"app"`
	Error   string  `json:"error"`
}

type DeleteService struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type LinkService struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type UnlinkService struct {
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

func postgresCreate(name string) {
	serviceCreate("postgres", name)
}

func postgresDelete(name string) {
	serviceDelete("postgres", name)
}

func postgresExists(name string) {
	serviceExists("postgres", name)
}

func postgresLink(serviceName string, appName string) {
	serviceLink("postgres", serviceName, appName)
}

func postgresList() {
	serviceList("postgres")
}

func postgresUnlink(serviceName string, appName string) {
	serviceUnlink("postgres", serviceName, appName)
}

func redisCreate(name string) {
	serviceCreate("redis", name)
}

func redisDelete(name string) {
	serviceDelete("redis", name)
}

func redisExists(name string) {
	serviceExists("redis", name)
}

func redisLink(serviceName string, appName string) {
	serviceLink("redis", serviceName, appName)
}

func redisList() {
	serviceList("redis")
}

func redisUnlink(serviceName string, appName string) {
	serviceUnlink("redis", serviceName, appName)
}

func serviceCreate(serviceType string, name string) {
	if name == "" {
		log.Fatal(errors.New("No name specified"))
	}

	mutation := `
mutation AddService($name: String!, $serviceType: ServiceType = %s) {
  addService(name: $name, serviceType: $serviceType) {
    service {
      name
      serviceType
      created
    }
    error
  }
}
`
	req := graphql.NewRequest(fmt.Sprintf(mutation, serviceType))

	req.Var("name", name)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Basic "+basicAuth(Username, ApiKey))

	ctx := context.Background()

	var respData AddServiceResponse

	client, err := graphqlClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	if respData.AddService.Error != "" {
		fmt.Printf(" !    %v\n", respData.AddService.Error)
	} else {
		fmt.Printf("====> %v created!\n", name)
	}
}

func serviceDelete(serviceType string, name string) {
	if name == "" {
		log.Fatal(errors.New("No name specified"))
	}

	mutation := `
mutation DeleteService($name: String!, $serviceType: ServiceType = %s) {
  deleteService(name: $name, serviceType: $serviceType) {
    ok
    error
  }
}
`
	req := graphql.NewRequest(fmt.Sprintf(mutation, serviceType))

	req.Var("name", name)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Basic "+basicAuth(Username, ApiKey))

	ctx := context.Background()

	var respData DeleteServiceResponse

	client, err := graphqlClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	if respData.DeleteService.Error != "" {
		fmt.Printf(" !    %v\n", respData.DeleteService.Error)
	} else {
		fmt.Printf("====> %v deleted!\n", name)
	}
}

func serviceExists(serviceType string, name string) {
	if name == "" {
		log.Fatal(errors.New("No name specified"))
	}

	respData, err := fetchServices()
	if err != nil {
		log.Fatal(err)
	}

	for _, service := range respData.Services {
		if service.Name == name && service.ServiceType == serviceType {
			fmt.Printf("%v exists\n", service.Name)
			os.Exit(0)
		}
	}
	fmt.Printf("%v not found. Possible causes:\n", name)
	fmt.Printf("- You may not have been granted access to this service.")
	fmt.Printf("- The service may not exist.")
	os.Exit(1)
}

func serviceList(serviceType string) {
	respData, err := fetchServices()
	if err != nil {
		log.Fatal(err)
	}

	for _, service := range respData.Services {
		if service.ServiceType == serviceType {
			fmt.Printf("%v\n", service.Name)
		}
	}
}

func serviceLink(serviceType string, serviceName string, appName string) {
	if serviceName == "" {
		log.Fatal(errors.New("No name specified"))
	}

	if appName == "" {
		log.Fatal(errors.New("No app specified"))
	}

	mutation := `
mutation LinkService($appname: String!, $serviceName: String!, $serviceType: ServiceType = %s) {
  linkService(appname: $appname, serviceType: $serviceType, serviceName: $serviceName) {
    ok
    error
  }
}
`
	req := graphql.NewRequest(fmt.Sprintf(mutation, serviceType))

	req.Var("appname", appName)
	req.Var("serviceName", serviceName)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Basic "+basicAuth(Username, ApiKey))

	ctx := context.Background()

	var respData LinkServiceResponse

	client, err := graphqlClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	if respData.LinkService.Error != "" {
		fmt.Printf(" !    %v\n", respData.LinkService.Error)
	} else {
		fmt.Printf("====> %v linked!\n", appName)
	}
}

func serviceUnlink(serviceType string, serviceName string, appName string) {
	if serviceName == "" {
		log.Fatal(errors.New("No name specified"))
	}

	if appName == "" {
		log.Fatal(errors.New("No app specified"))
	}

	mutation := `
mutation UnlinkService($appname: String!, $serviceName: String!, $serviceType: ServiceType = %s) {
  unlinkService(appname: $appname, serviceType: $serviceType, serviceName: $serviceName) {
    ok
    error
  }
}
`
	req := graphql.NewRequest(fmt.Sprintf(mutation, serviceType))

	req.Var("appname", appName)
	req.Var("serviceName", serviceName)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Basic "+basicAuth(Username, ApiKey))

	ctx := context.Background()

	var respData UnlinkServiceResponse

	client, err := graphqlClient()
	if err != nil {
		log.Fatal(err)
	}
	if err := client.Run(ctx, req, &respData); err != nil {
		log.Fatal(err)
	}

	if respData.UnlinkService.Error != "" {
		fmt.Printf(" !    %v\n", respData.UnlinkService.Error)
	} else {
		fmt.Printf("====> %v unlinked!\n", appName)
	}
}

func fetchServices() (ServicesResponse, error) {
	req := graphql.NewRequest(`
{
    services {
        name
        serviceType
        created
    }
}
`)

	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Authorization", "Basic "+basicAuth(Username, ApiKey))

	ctx := context.Background()

	var respData ServicesResponse

	client, err := graphqlClient()
	if err != nil {
		return respData, err
	}

	err = client.Run(ctx, req, &respData)
	return respData, err
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
	fmt.Printf("%v not found. Possible causes:\n", name)
	fmt.Printf("- You may not have been granted access to this app.")
	fmt.Printf("- The app may not exist (or may not have been deployed yet).")
	fmt.Printf("- The app is broken and could not be started.")
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

	appsListCmd := parser.NewCommand("apps:list", "List all apps")

	appsCreateCmd := parser.NewCommand("apps:create", "Create an app")
	appsCreateCmdName := appsCreateCmd.String("", "name", &argparse.Options{Help: "Name of app"})

	appsDeleteCmd := parser.NewCommand("apps:delete", "Delete an app")
	appsDeleteCmdName := appsDeleteCmd.String("", "name", &argparse.Options{Help: "Name of app"})

	appExistsCmd := parser.NewCommand("apps:exists", "Check if an app exists")
	appExistsCmdName := appExistsCmd.String("", "name", &argparse.Options{Help: "Name of app"})

	postgresCreateCmd := parser.NewCommand("postgres:create", "Create a postgres service")
	postgresCreateCmdName := postgresCreateCmd.String("", "name", &argparse.Options{Help: "Name of service"})

	postgresDeleteCmd := parser.NewCommand("postgres:delete", "Delete a postgres service")
	postgresDeleteCmdName := postgresDeleteCmd.String("", "name", &argparse.Options{Help: "Name of service"})

	postgresExistsCmd := parser.NewCommand("postgres:exists", "Check if a postgres service exists")
	postgresExistsCmdName := postgresExistsCmd.String("", "name", &argparse.Options{Help: "Name of service"})

	postgresLinkCmd := parser.NewCommand("postgres:link", "Link a postgres service to an app")
	postgresLinkCmdName := postgresLinkCmd.String("", "name", &argparse.Options{Help: "Name of service"})
	postgresLinkCmdApp := postgresLinkCmd.String("", "app", &argparse.Options{Help: "Name of app"})

	postgresListCmd := parser.NewCommand("postgres:list", "List all postgres services")

	postgresUnlinkCmd := parser.NewCommand("postgres:unlink", "Unlink a postgres service to an app")
	postgresUnlinkCmdName := postgresUnlinkCmd.String("", "name", &argparse.Options{Help: "Name of service"})
	postgresUnlinkCmdApp := postgresUnlinkCmd.String("", "app", &argparse.Options{Help: "Name of app"})

	redisCreateCmd := parser.NewCommand("redis:create", "Create a redis service")
	redisCreateCmdName := redisCreateCmd.String("", "name", &argparse.Options{Help: "Name of service"})

	redisDeleteCmd := parser.NewCommand("redis:delete", "Delete a redis service")
	redisDeleteCmdName := redisDeleteCmd.String("", "name", &argparse.Options{Help: "Name of service"})

	redisExistsCmd := parser.NewCommand("redis:exists", "Check if a redis service exists")
	redisExistsCmdName := redisExistsCmd.String("", "name", &argparse.Options{Help: "Name of service"})

	redisLinkCmd := parser.NewCommand("redis:link", "Link a redis service to an app")
	redisLinkCmdName := redisLinkCmd.String("", "name", &argparse.Options{Help: "Name of service"})
	redisLinkCmdApp := redisLinkCmd.String("", "app", &argparse.Options{Help: "Name of app"})

	redisListCmd := parser.NewCommand("redis:list", "List all redis services")

	redisUnlinkCmd := parser.NewCommand("redis:unlink", "Unlink a redis service to an app")
	redisUnlinkCmdName := redisUnlinkCmd.String("", "name", &argparse.Options{Help: "Name of service"})
	redisUnlinkCmdApp := redisUnlinkCmd.String("", "app", &argparse.Options{Help: "Name of app"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		return
	}

	if appsListCmd.Happened() {
		appsList()
	} else if appsCreateCmd.Happened() {
		appsCreate(*appsCreateCmdName)
	} else if appsDeleteCmd.Happened() {
		appsDelete(*appsDeleteCmdName)
	} else if appExistsCmd.Happened() {
		appExists(*appExistsCmdName)
	} else if postgresCreateCmd.Happened() {
		postgresCreate(*postgresCreateCmdName)
	} else if postgresDeleteCmd.Happened() {
		postgresDelete(*postgresDeleteCmdName)
	} else if postgresExistsCmd.Happened() {
		postgresExists(*postgresExistsCmdName)
	} else if postgresLinkCmd.Happened() {
		postgresLink(*postgresLinkCmdName, *postgresLinkCmdApp)
	} else if postgresListCmd.Happened() {
		postgresList()
	} else if postgresUnlinkCmd.Happened() {
		postgresUnlink(*postgresUnlinkCmdName, *postgresUnlinkCmdApp)
	} else if redisCreateCmd.Happened() {
		redisCreate(*redisCreateCmdName)
	} else if redisDeleteCmd.Happened() {
		redisDelete(*redisDeleteCmdName)
	} else if redisExistsCmd.Happened() {
		redisExists(*redisExistsCmdName)
	} else if redisLinkCmd.Happened() {
		redisLink(*redisLinkCmdName, *redisLinkCmdApp)
	} else if redisListCmd.Happened() {
		redisList()
	} else if redisUnlinkCmd.Happened() {
		redisUnlink(*redisUnlinkCmdName, *redisUnlinkCmdApp)
	} else {
		err := fmt.Errorf("bad arguments, please check usage")
		fmt.Print(parser.Usage(err))
	}
}

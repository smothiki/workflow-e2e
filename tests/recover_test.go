package tests

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/deis/workflow-e2e/tests/cmd"
	"github.com/deis/workflow-e2e/tests/cmd/apps"
	"github.com/deis/workflow-e2e/tests/cmd/auth"
	"github.com/deis/workflow-e2e/tests/cmd/git"
	"github.com/deis/workflow-e2e/tests/cmd/helm"
	"github.com/deis/workflow-e2e/tests/cmd/keys"
	"github.com/deis/workflow-e2e/tests/model"
	"github.com/deis/workflow-e2e/tests/settings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func parentdir() {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
	}
	dir := filepath.Dir(pwd)
	os.Chdir(dir)
	os.RemoveAll("example-dockerfile-http")
	os.RemoveAll("example-go")
}

func cleanup() {
	os.RemoveAll("example-dockerfile-http")
	os.RemoveAll("example-go")
}

func curl(endpoint string) chan bool {
	tick := time.Tick(2000 * time.Millisecond)
	timeout := time.Tick(70000 * time.Millisecond)
	ok := make(chan bool)
	go func() {
		for {
			select {
			case <-tick:
				res, _ := http.Get(endpoint)
				if (res.StatusCode >= 200) && (res.StatusCode <= 399) {
					ok <- true
					close(ok)
					return
				}
			case <-timeout:
				ok <- false
				close(ok)
				return
			default:
			}
		}
	}()
	return ok
}

func waitforcontroller() (string, error) {
	IP, err := cmd.Execute(`kubectl get service -o jsonpath={.status.loadBalancer.ingress[0].ip}  deis-router --namespace=deis`)
	if err != nil {
		return IP, err
	}
	ok := curl("http://deis." + IP + ".nip.io/healthz")
	status := <-ok
	if !status {
		return "", fmt.Errorf("took much time to start controller")
	}
	if settings.GetDeisControllerURL() != "deis."+IP+".nip.io" {
		settings.SetDeisControllerURL(IP)
	}
	return IP, nil
}

var _ = Describe("Recover", func() {

	Context("with an existing user and a public key", func() {

		var user model.User
		var keyPath string
		Specify("that user can register and create key", func() {
			user = auth.Register()
			_, keyPath = keys.Add(user)
			cleanup()
		})

		Context("deploys a deis pull, dockerfile, buildpack apps", func() {

			// var app1 model.App
			var app2 model.App
			var app3 model.App

			// Specify("that user can deploy using deis pull", func() {
			// 	app1 = apps.Create(user, "--no-remote")
			// 	builds.Create(user, app1)
			// })
			Specify("that user can deploy buildpack app using a git push", func() {
				output, err := cmd.Execute(`git clone https://github.com/deis/example-go.git`)
				Expect(err).NotTo(HaveOccurred(), output)
				os.Chdir("example-go")
				app2 = apps.Create(user)
				git.Push(user, keyPath)
			})
			Specify("that user can deploy dockerfile app using a git push", func() {
				parentdir()
				output, err := cmd.Execute(`git clone https://github.com/deis/example-dockerfile-http`)
				Expect(err).NotTo(HaveOccurred(), output)
				os.Chdir("example-dockerfile-http")
				app3 = apps.Create(user)
				git.Push(user, keyPath)
			})
			Specify("helm uninstall workflow-dev -y -n  deis", func() {
				helm.Uninstall("workflow-dev")
			})
			Specify("helm install workflow-dev ", func() {
				helm.Install("workflow-dev")
				waitforcontroller()
			})
			Specify("login with existing user", func() {
				auth.Login(user)
			})
			Specify("should open apps", func() {
				auth.Login(user)
			})
			Specify("should be able to curl apps", func() {
				auth.Login(user)
				// app.Open(user,app1)
				apps.Open(user, app2)
				apps.Open(user, app3)
			})
		})
	})
})

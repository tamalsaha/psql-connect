package main

import (
	"fmt"
	"log"
	"path/filepath"

	shell "github.com/codeskyblue/go-sh"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"kmodules.xyz/client-go/tools/portforward"
)

func main() {
	masterURL := ""
	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags(masterURL, kubeconfigPath)
	if err != nil {
		log.Fatalf("Could not get Kubernetes config: %s", err)
	}

	// kubedb connect -it -n demo postgres quick-postgres

	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalln(err)
	}

	ns := "demo"
	podName := "quick-postgres-0"
	port := 5432
	secretName := "quick-postgres-auth"

	tunnel := portforward.NewTunnel(client.RESTClient(), config, ns, podName, port)
	err = tunnel.ForwardPort()
	if err != nil {
		log.Fatalln(err)
	}
	defer tunnel.Close()

	auth, err := client.CoreV1().Secrets(ns).Get(secretName, metav1.GetOptions{})
	if err != nil {
		log.Fatalln(err)
	}

	sh := shell.NewSession()
	sh.SetEnv("PGPASSWORD", string(auth.Data["POSTGRES_PASSWORD"]))
	sh.ShowCMD = true
	err = sh.Command("docker", "run", "--network=host", "-it",
		"postgres:10.2-alpine",
		"psql",
		"--host=127.0.0.1",
		fmt.Sprintf("--port=%d", tunnel.Remote),
		"--username=postgres").Run()
	if err != nil {
		log.Fatalln(err)
	}
}

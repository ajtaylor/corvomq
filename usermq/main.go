package main

import (
	"log"
	"os"
	"os/exec"

	"github.com/minio/minio-go"
)

var (
	accessKey = os.Getenv("SPACES_ACCESS_KEY")
	secKey    = os.Getenv("SPACES_SECURITY_KEY")
)

func runCmd(cmd *exec.Cmd) {
	err := cmd.Run()
	if err != nil {
		log.Println(err)
	}
}

func main() {
	var cmd *exec.Cmd
	var hostname string

	hostname, _ = os.Hostname()

	// Connect to Spaces
	client, err := minio.New("nyc3.digitaloceanspaces.com", accessKey, secKey, true)
	if err != nil {
		log.Printf("Error connecting to Spaces: %v", err.Error())
	}

	// Check user files are ready
	checkFile(client, "cloud-init/"+hostname+"_meta-data")
	checkFile(client, "cloud-init/"+hostname+"_user-data")

	// Comment out crontab for this job
	// cmd = exec.Command("/bin/sh", "-c", "crontab -l | sed 's;^\\*.*\\/opt\\/get_files;# &;' | crontab -")
	// runCmd(cmd)

	// Remove cloud-init installation instance
	cmd = exec.Command("/bin/sh", "-c", "rm -rf /var/lib/cloud/instance/*")
	runCmd(cmd)

	// Get files from Spaces
	getFile(client, "meta-data",
		"cloud-init/"+hostname+"_meta-data",
		"/var/lib/cloud/seed/nocloud-net/meta-data")

	getFile(client, "user-data",
		"cloud-init/"+hostname+"_user-data",
		"/var/lib/cloud/seed/nocloud-net/user-data")

	getFile(client, "gnatsd",
		"download/nats/gnatsd-v1.0.2-linux-amd64.zip",
		"/srv/nats/gnatsd-v1.0.2-linux-amd64.zip")

	getFile(client, "nats-exporter",
		"download/prometheus/prometheus-nats-exporter.gz",
		"/srv/prometheus/nats_exporter/prometheus-nats-exporter.gz")

	getFile(client, "node-exporter",
		"download/prometheus/node_exporter-0.14.0.linux-amd64.tar.gz",
		"/srv/prometheus/nats_exporter/node_exporter-0.14.0.linux-amd64.tar.gz")

	// Reboot
	cmd = exec.Command("/bin/sh", "-c", "reboot")
	runCmd(cmd)
}

func getFile(client *minio.Client, label string, source string, destination string) {
	log.Printf("Getting %v...", label)

	err := client.FGetObject("corvomq", source, destination, minio.GetObjectOptions{})
	if err != nil {
		log.Printf("Error getting %v: %v", label, err.Error())
	}
}

func checkFile(client *minio.Client, source string) {
	log.Printf("Checking %v...", source)

	_, err := client.StatObject("corvomq", source, minio.StatObjectOptions{})
	if err != nil {
		log.Fatalf("File check error: %v, %v", source, err.Error())
	}
}

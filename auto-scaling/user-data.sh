#!/bin/bash
sudo yum update -y
sudo yum install golang -y
mkdir /home/ec2-user/go
cd /home/ec2-user/go
cat <<EOT >> main.go
package main

import (
	"io"
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", Get)
	http.ListenAndServe(":80", nil)
}

func Get(responseWriter http.ResponseWriter, _ *http.Request) {
	instanceID, err := GetInstanceID()
	if err != nil {
		responseWriter.WriteHeader(500)
		responseWriter.Write([]byte(ToHTML(err.Error())))
		return
	}
	responseWriter.Write([]byte(ToHTML(fmt.Sprintf("Hello from %s\n", instanceID))))
}

func GetInstanceID() (string, error) {
	response, err := http.Get("http://169.254.169.254/latest/meta-data/instance-id")
	if err != nil {
		return "", fmt.Errorf("could not send request: %w", err)
	}
	instanceIDBytes, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("could not read response body: %w", err)
	}
	return string(instanceIDBytes), nil
}

func ToHTML(value string) string {
	return fmt.Sprintf(
		"<!DOCTYPE html>" +
			"<html>\n" +
			"  <head>\n" +
			"    <title>AWS SAA Labs</title>\n" +
			"  </head>\n" +
			"  <body>\n" +
			"    <h1>%s</h1>\n" +
			"  </body>\n" +
			"</html>\n",
		value,
	)
}
EOT

sudo cat <<EOT >> /etc/systemd/system/http-server.service
[Unit]
Description=HTTP Server
After=network.target
StartLimitIntervalSec=0

[Service]
Type=simple
Restart=always
RestartSec=1
User=root
ExecStart=/usr/bin/env go run /home/ec2-user/go/main.go

[Install]
WantedBy=multi-user.target
EOT

sudo systemctl start http-server
sudo systemctl enable http-server

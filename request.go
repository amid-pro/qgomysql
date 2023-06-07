package qgomysql

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"

	"golang.org/x/crypto/ssh"
)

func (config *Config) getData() {
	defer config.dataToJson()
	if config.raw != nil {
		return
	}
	config.query = "mysql -H -B -u" + config.MYSQL.User + " -p'" + config.MYSQL.Password + "' -e '" + config.query + "'"
	if config.SSH.Use_ssh {
		config.sshExec()
		return
	}
	config.requestExec()
}

func (config *Config) dataToJson() {
	out, err := json.Marshal(config.raw)
	config.raw = nil
	if err != nil {
		json_error, _ := json.Marshal(map[string]error{
			"Data error": err,
		})

		config.result = string(json_error)
		return
	}
	config.result = string(out)
}

func (config *Config) requestExec() {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", config.query)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	if err != nil {
		config.raw = append(config.raw, map_si{
			"Execute error": stderr.String(),
		})
		return
	}

	exit_code := cmd.ProcessState.ExitCode()
	if exit_code != 0 {
		config.raw = append(config.raw, map_si{
			"Exit code": exit_code,
		})
		return
	}

	config.result = stdout.String()
	config.tableToMap()
}

func (config *Config) sshExec() {
	ssh_config := &ssh.ClientConfig{
		User: config.SSH.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.SSH.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", config.SSH.Address+":"+fmt.Sprintf("%v", config.SSH.Port), ssh_config)
	if err != nil {
		config.raw = append(config.raw, map_si{
			"SSH dial": err.Error(),
		})
		return
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		config.raw = append(config.raw, map_si{
			"SSH client": err.Error(),
		})
		return
	}
	defer session.Close()

	var stdout bytes.Buffer
	session.Stdout = &stdout
	if err := session.Run(config.query); err != nil {
		config.raw = append(config.raw, map_si{
			"SSH run": err.Error(),
		})
		return
	}

	config.result = stdout.String()
	config.tableToMap()
}

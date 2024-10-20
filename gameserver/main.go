package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"path"
	"reflect"
	"syscall"
)

const SERVER_PATH = "/server"

var SERVER_CMD *exec.Cmd

type RustServerConfig struct {
	ServerHostname          string `json:"server.hostname"`
	ServerDescription       string `json:"server.description"`
	ServerUrl               string `json:"server.url"`
	ServerHeaderimage       string `json:"server.headerimage"`
	ServerTags              string `json:"server.tags"`
	ServerSeed              int    `json:"server.seed"`
	ServerWorldsize         int    `json:"server.worldsize"`
	ServerMaxplayers        int    `json:"server.maxplayers"`
	ServerTickrate          int    `json:"server.tickrate"`
	ServerSecure            int    `json:"server.secure"`
	ServerPort              int    `json:"server.port"`
	ServerQueryport         int    `json:"server.queryport"`
	ServerAppport           int    `json:"app.port"`
	ServerRconport          int    `json:"rcon.port"`
	ServerRconPassword      string `json:"rcon.password"`
	ServerAnalyticsUrl      string `json:"analytics.server_analytics_url"`
	ServerAnalyticsSecret   string `json:"analytics.analytics_secret"`
	ReportServerEndpoint    string `json:"server.reportsserverendpoint"`
	ReportServerEndpointKey string `json:"server.reportsserverendpointkey"`
	ServerStats             int    `json:"server.stats"`
}

func (h RustServerConfig) BuildArgs() []string {
	t := reflect.TypeOf(h)
	v := reflect.ValueOf(h)
	vals := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		jsontag := t.Field(i).Tag.Get("json")
		value := v.Field(i)
		vals = append(vals, fmt.Sprintf("+%s", jsontag))
		if value.CanInt() {
			vals = append(vals, fmt.Sprintf("%d", value.Int()))
		} else {
			vals = append(vals, value.String())
		}
		fmt.Printf("%s %s\n", jsontag, vals[len(vals)-1])
	}
	return vals
}

func (h RustServerConfig) BuildConfig() string {
	t := reflect.TypeOf(h)
	v := reflect.ValueOf(h)
	text := ""
	for i := 0; i < v.NumField(); i++ {
		key := t.Field(i).Tag.Get("json")
		field := v.Field(i)
		value := ""
		if field.CanInt() {
			value = fmt.Sprintf("\"%d\"", field.Int())
		} else {
			value = fmt.Sprintf("\"%s\"", field.String())
		}
		text = text + fmt.Sprintf("%s %s\n", key, value)
	}
	return text
}

func main() {
	var serverConfig RustServerConfig
	file, err := os.Open("./server_config.json")
	if err != nil {
		log.Fatalf("Failed to open server config file: %s", err.Error())
	}

	err = json.NewDecoder(file).Decode(&serverConfig)
	if err != nil {
		log.Fatalf("Failed to decode server config file: %s", err.Error())
	}
	err = CreateCfgFile(&serverConfig)
	if err != nil {
		log.Fatalf("Failed to create server config file: %s", err.Error())
	}

	go RunServer(&serverConfig)

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	SERVER_CMD.Process.Kill()
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CreateCfgFile(serverConfig *RustServerConfig) error {
	//"/serverauto.cfg"
	cfg_dir_path := path.Join(SERVER_PATH, "/server/my_server_identity/cfg")
	ex, err := exists(cfg_dir_path)
	if err != nil {
		return fmt.Errorf("failed to check exists: %v", err)
	}
	if !ex {
		err = os.MkdirAll(cfg_dir_path, 0777)
		if err != nil {
			return fmt.Errorf("failed to mkdirall: %v", err)
		}
	}

	os.Remove(path.Join(cfg_dir_path, "/serverauto.cfg"))

	file_content := serverConfig.BuildConfig()

	file, err := os.Create(path.Join(cfg_dir_path, "/serverauto.cfg"))
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	_, err = file.WriteString(file_content)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}
	err = file.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %v", err)
	}
	return nil
}

// RUN echo 'export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/server/RustDedicated_Data/Plugins' > /server/bootstrap.sh
// RUN echo 'export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/server/RustDedicated_Data/Plugins/x86_64' >> /server/bootstrap.sh
func RunServer(serverConfig *RustServerConfig) error {
	args := serverConfig.BuildArgs()

	SERVER_CMD = exec.Command(path.Join(SERVER_PATH, "/RustDedicated"), args...)
	SERVER_CMD.Stdout = os.Stdout
	SERVER_CMD.Stderr = os.Stderr
	SERVER_CMD.Env = append(SERVER_CMD.Env, "LD_LIBRARY_PATH=:/server/RustDedicated_Data/Plugins:/server/RustDedicated_Data/Plugins/x86_64")
	err := SERVER_CMD.Run()
	if err != nil {
		return err
	}
	return nil
}

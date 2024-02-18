package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"reflect"
	"syscall"
)

type RustServerConfig struct {
	ServerHostname        string `json:"server.hostname"`
	ServerDescription     string `json:"server.description"`
	ServerUrl             string `json:"server.url"`
	ServerHeaderimage     string `json:"server.headerimage"`
	ServerTags            string `json:"server.tags"`
	ServerSeed            int    `json:"server.seed"`
	ServerWorldsize       int    `json:"server.worldsize"`
	ServerMaxplayers      int    `json:"server.maxplayers"`
	ServerTickrate        int    `json:"server.tickrate"`
	ServerSecure          int    `json:"server.secure"`
	ServerPort            int    `json:"server.port"`
	ServerQueryport       int    `json:"server.queryport"`
	ServerAppport         int    `json:"app.port"`
	ServerRconport        int    `json:"rcon.port"`
	ServerRconPassword    string `json:"rcon.password"`
	ServerAnalyticsUrl    string `json:"server_analytics_url"`
	ServerAnalyticsSecret string `json:"analytics_secret"`
	ServerStats           int    `json:"server.stats"`
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

func main() {
	var serverConfig RustServerConfig
	file, err := os.Open("./server_config.json")
	if err != nil {
		log.Printf("Failed to open server config file: %s", err.Error())
		return
	}

	err = json.NewDecoder(file).Decode(&serverConfig)
	if err != nil {
		log.Printf("Failed to decode server config file: %s", err.Error())
		return
	}

	sigCh := make(chan os.Signal, 1)

	go RunServer(serverConfig, sigCh)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
}

// RUN echo 'export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/server/RustDedicated_Data/Plugins' > /server/bootstrap.sh
// RUN echo 'export LD_LIBRARY_PATH=$LD_LIBRARY_PATH:/server/RustDedicated_Data/Plugins/x86_64' >> /server/bootstrap.sh
func RunServer(serverConfig RustServerConfig, sig chan os.Signal) {
	args := serverConfig.BuildArgs()

	cmd := exec.Command("/server/RustDedicated", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, "LD_LIBRARY_PATH=:/server/RustDedicated_Data/Plugins:/server/RustDedicated_Data/Plugins/x86_64")
	err := cmd.Run()
	if err != nil {
		log.Printf("Failed to start server: %s", err.Error())
	}
	sig <- syscall.SIGQUIT
}

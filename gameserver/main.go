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
	"strings"
	"syscall"
)

const SERVER_PATH = "/server"

var SERVER_CMD *exec.Cmd

type RustServerConfig struct {
	ServerAnalyticsUrl                      string `json:"analytics.server_analytics_url" default:"http://metrics-server:5555/event"`
	ServerStats                             string `json:"server.stats" default:"1"`
	ServerAppport                           string `json:"app.port" default:"10012"`
	ServerRconport                          string `json:"rcon.port" default:"10013"`
	ServerRconip                            string `json:"rcon.ip" default:"0.0.0.0"`
	ServerRconWeb                           string `json:"rcon.web" default:"True"`
	ServerRconPassword                      string `json:"rcon.password" default:"sahisahdfbasb37"`
	Backtracklength                         string `json:"boombox.backtracklength" default:"30"`
	Serverurllist                           string `json:"boombox.serverurllist" default:""`
	Allowdesigning                          string `json:"ai.allowdesigning" default:"True"`
	Allusers                                string `json:"creative.allusers" default:"False"`
	Freebuild                               string `json:"creative.freebuild" default:"False"`
	Freeplacement                           string `json:"creative.freeplacement" default:"False"`
	Freerepair                              string `json:"creative.freerepair" default:"False"`
	Unlimitedio                             string `json:"creative.unlimitedio" default:"False"`
	Delete_after_upload                     string `json:"demo.delete_after_upload" default:"True"`
	Server_demo_cleanup_interval            string `json:"demo.server_demo_cleanup_interval" default:"20"`
	Server_demo_disk_space_gb               string `json:"demo.server_demo_disk_space_gb" default:"30"`
	Recordlist                              string `json:"demo.recordlist" default:""`
	Recordlistmode                          string `json:"demo.recordlistmode" default:"0"`
	Server_flush_seconds                    string `json:"demo.server_flush_seconds" default:"300"`
	Full_server_demo                        string `json:"demo.full_server_demo" default:"False"`
	Upload_demos                            string `json:"demo.upload_demos" default:"True"`
	Upload_url                              string `json:"demo.upload_url" default:""`
	Zip_demos                               string `json:"demo.zip_demos" default:"True"`
	Limit                                   string `json:"fps.limit" default:"240"`
	Maxspraysperplayer                      string `json:"global.maxspraysperplayer" default:"40"`
	Perf                                    string `json:"global.perf" default:"0"`
	Sprayduration                           string `json:"global.sprayduration" default:"10800"`
	Sprayoutofauthmultiplier                string `json:"global.sprayoutofauthmultiplier" default:"0.5"`
	Serversideragdolls                      string `json:"physics.serversideragdolls" default:"False"`
	Treecollision                           string `json:"physics.treecollision" default:"True"`
	Woundforever                            string `json:"player.woundforever" default:"False"`
	ServerHostname                          string `json:"server.hostname" default:"CRUST MAX3 X3 [SHIELDS|PROTECTION]"`
	ServerDescription                       string `json:"server.description" default:"Мод сервер с уникальной системой рейдов и защиты домов.\n\n• Макс команда: 3 человека\n• Рейты добычи: x3\n• Рейт добычи серы: x5\n• Скорость крафта: x2\n• Оффлайн укрепление построек\n\nВайпы каждые 2 недели, глобальный вайп раз в месяц\n"`
	ServerUrl                               string `json:"server.url" default:"https://discord.com/invite/UK6fJbmPp3"`
	ServerSeed                              string `json:"server.seed" default:"2174689"`
	ServerWorldsize                         string `json:"server.worldsize" default:"4250"`
	ServerMaxplayers                        string `json:"server.maxplayers" default:"250"`
	ServerTickrate                          string `json:"server.tickrate" default:"50"`
	ServerSecure                            string `json:"server.secure" default:"0"`
	ServerPort                              string `json:"server.port" default:"10010"`
	ServerQueryport                         string `json:"server.queryport" default:"10011"`
	Arrowarmor                              string `json:"server.arrowarmor" default:"1"`
	Arrowdamage                             string `json:"server.arrowdamage" default:"1"`
	Artificialtemperaturegrowablerange      string `json:"server.artificialtemperaturegrowablerange" default:"4"`
	Autouploadmap                           string `json:"server.autouploadmap" default:"True"`
	Bleedingarmor                           string `json:"server.bleedingarmor" default:"1"`
	Bleedingdamage                          string `json:"server.bleedingdamage" default:"1"`
	Bulletarmor                             string `json:"server.bulletarmor" default:"1"`
	Bulletdamage                            string `json:"server.bulletdamage" default:"1"`
	Canequipbackpacksinair                  string `json:"server.canequipbackpacksinair" default:"False"`
	Ceilinglightgrowablerange               string `json:"server.ceilinglightgrowablerange" default:"3"`
	Ceilinglightheightoffset                string `json:"server.ceilinglightheightoffset" default:"3"`
	Conveyormovefrequency                   string `json:"server.conveyormovefrequency" default:"5"`
	Crawlingenabled                         string `json:"server.crawlingenabled" default:"True"`
	Defaultblueprintresearchcost            string `json:"server.defaultblueprintresearchcost" default:"10"`
	Enforcepipechecksonbuildingblockchanges string `json:"server.enforcepipechecksonbuildingblockchanges" default:"True"`
	Favoritesendpoint                       string `json:"server.favoritesendpoint" default:""`
	Funwaterdamagethreshold                 string `json:"server.funwaterdamagethreshold" default:"0.8"`
	Funwaterwetnessgain                     string `json:"server.funwaterwetnessgain" default:"0.05"`
	Headerimage                             string `json:"server.headerimage" default:"https://i.imgur.com/hXUzeof.png"`
	Incapacitatedrecoverchance              string `json:"server.incapacitatedrecoverchance" default:"0.1"`
	Industrialcrafterfrequency              string `json:"server.industrialcrafterfrequency" default:"5"`
	Industrialframebudgetms                 string `json:"server.industrialframebudgetms" default:"0.5"`
	Industrialtransferstricttimelimits      string `json:"server.industrialtransferstricttimelimits" default:"False"`
	Logoimage                               string `json:"server.logoimage" default:""`
	Maximummapmarkers                       string `json:"server.maximummapmarkers" default:"5"`
	Maximumpings                            string `json:"server.maximumpings" default:"5"`
	Maxitemstacksmovedpertickindustrial     string `json:"server.maxitemstacksmovedpertickindustrial" default:"12"`
	Meleearmor                              string `json:"server.meleearmor" default:"1"`
	Meleedamage                             string `json:"server.meleedamage" default:"1"`
	Motd                                    string `json:"server.motd" default:""`
	Nonplanterdeathchancepertick            string `json:"server.nonplanterdeathchancepertick" default:"0.005"`
	Oilrig_radiation_amount_scale           string `json:"server.oilrig_radiation_amount_scale" default:"1"`
	Oilrig_radiation_time_scale             string `json:"server.oilrig_radiation_time_scale" default:"1"`
	Optimalplanterqualitysaturation         string `json:"server.optimalplanterqualitysaturation" default:"0.6"`
	Parachuterepacktime                     string `json:"server.parachuterepacktime" default:"8"`
	Pingduration                            string `json:"server.pingduration" default:"10"`
	Playerserverfall                        string `json:"server.playerserverfall" default:"True"`
	Printreportstoconsole                   string `json:"server.printreportstoconsole" default:"False"`
	Reportsserverendpoint                   string `json:"server.reportsserverendpoint" default:"http://metrics-server:5555/feedback"`
	Reportsserverendpointkey                string `json:"server.reportsserverendpointkey" default:"aboba"`
	Rewounddelay                            string `json:"server.rewounddelay" default:"60"`
	Savebackupcount                         string `json:"server.savebackupcount" default:"2"`
	Server_id                               string `json:"server.server_id" default:"cbfd2d65e7d34b129a28d76aafc54160"`
	Showholstereditems                      string `json:"server.showholstereditems" default:"True"`
	Sprinklereyeheightoffset                string `json:"server.sprinklereyeheightoffset" default:"3"`
	Sprinklerradius                         string `json:"server.sprinklerradius" default:"3"`
	Tags                                    string `json:"server.tags" default:"biweekly,EU"`
	Tutorialenabled                         string `json:"server.tutorialenabled" default:"False"`
	Watercontainersleavewaterbehind         string `json:"server.watercontainersleavewaterbehind" default:"False"`
	Workbench1taxrate                       string `json:"server.workbench1taxrate" default:"0"`
	Workbench2taxrate                       string `json:"server.workbench2taxrate" default:"10"`
	Workbench3taxrate                       string `json:"server.workbench3taxrate" default:"20"`
	Woundedmaxfoodandwaterbonus             string `json:"server.woundedmaxfoodandwaterbonus" default:"0.25"`
	Woundedrecoverchance                    string `json:"server.woundedrecoverchance" default:"0.2"`
	Woundingenabled                         string `json:"server.woundingenabled" default:"True"`
	Server_allow_steam_nicknames            string `json:"steam.server_allow_steam_nicknames" default:"True"`
	Analytics_header                        string `json:"analytics.analytics_header" default:"X-API-KEY"`
	Analytics_secret                        string `json:"analytics.analytics_secret" default:"aboba"`
	Analytics_bulk_upload_url               string `json:"analytics.analytics_bulk_upload_url" default:""`
	Gameplay_analytics                      string `json:"analytics.gameplay_analytics" default:"False"`
	Gameplay_tick_analytics                 string `json:"analytics.gameplay_tick_analytics" default:"False"`
	Performance_analytics                   string `json:"analytics.performance_analytics" default:"True"`
	Runtime_profiling                       string `json:"profile.runtime_profiling" default:"0"`
	Runtime_profiling_persist               string `json:"profile.runtime_profiling_persist" default:"False"`
	Dynamicpricingenabled                   string `json:"npcvendingmachine.dynamicpricingenabled" default:"True"`
	Maximumpricemultiplier                  string `json:"npcvendingmachine.maximumpricemultiplier" default:"2"`
	Minimumpricemultiplier                  string `json:"npcvendingmachine.minimumpricemultiplier" default:"0.5"`
	Pricedecreaseamount                     string `json:"npcvendingmachine.pricedecreaseamount" default:"0.05"`
	Priceincreaseamount                     string `json:"npcvendingmachine.priceincreaseamount" default:"0.1"`
	Priceupdatefrequencybiweekly            string `json:"npcvendingmachine.priceupdatefrequencybiweekly" default:"2"`
	Priceupdatefrequencydefault             string `json:"npcvendingmachine.priceupdatefrequencydefault" default:"3"`
	Priceupdatefrequencyweekly              string `json:"npcvendingmachine.priceupdatefrequencyweekly" default:"1"`
	Startingpricemultiplier                 string `json:"npcvendingmachine.startingpricemultiplier" default:"2"`
	Bypassrepack                            string `json:"parachute.bypassrepack" default:"False"`
	Landinganimations                       string `json:"parachute.landinganimations" default:"False"`
	Max_speed                               string `json:"travellingvendor.max_speed" default:"5"`
	Enforcetrespasschecks                   string `json:"tutorialisland.enforcetrespasschecks" default:"True"`
	Spawntutorialislandfornewplayer         string `json:"tutorialisland.spawntutorialislandfornewplayer" default:"True"`
	Racetimeout                             string `json:"waypointrace.racetimeout" default:"900"`
}

func (h RustServerConfig) BuildArgs() []string {
	args := make([]string, 0)

	args = append(args, "+rcon.port", h.ServerRconport)
	args = append(args, "+rcon.web", h.ServerRconWeb)
	args = append(args, "+rcon.ip", h.ServerRconip)
	args = append(args, "+rcon.password", h.ServerRconPassword)

	return args
}
func (h RustServerConfig) BuildConfig() string {
	t := reflect.TypeOf(h)
	v := reflect.ValueOf(h)
	text := ""
	for i := 0; i < v.NumField(); i++ {
		key := t.Field(i).Tag.Get("json")
		field := v.Field(i)
		defaultValue := t.Field(i).Tag.Get("default")

		defaultValue = strings.ReplaceAll(defaultValue, "\n", "\\n")

		value := ""

		if len(field.String()) == 0 {
			value = fmt.Sprintf("\"%s\"", defaultValue)
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
		f, err := os.Create("./server_config.json")
		if err == nil {
			b, _ := json.MarshalIndent(&RustServerConfig{}, "", "\t")
			f.Write(b)
			f.Close()
		}
		log.Fatalf("Failed to open server config file, creating: %s", err.Error())
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
	fmt.Println(file_content)
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
	SERVER_CMD = exec.Command(path.Join(SERVER_PATH, "/RustDedicated"), serverConfig.BuildArgs()...)
	SERVER_CMD.Stdout = os.Stdout
	SERVER_CMD.Stderr = os.Stderr
	SERVER_CMD.Env = append(SERVER_CMD.Env, "LD_LIBRARY_PATH=:/server/RustDedicated_Data/Plugins:/server/RustDedicated_Data/Plugins/x86_64")
	err := SERVER_CMD.Run()
	if err != nil {
		return err
	}
	return nil
}

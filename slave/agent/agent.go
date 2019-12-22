package agent

import (
	"bytes"
	"encoding/json"
	"github.com/inconshreveable/go-update"
	"github.com/r3labs/sse"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/response"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"
)

var (
	metricsFrequency int
)

func Start(hostToken string) {
	logrus.Info("Starting Agent")
	config.HOST_TOKEN = hostToken

	resp, err := http.Get(config.API_URL + "/v1/host/auth/" + config.HOST_TOKEN)
	if err != nil {
		lib.FailOnError(err, "Failed to get host token")
	}

	var token response.Token
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		lib.FailOnError(err, "Failed to decode credentials json")
	}
	config.API_TOKEN = token.Token

	go agent()

	logrus.Info("Send AGENT_METRICS_FREQUENCY")
	taskMsg := task.Message{
		TaskId: task.AGENT_METRICS_FREQUENCY,
	}

	jsonStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Error(err)
	}
	//TODO: do something with status code
	_, err = lib.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, jsonStr)
	if err != nil {
		logrus.Error(err)
	}

	sendOSInfo()
	go heartbeat()

	select {}
}

func agent() {
	client := sse.NewClient(config.API_URL + "/v1/events/" + config.HOST_TOKEN)
	client.Headers = map[string]string{
		"Authorization": "Bearer " + config.API_TOKEN,
	}
	err := client.SubscribeRaw(func(msg *sse.Event) {
		logrus.Info(msg.ID)
		logrus.Info(string(msg.Data))
		logrus.Info(string(msg.Event))

		taskMsg := task.Message{}
		err := taskMsg.Deserialize(string(msg.Data))
		if err != nil {
			logrus.Error(err)
		}
		taskId, _ := strconv.Atoi(string(msg.Event))

		switch taskId {
		case task.WRAPPER_START:
			logrus.Info("START_WRAPPER")
			//TODO: move gsID to env variable
			cmd := exec.Command("slave", "wrapper", taskMsg.GameServerID)

			cmd.Env = os.Environ()
			cmd.Env = append(cmd.Env, "HOST_TOKEN="+config.HOST_TOKEN)
			cmd.Env = append(cmd.Env, "API_TOKEN="+config.API_TOKEN)

			//TODO: FOR TESTING ONLY
			var stdBuffer bytes.Buffer
			mw := io.MultiWriter(os.Stdout, &stdBuffer)
			cmd.Stdout = mw
			cmd.Stderr = mw

			if err := cmd.Start(); err != nil {
				logrus.Error(err)
			}

			cmd.Process.Release()
		case task.GAME_INSTALL:
			logrus.Info("INSTALL_GAME_SERVER")

			gsPath := filepath.Clean(config.GS_DIR) + "/" + taskMsg.GameServerID
			if _, err := os.Stat(gsPath); os.IsNotExist(err) {
				//TODO: Set correct perms
				_ = os.Mkdir(gsPath, 0777)
			}

			err = taskMsg.Game.Install(filepath.Clean(config.STORAGE_DIR), gsPath)
			if err != nil {
				logrus.Error(err)
			}

			logrus.Info("Installation finished")
		case task.AGENT_METRICS_FREQUENCY:
			logrus.Info("AGENT_METRICS_FREQUENCY")

			metricsFrequency = taskMsg.MetricFrequency

			go metrics()
		case task.AGENT_REMOVE_GS:
			logrus.Info("AGENT_REMOVE_GS")

			gsId := taskMsg.GameServerID
			gsPath := filepath.Clean(config.GS_DIR) + "/" + gsId
			gsTrashPath := filepath.Clean(config.TRASH_DIR) + "/" + gsId

			err := os.Rename(gsPath, gsTrashPath)
			if err != nil {
				logrus.Error(err)
			}
		case task.SLAVE_UPDATE:
			if config.GIT_COMMIT == taskMsg.Commit {
				return
			}

			resp, err := http.Get(taskMsg.Url)
			if err != nil {
				logrus.Error(err)
			}
			defer resp.Body.Close()

			err = update.Apply(resp.Body, update.Options{})
			if err != nil {
				logrus.Error(err)
			}
			os.Exit(0)
		}
	})
	if err != nil {
		lib.FailOnError(err, "Failed to subscribe a channel")
	}
}

func metrics() {
	for {
		<-time.After(time.Duration(metricsFrequency) * time.Second)

		virtualMemory, _ := mem.VirtualMemory()
		ramFree := virtualMemory.Free / 1024 / 1024
		ramTotal := virtualMemory.Total / 1024 / 1024
		cpuPercent, _ := cpu.Percent(time.Duration(1)*time.Second, false)
		cpuTimes, _ := cpu.Times(false)

		diskUsage, _ := disk.Usage("/")
		diskFree := diskUsage.Free / 1024 / 1024
		diskTotal := diskUsage.Total / 1024 / 1024
		diskUsed := diskUsage.Used / 1024 / 1024

		taskMsg := task.Message{
			TaskId:    task.AGENT_METRICS,
			CpuUsage:  int(cpuPercent[0]),
			RamFree:   int(ramFree),
			RamTotal:  int(ramTotal),
			DiskFree:  int(diskFree),
			DiskTotal: int(diskTotal),
			DiskUsed:  int(diskUsed),
			User:      int(cpuTimes[0].User),
			System:    int(cpuTimes[0].System),
			Idle:      int(cpuTimes[0].Idle),
			Nice:      int(cpuTimes[0].Nice),
			Iowait:    int(cpuTimes[0].Iowait),
			Irq:       int(cpuTimes[0].Irq),
			Softirq:   int(cpuTimes[0].Softirq),
			Steal:     int(cpuTimes[0].Steal),
			Guest:     int(cpuTimes[0].Guest),
			GuestNice: int(cpuTimes[0].GuestNice),
		}

		jsonStr, err := taskMsg.Serialize()
		if err != nil {
			logrus.Error(err)
		}
		//TODO: do something with status code
		_, err = lib.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, jsonStr)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func heartbeat() {
	for {
		<-time.After(5 * time.Second)

		logrus.Info("Sending heartbeat")

		taskMsg := task.Message{
			TaskId:    task.AGENT_HEARTBEAT,
			Timestamp: time.Now().Unix(),
		}

		jsonStr, err := taskMsg.Serialize()
		if err != nil {
			logrus.Error(err)
		}
		//TODO: do something with status code
		_, err = lib.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, jsonStr)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func sendOSInfo() {
	h, _ := host.Info()

	taskMsg := task.Message{
		TaskId: task.AGENT_OS,

		OS:              h.OS,
		Platform:        h.Platform,
		PlatformFamily:  h.PlatformFamily,
		PlatformVersion: h.PlatformVersion,
		KernelVersion:   h.KernelVersion,
		KernelArch:      h.KernelArch,
	}

	jsonStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Error(err)
	}
	//TODO: do something with status code
	_, err = lib.SendTaskData(config.API_URL+"/v1/events/"+config.HOST_TOKEN, config.API_TOKEN, jsonStr)
	if err != nil {
		logrus.Error(err)
	}
}

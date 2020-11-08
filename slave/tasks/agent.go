package tasks

import (
	"context"
	"github.com/inconshreveable/go-update"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib/ahttp"
	"gitlab.com/systemz/aimpanel2/lib/filemanager"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func AgentTaskHandler(taskMsg task.Message) {
	switch taskMsg.TaskId {
	case task.AGENT_START_GS:
		// kill all instances of GS before starting it, just in case
		// FIXME detect current state before starting server
		GsKill(taskMsg.GameServerID)
		// TODO check why we need this
		model.SetGsGame(taskMsg.GameServerID, taskMsg.Game)
		model.SetGsStart(taskMsg.GameServerID, 1)
		// start game server by running wrapper
		StartWrapperInDocker(taskMsg)
		if false {
			StartWrapperExecRaw(taskMsg)
		}
	case task.AGENT_INSTALL_GS:
		GsInstall(taskMsg)
	case task.AGENT_UPDATE:
		SelfUpdate(taskMsg)
	case task.AGENT_REMOVE_GS:
		GsRemove(taskMsg)
	case task.AGENT_FILE_LIST_GS:
		GsFileList(taskMsg.GameServerID)
	case task.AGENT_METRICS_FREQUENCY:
		go AgentMetrics(taskMsg.MetricFrequency)
	case task.AGENT_BACKUP_GS:
		GsBackup(taskMsg.GameServerID)
	case task.AGENT_BACKUP_RESTORE_GS:
		GsBackupRestore(taskMsg.GameServerID, taskMsg.BackupFilename)
	case task.AGENT_GET_UPDATE:
		go AgentGetUpdate(taskMsg)
	case task.AGENT_BACKUP_LIST_GS:
		go AgentSendGsBackupList(taskMsg.GameServerID)
	case task.AGENT_CLEAN_REINSTALL_GS:
		go GsCleanReinstall(taskMsg)
	}
}

func GsRemove(taskMsg task.Message) {
	gsId := taskMsg.GameServerID
	gsPath := filepath.Clean(config.GS_DIR) + "/" + gsId
	gsTrashPath := filepath.Clean(config.TRASH_DIR) + "/" + gsId

	err := os.Rename(gsPath, gsTrashPath)
	if err != nil {
		logrus.Error(err)
	}
}

func GsInstall(taskMsg task.Message) {
	gsPath := filepath.Clean(config.GS_DIR) + "/" + taskMsg.GameServerID
	if _, err := os.Stat(gsPath); os.IsNotExist(err) {
		//TODO: Set correct perms
		_ = os.Mkdir(gsPath, 0777)
	}

	err := taskMsg.Game.Install(filepath.Clean(config.STORAGE_DIR), gsPath)
	if err != nil {
		logrus.Error(err)
	}
}

func SelfUpdate(taskMsg task.Message) {
	if config.GIT_COMMIT == "" {
		logrus.Warning("version of slave is empty, ignoring update")
		return
	}

	if config.GIT_COMMIT == taskMsg.Commit {
		logrus.Warning("version of new slave is same as current, ignoring update")
		return
	}
	if taskMsg.Commit == "" {
		logrus.Warning("version of new slave is empty, ignoring update")
		return
	}

	logrus.Infof("downloading new version %v from %v", taskMsg.Commit, taskMsg.Url)
	resp, err := http.Get(taskMsg.Url)
	if err != nil {
		logrus.Error(err)
	}
	defer resp.Body.Close()

	err = update.Apply(resp.Body, update.Options{
		TargetPath:  "",
		TargetMode:  0,
		Checksum:    nil,
		PublicKey:   nil,
		Signature:   nil,
		Verifier:    nil,
		Hash:        0,
		Patcher:     nil,
		OldSavePath: "",
	})
	if err != nil {
		logrus.Error(err)
	}
	logrus.Info("shutting down agent to apply update")
	os.Exit(0)
}

func SelfHeal() {
	//Check if directory exist - if not create it
	if _, err := os.Stat(config.STORAGE_DIR); err != nil {
		if err := os.MkdirAll(config.STORAGE_DIR, 0755); err != nil {
			logrus.Warnf("Could not create %s directory", config.STORAGE_DIR)
		}
	}

	if _, err := os.Stat(config.GS_DIR); err != nil {
		if err := os.MkdirAll(config.GS_DIR, 0755); err != nil {
			logrus.Warnf("Could not create %s directory", config.GS_DIR)
		}
	}

	if _, err := os.Stat(config.BACKUP_DIR); err != nil {
		if err := os.MkdirAll(config.BACKUP_DIR, 0755); err != nil {
			logrus.Warnf("Could not create %s directory", config.BACKUP_DIR)
		}
	}

	if _, err := os.Stat(config.TRASH_DIR); err != nil {
		if err := os.MkdirAll(config.TRASH_DIR, 0755); err != nil {
			logrus.Warnf("Could not create %s directory", config.TRASH_DIR)
		}
	}
}

func GsFileList(gsId string) {
	logrus.Infof("File list for GS ID %v started", gsId)

	node, err := filemanager.NewTree(config.GS_DIR+"/"+gsId, 100, 64)
	if err != nil {
		logrus.Error(err)
	}

	taskMsg := task.Message{
		TaskId:       task.AGENT_FILE_LIST_GS,
		GameServerID: gsId,
		Files:        node,
	}

	_, err = ahttp.SendTaskData("/v1/events/"+config.HOST_TOKEN, config.HW_ID, taskMsg)
	if err != nil {
		logrus.Error(err)
	}

	logrus.Infof("File list for GS ID %v finished", gsId)
}

func AgentShutdown() {
	logrus.Info("Send AGENT_SHUTDOWN")
	taskMsg := task.Message{
		TaskId: task.AGENT_SHUTDOWN,
	}

	//TODO: do something with status code
	_, err := ahttp.SendTaskData("/v1/events/"+config.HOST_TOKEN, config.HW_ID, taskMsg)
	if err != nil {
		logrus.Error(err)
	}

	//GsStop("all")

	<-time.After(5 * time.Second)

	os.Exit(1)
}

func AgentShutdownTrigger() {
	taskMsg := task.Message{
		// FIXME other task IDs for user CLI actions
		TaskId: task.AGENT_SHUTDOWN,
	}
	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Errorf("preparing msg failed: %v", err)
		return
	}

	res, err := model.Redis.Publish(context.TODO(), config.REDIS_PUB_SUB_AGENT_CH, taskMsgStr).Result()
	if err != nil {
		logrus.Errorf("sending msg failed: %v", err)
	}
	logrus.Infof("Task sent to %v processes", res)
}

func AgentMetrics(metricsFrequency int) {
	for {
		<-time.After(time.Duration(metricsFrequency) * time.Second)

		virtualMemory, _ := mem.VirtualMemory()
		ramFree := virtualMemory.Free / 1024 / 1024
		ramTotal := virtualMemory.Total / 1024 / 1024
		ramCache := virtualMemory.Cached / 1024 / 1024
		ramBuffers := virtualMemory.Buffers / 1024 / 1024
		ramAvailable := virtualMemory.Available / 1024 / 1024
		cpuPercent, _ := cpu.Percent(time.Duration(1)*time.Second, false)
		cpuTimes, _ := cpu.Times(false)

		diskUsage, _ := disk.Usage("/")
		diskFree := diskUsage.Free / 1024 / 1024
		diskTotal := diskUsage.Total / 1024 / 1024
		diskUsed := diskUsage.Used / 1024 / 1024

		taskMsg := task.Message{
			TaskId:       task.AGENT_METRICS,
			CpuUsage:     int(cpuPercent[0]),
			RamFree:      int(ramFree),
			RamCache:     int(ramCache),
			RamBuffers:   int(ramBuffers),
			RamTotal:     int(ramTotal),
			RamAvailable: int(ramAvailable),
			DiskFree:     int(diskFree),
			DiskTotal:    int(diskTotal),
			DiskUsed:     int(diskUsed),
			User:         int(cpuTimes[0].User),
			System:       int(cpuTimes[0].System),
			Idle:         int(cpuTimes[0].Idle),
			Nice:         int(cpuTimes[0].Nice),
			Iowait:       int(cpuTimes[0].Iowait),
			Irq:          int(cpuTimes[0].Irq),
			Softirq:      int(cpuTimes[0].Softirq),
			Steal:        int(cpuTimes[0].Steal),
			Guest:        int(cpuTimes[0].Guest),
			GuestNice:    int(cpuTimes[0].GuestNice),
		}
		//TODO: do something with status code
		_, err := ahttp.SendTaskData("/v1/events/"+config.HOST_TOKEN, config.HW_ID, taskMsg)
		if err != nil {
			logrus.Error(err)
		}
	}
}

func AgentSendOSInfo() {
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
	//TODO: do something with status code
	_, err := ahttp.SendTaskData("/v1/events/"+config.HOST_TOKEN, config.HW_ID, taskMsg)
	if err != nil {
		logrus.Error(err)
	}
}

func AgentGetUpdate(taskMsg task.Message) {
	_, err := ahttp.SendTaskData("/v1/events/"+config.HOST_TOKEN, config.HW_ID, taskMsg)
	if err != nil {
		logrus.Error(err)
	}
}

//Backups
func GsBackupTrigger(gsId string) {
	taskMsg := task.Message{
		// FIXME other task IDs for user CLI actions
		TaskId:       task.AGENT_BACKUP_GS,
		GameServerID: gsId,
	}
	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Errorf("preparing msg failed: %v", err)
		return
	}
	res, err := model.Redis.Publish(context.TODO(), config.REDIS_PUB_SUB_AGENT_CH, taskMsgStr).Result()
	if err != nil {
		logrus.Errorf("sending msg failed: %v", err)
	}
	logrus.Infof("Task sent to %v processes", res)
}

func GsBackup(gsId string) {
	logrus.Infof("Backup for GS ID %v started", gsId)

	// prepare destination name and path for backup
	unixTimestamp := strconv.Itoa(int(time.Now().Unix()))
	// FIXME add human readable UTC date at the end
	backupFilename := unixTimestamp + "_" + gsId + ".tar.gz"
	backupPath := config.BACKUP_DIR + backupFilename
	inputDirPath := strings.TrimRight(config.GS_DIR+gsId, "/")

	// create backup
	TarGz(backupPath, inputDirPath, true)

	// all done!
	logrus.Infof("Backup for GS ID %v finished", gsId)
}


func GsBackupRestoreTrigger(gsId string, backupFilename string) {
	taskMsg := task.Message{
		// FIXME other task IDs for user CLI actions
		TaskId:         task.AGENT_BACKUP_RESTORE_GS,
		GameServerID:   gsId,
		BackupFilename: backupFilename,
	}
	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Errorf("preparing msg failed: %v", err)
		return
	}
	res, err := model.Redis.Publish(context.TODO(), config.REDIS_PUB_SUB_AGENT_CH, taskMsgStr).Result()
	if err != nil {
		logrus.Errorf("sending msg failed: %v", err)
	}
	logrus.Infof("Task sent to %v processes", res)
}

func GsBackupRestore(gsId string, backupFilename string) {
	logrus.Infof("Backup restore for GS ID %v started", gsId)

	GsCleanFiles(gsId)

	//extract backup to gs dir
	UnTar(filepath.Join(config.BACKUP_DIR, backupFilename), filepath.Join(config.GS_DIR, gsId))

	logrus.Infof("Backup restore for GS ID %v finished", gsId)
}

func GsCleanFiles(gsId string) {
	logrus.Infof("Cleaning files for GS ID %v started", gsId)

	//remove all files in gs dir
	files, err := filepath.Glob(filepath.Join(config.GS_DIR, gsId, "*"))
	if err != nil {
		logrus.Error(err)
	}

	for _, file := range files {
		err = os.RemoveAll(file)
		if err != nil {
			logrus.Error(err)
		}
	}

	logrus.Infof("Cleaning files for GS ID %v finished", gsId)
}

func AgentSendGsBackupList(gsId string) {
	//Get file list which ends with "_{gsId}.tar.gz"
	var files []string
	filepath.Walk(config.BACKUP_DIR, func(path string, fi os.FileInfo, err error) error {
		logrus.Info(fi.Name())
		if !fi.IsDir() {
			if strings.HasSuffix(fi.Name(), "_"+gsId+".tar.gz") {
				files = append(files, fi.Name())
			}
		}
		return nil
	})

	taskMsg := task.Message{
		TaskId:       task.AGENT_BACKUP_LIST_GS,
		GameServerID: gsId,
		Backups:      files,
	}
	//TODO: do something with status code
	_, err := ahttp.SendTaskData("/v1/events/"+config.HOST_TOKEN, config.HW_ID, taskMsg)
	if err != nil {
		logrus.Error(err)
	}
}

func GsCleanReinstall(taskMsg task.Message) {
	GsCleanFiles(taskMsg.GameServerID)
	GsInstall(taskMsg)
}

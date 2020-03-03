package tasks

import (
	"github.com/inconshreveable/go-update"
	"github.com/sirupsen/logrus"
	"gitlab.com/systemz/aimpanel2/lib"
	"gitlab.com/systemz/aimpanel2/lib/task"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"
)

// tasks below will be eventually finished by agent
func StartWrapper(taskMsg task.Message) {
	cred := &syscall.Credential{
		Uid:         uint32(syscall.Getuid()),
		Gid:         uint32(syscall.Getgid()),
		Groups:      []uint32{},
		NoSetGroups: true,
	}

	sysproc := &syscall.SysProcAttr{
		Credential: cred,
		Noctty:     true,
		Setpgid:    true,
	}

	attr := os.ProcAttr{
		Dir: ".",
		Env: os.Environ(),
		Files: []*os.File{
			os.Stdin,
			nil,
			nil,
		},
		Sys: sysproc,
	}
	attr.Env = append(attr.Env, "HOST_TOKEN="+config.HOST_TOKEN)
	attr.Env = append(attr.Env, "API_TOKEN="+config.API_TOKEN)

	//TODO: move gsID to env variable
	process, err := os.StartProcess("/usr/local/bin/slave", []string{"/usr/local/bin/slave", "wrapper", taskMsg.GameServerID}, &attr)
	if err != nil {
		logrus.Error(err)
	}

	go func() {
		_, err := process.Wait()
		if err != nil {
			logrus.Error(err)
		}
	}()
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
	if config.GIT_COMMIT == taskMsg.Commit {
		return
	}
	if config.GIT_COMMIT == "" {
		return
	}

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

func GsBackupTrigger(gsId string) {
	taskMsg := task.Message{
		// FIXME other task IDs for user CLI actions
		TaskId:       task.GAME_MAKE_BACKUP,
		GameServerID: gsId,
	}
	taskMsgStr, err := taskMsg.Serialize()
	if err != nil {
		logrus.Errorf("preparing msg failed: %v", err)
		return
	}
	res, err := model.Redis.Publish(config.REDIS_PUB_SUB_CH, taskMsgStr).Result()
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

func AgentShutdown() {
	logrus.Info("Send AGENT_SHUTDOWN")
	taskMsg := task.Message{
		TaskId: task.AGENT_SHUTDOWN,
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

	res, err := model.Redis.Publish(config.REDIS_PUB_SUB_CH, taskMsgStr).Result()
	if err != nil {
		logrus.Errorf("sending msg failed: %v", err)
	}
	logrus.Infof("Task sent to %v processes", res)
}

package cmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gitlab.com/systemz/aimpanel2/slave/config"
	"gitlab.com/systemz/aimpanel2/slave/model"
	"gitlab.com/systemz/aimpanel2/slave/tasks/gdrive"
	"os"
)

func init() {
	rootCmd.AddCommand(backupUploadCmd)
}

var backupUploadCmd = &cobra.Command{
	Use:   "backup-upload <BACKUP FILENAME> <BACKUP FILENAME> ...",
	Short: "Upload backup",
	Long:  "Upload game server backup to Google Drive",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		model.InitRedis()
		for _, backupFilename := range args {
			filePath := config.BACKUP_DIR + backupFilename
			logrus.Infof("Uploading to drive: %v", filePath)

			f, err := os.Open(filePath)
			if err != nil {
				logrus.Panicf("File can't be open/found: %v", err)
			}

			// FIXME a lot of backups to upload can cause a lot of trouble not closing this sooner
			defer f.Close()

			service := gdrive.ClientInit()
			// FIXME get md5sum from upload to store in DB
			_, err = gdrive.UploadFile(service, f, backupFilename)
			if err != nil {
				logrus.Panicf("File can't uploaded: %v", err)
			}
			logrus.Info("Backup uploaded successfully")
		}
	},
}

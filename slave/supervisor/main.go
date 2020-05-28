package supervisor

func Start() {
	InstallPackages()
	ReportInit()
	WatchSshd()
}

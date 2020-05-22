package supervisor

func Start() {
	ReportInit()
	WatchSshd()
}

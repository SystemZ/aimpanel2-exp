package supervisor

//var (
//	reportRateLimit ratelimit.Limit
//)

func ReportInit() {
	//reportRateLimit = ratelimit.CreateLimit("1r/h")
}

func ReportIp(ip string, logMsg string) {
	//err := reportRateLimit.Hit(ip)
	// skip if recently reported
	//if err != nil {
	//	return
	//logrus.Info("reported " + ip)
}

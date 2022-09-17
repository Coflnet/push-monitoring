package usecase

func StartMonitoringServices() {
	go monitorModSocket()
}

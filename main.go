package main

func main() {
	settings := loadSettings()
	dataReady := make(chan struct{})

	go updateLaunchesPeriodically(settings.RefreshInterval, dataReady)

	// Wait here for the first fetch to complete before starting the HTTP server
	<-dataReady

	// Start the HTTP server to serve the launch data
	serve(settings.HttpPort, settings.RefreshInterval, settings.LogRequests)
}

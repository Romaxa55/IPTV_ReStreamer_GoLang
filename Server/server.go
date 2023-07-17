package Server

import (
	"IPTV_ReStreamer_GoLang/Config"
	"IPTV_ReStreamer_GoLang/Logger"
	"fmt"
	"net/http"
)

func StartServer() {
	config := Config.GetServerConfig()
	Logger.Log.Info(fmt.Sprintf("Starting server at http://%s:%s", config.IP, config.Port))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/startStream", startStreamHandler)
	http.HandleFunc("/videoFiles", videoFilesHandler)
	http.HandleFunc("/status_health", statusHealthHandler)
	// Добавьте другие обработчики здесь

	err := http.ListenAndServe(config.IP+":"+config.Port, nil)
	if err != nil {
		Logger.Log.Error("Server failed to start: ", err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the Home Page!")
}

func startStreamHandler(w http.ResponseWriter, r *http.Request) {
	// Ваш код здесь
}

func videoFilesHandler(w http.ResponseWriter, r *http.Request) {
	// Ваш код здесь
}

func statusHealthHandler(w http.ResponseWriter, r *http.Request) {
	// Ваш код здесь
}

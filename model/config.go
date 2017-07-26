package model

type FimpUiConfigs struct {
	ReportLogFiles     []string `json:"report_log_files"`
	ReportLogSizeLimit int64    `json:"report_log_size_limit"`
	VinculumAddress    string   `json:"vinculum_address"`
	MqttServerURI      string   `json:"mqtt_server_uri"`
	FlowStorageDir     string 	`json:"flow_storage_dir"`
	MqttClientIdPrefix string   `json:"mqtt_client_id_prefix"`
	LogFile            string   `json:"log_file"`
	LogLevel           string   `json:"log_level"`
}

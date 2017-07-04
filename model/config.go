package model

type FimpUiConfigs struct {
	ReportLogFiles []string `json:"report_log_files"`
	ReportLogSizeLimit int64 `json:"report_log_size_limit"`
	VinculumAddress string `json:"vinculum_address"`
}
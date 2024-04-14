package util

type APIError struct {
	Status   string `json:"status"`
	Code     int    `json:"code"`
	Internal string `json:"internal"`
	Message  string `json:"message"`
}

type Log struct {
	Message string `json:"message"`
	Time    string `json:"time"`
}

type config struct {
	Credentials       map[string]string `json:"credentials"`
	BucketName        string            `json:"bucketName"`
	LogFileDir        string            `json:"logFileDir"`
	TimeLayout        string            `json:"timeLayout"`
	Port              string            `json:"port"`
	WriteLogInterval  int               `json:"writeLogInterval"`  // minutes
	UploadLogInterval int               `json:"uploadLogInterval"` // minutes
	ServerStorage     bool              `json:"serverStorage"`
}

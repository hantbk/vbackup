package global

const (
	DownloadUrl     = "https://github.com/hantbk/vbackup/releases/download"                                   // Download URL for installation packages
	ServiceFileUrl  = "https://raw.githubusercontent.com/hantbk/vbackup/master/vbackup.service"               // URL for service file download
	LatestUrl       = "https://api.github.com/repos/hantbk/vbackup/releases?page=1&per_page=1&direction=desc" // URL to get the latest version
	OTPSecretLength = 16                                                                                      // Secret length
	OTPDigits       = 6                                                                                       // Number of OTP digits
)

var EmailChannel = make(chan string, 1) // Buffer of 1 to avoid blocking

package guard

type CaptchaDTO struct {
	Id          string
	ImageBase64 string
}

type ExpiringToken struct {
	Value   string
	Exipres int64
}

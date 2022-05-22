package sign

const (
	// KeyNameTimeStamp Timestamp field name
	KeyNameTimeStamp = "timestamp"
	// KeyNameNonceStr temp str field
	KeyNameNonceStr = "nonce_str"
	// KeyNameAppID app id field
	KeyNameAppID = "app_id"
	// KeyNameSign sign field
	KeyNameSign = "sign"
)

// DefaultKeyName Fields required for signature
type DefaultKeyName struct {
	Timestamp string
	NonceStr  string
	AppID     string
	Sign      string
}

func newDefaultKeyName() *DefaultKeyName {
	return &DefaultKeyName{
		Timestamp: KeyNameTimeStamp,
		NonceStr:  KeyNameNonceStr,
		AppID:     KeyNameAppID,
		Sign:      KeyNameSign,
	}
}

// SetKeyNameTimestamp set timestamp
func (d *DefaultKeyName) SetKeyNameTimestamp(name string) {
	d.Timestamp = name
}

// SetKeyNameNonceStr set random string
func (d *DefaultKeyName) SetKeyNameNonceStr(name string) {
	d.NonceStr = name
}

// SetKeyNameAppID set app id
func (d *DefaultKeyName) SetKeyNameAppID(name string) {
	d.AppID = name
}

// SetKeyNameSign set signature
func (d *DefaultKeyName) SetKeyNameSign(name string) {
	d.Sign = name
}

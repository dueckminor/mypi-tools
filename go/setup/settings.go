package setup

type SettingsAlpine struct {
	Version string
	Arch    string
}

type SettingsNetwork struct {
	Hostname string
}

type SettingsWLAN struct {
	SSID     string
	Password string
}

type SettingsSSH struct {
	AuthorizedKeys string
}

type Settings struct {
	Alpine  SettingsAlpine
	Network SettingsNetwork
	WLAN    SettingsWLAN
	SSH     SettingsSSH
}

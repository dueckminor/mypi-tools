package ccu

type DeviceDescription struct {
	Type              string   `xmlrpc:"TYPE"`
	Address           string   `xmlrpc:"ADDRESS"`
	RFAddress         int      `xmlrpc:"RF_ADDRESS"`
	Children          []string `xmlrpc:"CHILDREN"`
	Parent            string   `xmlrpc:"PARENT,omitempty"`
	ParentType        string   `xmlrpc:"PARENT_TYPE,omitempty"`
	Index             int      `xmlrpc:"INDEX,omitempty"`
	AESActive         int      `xmlrpc:"AES_ACTIVE"`
	ParamSets         []string `xmlrpc:"PARAMSETS"`
	Firmware          string   `xmlrpc:"FIRMWARE,omitempty"`
	AvailableFirmware string   `xmlrpc:"AVAILABLE_FIRMWARE,omitempty"`
	Updateable        int      `xmlrpc:"UPDATABLE,omitempty"`
	Version           int      `xmlrpc:"VERSION"`
	Flags             int      `xmlrpc:"FLAGS"`
	LinkSourceRoles   string   `xmlrpc:"LINK_SOURCE_ROLES,omitempty"`
	LinkTargetRoles   string   `xmlrpc:"LINK_TARGET_ROLES,omitempty"`
	Direction         int      `xmlrpc:"DIRECTION,omitempty"`
	Group             string   `xmlrpc:"GROUP,omitempty"`
	Team              string   `xmlrpc:"TEAM,omitempty"`
	TeamTag           string   `xmlrpc:"TEAM_TAG,omitempty"`
	TeamChannels      []string `xmlrpc:"TEAM_CHANNELS"`
	Interface         string   `xmlrpc:"INTERFACE,omitempty"`
	Roaming           int      `xmlrpc:"ROAMING,omitempty"`
	RxMode            int      `xmlrpc:"RX_MODE,omitempty"`
}

type ParameterDescription struct {
	Type string `xmlrpc:"TYPE"`
}

type ParamsetDescription map[string]*ParameterDescription

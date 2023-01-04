package fdisk

type windowsPartition struct {
	win32DiskPartition wmiObject
}

func (part *windowsPartition) GetDeviceName() string {
	return ""
}

func (part *windowsPartition) GetDeviceFileName() string {
	return ""
}

func (part *windowsPartition) GetName() string {
	return ""
}

func (part *windowsPartition) GetSize() int64 {
	return 0
}

func (part *windowsPartition) GetMountPoint() (mountPoint string, err error) {
	ps, err := NewPowerShell()
	if err != nil {
		return "", err
	}

	obj, err := ps.WmiQueryObject("ASSOCIATORS OF {Win32_DiskPartition.DeviceID='%s'} WHERE AssocClass = Win32_LogicalDiskToPartition", part.win32DiskPartition.GetString("DeviceID"))
	if err != nil {
		return "", err
	}

	return obj.GetString("DeviceID"), nil
}

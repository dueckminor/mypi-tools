package fdisk

import (
	"fmt"
	"strconv"
	"strings"
)

type windowsDisk struct {
	win32DiskDrive wmiObject
}

func (disk *windowsDisk) GetDeviceName() string {
	return fmt.Sprintf("Disk %v", disk.win32DiskDrive["Index"])
}
func (disk *windowsDisk) GetDeviceFileName() string {
	return disk.win32DiskDrive.GetString("DeviceID")
}
func (disk *windowsDisk) GetName() string {
	return fmt.Sprintf("Disk %v", disk.win32DiskDrive["Index"])
}
func (disk *windowsDisk) GetSize() int64 {
	return disk.win32DiskDrive.GetInt64("Size")
}
func (disk *windowsDisk) IsRemovable() bool {
	return disk.win32DiskDrive.GetString("MediaType") == "Removable Media"
}
func (disk *windowsDisk) GetPartitions() (result []Partition, err error) {
	ps, err := NewPowerShell()
	if err != nil {
		return nil, err
	}

	win32DiskPartitions, err := ps.WmiQueryArray("SELECT * FROM Win32_DiskPartition WHERE DiskIndex=%d", disk.win32DiskDrive.GetInt64("Index"))
	if err != nil {
		return nil, err
	}

	for _, win32DiskPartition := range win32DiskPartitions {
		result = append(result, &windowsPartition{
			win32DiskPartition: win32DiskPartition,
		})
	}

	return result, nil
}
func (disk *windowsDisk) InitializePartitions(Type string, partitionInfos ...PartitionInfo) error {
	return nil
}
func (disk *windowsDisk) CreatePartitions(partitionInfos ...PartitionInfo) error {
	return nil
}

// GetRemovableDisks return all removable disks
func GetDisks() (result []Disk, err error) {
	ps, err := NewPowerShell()
	if err != nil {
		return nil, err
	}
	defer ps.Close()
	win32DiskDrives, err := ps.GetArray("Get-WmiObject Win32_DiskDrive")
	if err != nil {
		return nil, err
	}

	for _, win32DiskDrive := range win32DiskDrives {
		if win32DiskDrive["Size"] != nil {
			result = append(result, &windowsDisk{
				win32DiskDrive: win32DiskDrive,
			})
		}
	}

	return result, nil
}

func getDeviceID(name string) (deviceID string, err error) {
	index := ""
	if strings.HasPrefix(name, "\\\\.\\PHYSICALDRIVE") {
		index = name[len("\\\\.\\PHYSICALDRIVE"):]
	} else if strings.HasPrefix(name, "Disk ") {
		index = name[5:]
	}
	if len(index) == 0 {
		return "", fmt.Errorf("'%s' is no valid disk name", name)
	}
	i, err := strconv.ParseInt(index, 10, 0)
	if err != nil || index != strconv.FormatInt(i, 10) {
		return "", fmt.Errorf("'%s' is no valid disk name", name)
	}
	return "\\\\.\\PHYSICALDRIVE" + index, nil
}

func sqlQuoteString(str string) (quoted string, err error) {
	b := strings.Builder{}
	b.WriteRune('\'')
	for _, r := range str {
		switch r {
		case '\\':
			b.WriteRune(r)
		case '"', '\'':
			return "", fmt.Errorf("can quote quoted")
		}
		b.WriteRune(r)
	}
	b.WriteRune('\'')
	return b.String(), nil
}

// GetDisk returns information about a single disk
func GetDisk(diskName string) (Disk, error) {

	deviceID, err := getDeviceID(diskName)
	if err != nil {
		return nil, err
	}

	deviceID, err = sqlQuoteString(deviceID)
	if err != nil {
		return nil, err
	}

	ps, err := NewPowerShell()
	if err != nil {
		return nil, err
	}
	defer ps.Close()

	cmdline := "Get-WmiObject -query \"SELECT * FROM Win32_DiskDrive WHERE DeviceID=" + deviceID + "\""
	win32DiskDrive, err := ps.GetObject(cmdline)
	if err != nil {
		return nil, err
	}

	return &windowsDisk{win32DiskDrive: win32DiskDrive}, nil
}

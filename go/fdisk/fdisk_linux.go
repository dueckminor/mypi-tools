package fdisk

import (
	"os/exec"
	"strconv"
	"strings"
)

// nolint: unused
type linuxSFDiskPartition struct {
	Node  string `json:"node"`
	Start string `json:"start"`
	Size  string `json:"size"`
	Type  string `json:"type"`
}

// nolint: unused
type linuxSFDiskPartitionTable struct {
	Label      string                 `json:"label"`
	Id         string                 `json:"id"`
	Device     string                 `json:"device"`
	Unit       string                 `json:"unit"`
	Partitions []linuxSFDiskPartition `json:"partitions"`
}

// nolint: unused
type linuxSFDiskInfo struct {
	PartitionTable linuxSFDiskPartitionTable `json:"partitiontable"`
}

// nolint: unused
func makeDevName(diskName string) string {
	if strings.HasPrefix(diskName, "/dev/") {
		return diskName
	}
	return "/dev/" + diskName
}

// nolint: unused
type linuxDisk struct {
	devName string
	size    int64
}

// nolint: unused
func (d *linuxDisk) GetDeviceName() string {
	return d.devName
}

// nolint: unused
func (d *linuxDisk) GetDeviceFileName() string {
	return d.devName
}

// nolint: unused
func (d *linuxDisk) GetName() string {
	return d.devName
}

// nolint: unused
func (d *linuxDisk) GetSize() int64 {
	return 0
}

// nolint: unused
func (d *linuxDisk) IsRemovable() bool {
	return false
}

// nolint: unused
func (d *linuxDisk) GetPartitions() ([]Partition, error) {
	return nil, nil
}

// nolint: unused
func (d *linuxDisk) InitializePartitions(Type string, partitionInfos ...PartitionInfo) error {
	return nil
}

// nolint: unused
func (d *linuxDisk) CreatePartitions(partitionInfos ...PartitionInfo) error {
	return nil
}

// nolint: unused
func newLinuxDisk(diskName string) (Disk, error) {
	diskName = makeDevName(diskName)

	out, err := exec.Command("blockdev", "--getsize64", diskName).Output()
	if err != nil {
		return nil, err
	}
	size, err := strconv.ParseInt(string(out), 10, 64)
	if err != nil {
		return nil, err
	}

	return &linuxDisk{
		devName: diskName,
		size:    size,
	}, nil
}

func GetDisks() ([]Disk, error) {
	return nil, nil
}

// GetDisk returns information about a single disk
func GetDisk(diskName string) (Disk, error) {
	return nil, nil
}

package fdisk

import (
	"bytes"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"howett.net/plist"
)

// 1            8192          532479   256.0 MiB   0700  Microsoft basic data
// 2          532480         4390911   1.8 GiB     8300  Linux filesystem
// 3         4390912        53945343   23.6 GiB    8300  Linux filesystem
// 4        53945344        62333951   4.0 GiB     8200  Linux swap

type PList struct {
	content map[string]interface{}
}

func (obj PList) followPath(path ...interface{}) (interface{}, error) {
	var where interface{}
	where = obj.content
	for _, pathElement := range path {
		if name, ok := pathElement.(string); ok {
			if dict, ok := where.(map[string]interface{}); ok {
				where, ok = dict[name]
				if !ok {
					return nil, nil
				}
			}
		}
		if index, ok := pathElement.(int); ok {
			if arr, ok := where.([]interface{}); ok {
				where = arr[index]
				if !ok {
					return nil, nil
				}
			}
		}
	}
	return where, nil
}

func (obj PList) getSliceOfStrings(path ...interface{}) (result []string, err error) {
	where, err := obj.followPath(path...)
	if arr, ok := where.([]interface{}); ok {
		result = make([]string, len(arr))
		for index, item := range arr {
			if name, ok := item.(string); ok {
				result[index] = name
			} else {
				return nil, nil
			}
		}
		return result, nil
	}
	return nil, nil
}

func (obj PList) getString(path ...interface{}) (result string, err error) {
	where, err := obj.followPath(path...)
	if result, ok := where.(string); ok {
		return result, nil
	}
	return "", nil
}

func (obj PList) getInt(path ...interface{}) (result int64, err error) {
	where, err := obj.followPath(path...)
	if result, ok := where.(int64); ok {
		return result, nil
	}
	if result, ok := where.(uint64); ok {
		return int64(result), nil
	}
	return 0, nil
}

func (obj PList) getBool(path ...interface{}) (result bool, err error) {
	where, err := obj.followPath(path...)
	if result, ok := where.(bool); ok {
		return result, nil
	}
	return false, nil
}

func callDiskutil(args ...string) (result PList, err error) {
	result = PList{}
	out, err := exec.Command("diskutil", args...).Output()
	if err != nil {
		return result, err
	}
	decoder := plist.NewDecoder(bytes.NewReader(out))
	err = decoder.Decode(&result.content)
	if err != nil {
		return result, err
	}
	return result, nil
}

type macosBlockDevice struct {
	deviceName     string
	deviceFileName string
	info           PList
}

func (dev macosBlockDevice) GetDeviceName() string {
	return dev.deviceName
}
func (dev macosBlockDevice) GetDeviceFileName() string {
	return dev.deviceFileName
}
func (dev macosBlockDevice) GetName() string {
	s, _ := dev.info.getString("IORegistryEntryName")
	return s
}
func (dev macosBlockDevice) GetSize() int64 {
	s, _ := dev.info.getInt("TotalSize")
	return s
}

type macosPartition struct {
	macosBlockDevice
}

func macosNewPartition(deviceName string) (result macosPartition, err error) {
	result.info, err = callDiskutil("info", "-plist", deviceName)
	if err != nil {
		return macosPartition{}, nil
	}
	result.deviceName = deviceName
	return result, err
}

func (partition macosPartition) GetMountPoint() (mountPoint string, err error) {
	return partition.info.getString("MountPoint")
}

type macosDisk struct {
	macosBlockDevice
	partitionNames []string
	partitions     []Partition
}

func newDiskMacos(deviceName string) (result macosDisk, err error) {
	result.info, err = callDiskutil("info", "-plist", deviceName)
	if err != nil {
		return macosDisk{}, err
	}

	deviceName, err = result.info.getString("DeviceIdentifier")
	if err != nil {
		return macosDisk{}, err
	}

	partinfo, err := callDiskutil("list", "-plist", deviceName)

	result.deviceName = deviceName
	result.partitionNames, err = partinfo.GetPartitionNames(deviceName)
	if err != nil {
		return macosDisk{}, err
	}
	return result, err
}

func (dev macosDisk) IsRemovable() bool {
	ejectableOnly, err := dev.info.getBool("EjectableOnly")
	if err == nil && !ejectableOnly {
		return false
	}
	removeable, _ := dev.info.getBool("Removable")
	removeableMedia, _ := dev.info.getBool("RemovableMediaOrExternalDevice")
	return removeable || removeableMedia
}

func (dev macosDisk) GetPartitions() ([]Partition, error) {
	if len(dev.partitionNames) > len(dev.partitions) {
		partitions := make([]Partition, len(dev.partitionNames))
		for index, name := range dev.partitionNames {
			partition, err := macosNewPartition(name)
			if err != nil {
				return nil, err
			}
			partitions[index] = partition
		}
		dev.partitions = partitions
	}
	return dev.partitions, nil
}

func (dev macosDisk) InitializePartitions(Type string, partitionInfos ...PartitionInfo) error {
	if !dev.IsRemovable() {
		return os.ErrPermission
	}

	args := []string{
		"partitionDisk",
		dev.deviceName,
		Type,
	}

	HaveRemaining := false
	for _, partitionInfo := range partitionInfos {
		Size := "R"
		if partitionInfo.Size > 0 {
			Size = strconv.FormatInt(partitionInfo.Size, 10) + "B"
		} else {
			HaveRemaining = true
		}
		args = append(args, partitionInfo.Format, partitionInfo.Name, Size)
	}
	if !HaveRemaining {
		args = append(args, "Free Space", "Free Space", "R")
	}

	cmd := exec.Command("diskutil", args...)
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	err = cmd.Wait()
	if err != nil {
		return err
	}
	return nil
}

func (dev macosDisk) CreatePartitions(partitionInfos ...PartitionInfo) error {
	return os.ErrInvalid
}

func getPartitionNames(diskDevice string, disksAndPartitions []string) (result []string) {
	for _, diskOrPartition := range disksAndPartitions {
		if (len(diskOrPartition) > len(diskDevice)) &&
			strings.HasPrefix(diskOrPartition, diskDevice) {
			result = append(result, diskOrPartition)
		}
	}
	return result
}

func (plist PList) GetDiskNames() ([]string, error) {
	return plist.getSliceOfStrings("WholeDisks")
}

func (plist PList) GetPartitionNames(diskName string) ([]string, error) {
	disksAndPartitions, err := plist.getSliceOfStrings("AllDisks")
	if err != nil {
		return nil, err
	}
	return getPartitionNames(diskName, disksAndPartitions), nil
}

func callDiskutilList(args ...string) (plist PList, err error) {
	cmd := []string{"list", "-plist"}
	cmd = append(cmd, args...)
	return callDiskutil(cmd...)
}

// GetRemovableDisks return all removable disks
func GetDisks() ([]Disk, error) {
	d, err := callDiskutilList()
	if err != nil {
		return nil, err
	}

	diskNames, err := d.GetDiskNames()
	if err != nil {
		return nil, err
	}

	result := make([]Disk, 0)

	for _, diskName := range diskNames {
		disk, err := newDiskMacos(diskName)
		if err != nil {
			continue
		}
		result = append(result, disk)
	}

	return result, nil
}

// GetDisk returns information about a single disk
func GetDisk(diskName string) (Disk, error) {
	return newDiskMacos(diskName)
}

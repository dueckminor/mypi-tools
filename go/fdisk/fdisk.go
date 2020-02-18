package fdisk

// BlockDevice represents a Partition or a Hard-Disk
type BlockDevice interface {
	GetDeviceName() string
	GetDeviceFileName() string
	GetName() string
	GetSize() int64
}

// Partition represents a Partition on a Hard-Disk
type Partition interface {
	BlockDevice
	GetMountPoint() (mountPoint string, err error)
}

type PartitionInfo struct {
	Size   int64
	Type   int
	Format string
	Name   string
}

// Disk represents a Hard-Disk
type Disk interface {
	BlockDevice
	IsRemovable() bool
	GetPartitions() ([]Partition, error)
	InitializePartitions(Type string, partitionInfos ...PartitionInfo) error
}

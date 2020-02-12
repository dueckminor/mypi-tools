#!/usr/bin/env bash

set -e

DISK=disk2
ALPINE_VERSION=3.11.3
ALPINE_ARCH=aarch64
RPI_HOSTNAME=rpi2

DEVICE="${1:-mmcblk0}"
BOOT_DEVICE="${DEVICE}p1"
ROOT_DEVICE="${DEVICE}p2"

DIR_THIS="$(cd "$(dirname "${BASH_SOURCE[0]}")"; pwd)"
DIR_ROOT="$(cd "${DIR_THIS}/../.."; pwd)"
DIR_TMP="${DIR_ROOT}/tmp"
DIR_APKOVL="${DIR_TMP}/apkovl"

DIR_MNT="/Volumes/RPI-BOOT"

umount_all() {
    mount | grep "^/dev/${DISK}s" | cut -d' ' -f 1 | xargs diskutil unmount force
}

umount_all

diskutil partitionDisk "${DISK}" MBR FAT32 RPI-BOOT 256M FREE FREE R

pushd "${DIR_MNT}"
    tar xvf "${DIR_ROOT}/tmp/alpine-rpi-${ALPINE_VERSION}-${ALPINE_ARCH}.tar.gz" --no-same-owner
popd

mkdir -p "${DIR_APKOVL}"
rm -rf "${DIR_APKOVL:?}/*"

cp -a "${DIR_THIS}/apkovl/"* "${DIR_APKOVL}"

cat > "${DIR_APKOVL}/mypi-setup/config" <<-__EOF__
BOOT_DEVICE_NAME="${BOOT_DEVICE}"
ROOT_DEVICE_NAME="${ROOT_DEVICE}"
#WLAN_SSID="${WLAN_SSID}"
#WLAN_PASSWORD="${WLAN_PASSWORD}"
__EOF__

cat > "${DIR_APKOVL}/mypi-setup/answer.txt" <<-__EOF__
KEYMAPOPTS="us us"
HOSTNAMEOPTS="-n ${RPI_HOSTNAME}"
INTERFACESOPTS="auto lo
iface lo inet loopback

auto eth0
iface eth0 inet dhcp
    hostname ${RPI_HOSTNAME}
"
TIMEZONEOPTS="-z UTC"
PROXYOPTS=none
APKREPOSOPTS="-f"
SSHDOPTS="-c openssh"
NTPOPTS="-c busybox"
APKCACHEOPTS="none"
LBUOPTS="none"
__EOF__

pushd "${DIR_APKOVL}"
    tar czf "${DIR_MNT}/${RPI_HOSTNAME}.apkovl.tar.gz" .
popd >/dev/null

umount_all

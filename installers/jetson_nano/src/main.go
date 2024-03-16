// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	_ "embed"
	"os"
	"path/filepath"

	"github.com/siderolabs/go-copy/copy"
	"github.com/siderolabs/talos/pkg/machinery/overlay"
	"github.com/siderolabs/talos/pkg/machinery/overlay/adapter"
	"golang.org/x/sys/unix"
)

const (
	// References
	// - https://github.com/u-boot/u-boot/blob/v2021.10/configs/p3450-0000_defconfig#L8
	// - https://github.com/u-boot/u-boot/blob/v2021.10/include/configs/tegra-common.h#L53
	// - https://github.com/u-boot/u-boot/blob/v2021.10/include/configs/tegra210-common.h#L49
	dtb = "nvidia/tegra210-p3450-0000.dtb"
)

func main() {
	adapter.Execute(&JetsonNanoInstaller{})
}

type JetsonNanoInstaller struct{}

type jetsonNanoExtraOptions struct{}

func (i *JetsonNanoInstaller) GetOptions(extra jetsonNanoExtraOptions) (overlay.Options, error) {
	return overlay.Options{
		Name: "jetson_nano",
		KernelArgs: []string{
			"console=tty0",
			"console=ttyS0,115200",
			"sysctl.kernel.kexec_load_disabled=1",
			"talos.dashboard.disabled=1",
		},
	}, nil
}

func (i *JetsonNanoInstaller) Install(options overlay.InstallOptions[jetsonNanoExtraOptions]) error {
	var f *os.File

	f, err := os.OpenFile(options.InstallDisk, os.O_RDWR|unix.O_CLOEXEC, 0o666)
	if err != nil {
		return err
	}

	defer f.Close() //nolint:errcheck

	// NB: In the case that the block device is a loopback device, we sync here
	// to ensure that the file is written before the loopback device is
	// unmounted.
	err = f.Sync()
	if err != nil {
		return err
	}

	src := filepath.Join(options.ArtifactsPath, "arm64/dtb", dtb)
	dst := filepath.Join(options.MountPrefix, "/boot/EFI/dtb", dtb)

	err = os.MkdirAll(filepath.Dir(dst), 0o600)
	if err != nil {
		return err
	}

	err = copy.File(src, dst)
	if err != nil {
		return err
	}

	return nil
}

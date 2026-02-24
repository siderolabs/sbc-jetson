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
	// References for Jetson Orin Nano (Tegra234)
	// - https://github.com/u-boot/u-boot/blob/master/configs/p3767-0000_defconfig
	// - https://developer.nvidia.com/embedded/learn/get-started-jetson-orin-nano-devkit
	dtb = "nvidia/tegra234-p3768-0000+p3767-0000.dtb"
)

func main() {
	adapter.Execute(&JetsonOrinNanoInstaller{})
}

type JetsonOrinNanoInstaller struct{}

type jetsonOrinNanoExtraOptions struct{}

func (i *JetsonOrinNanoInstaller) GetOptions(extra jetsonOrinNanoExtraOptions) (overlay.Options, error) {
	return overlay.Options{
		Name: "jetson_orin_nano",
		KernelArgs: []string{
			"console=tty0",
			"console=ttyS0,115200",
			"sysctl.kernel.kexec_load_disabled=1",
			"talos.dashboard.disabled=1",
			// Orin Nano specific optimizations
			"tegra_fbmem=0x800000@0x278800000",
			"lut_mem=0x2008@0x278000000",
		},
	}, nil
}

func (i *JetsonOrinNanoInstaller) Install(options overlay.InstallOptions[jetsonOrinNanoExtraOptions]) error {
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
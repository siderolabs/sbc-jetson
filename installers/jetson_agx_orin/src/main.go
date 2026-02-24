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
	// References for Jetson AGX Orin (Tegra234)
	// - https://github.com/u-boot/u-boot/blob/master/configs/p3701-0000_defconfig
	// - https://developer.nvidia.com/embedded/learn/get-started-jetson-agx-orin-devkit
	dtb = "nvidia/tegra234-p3701-0000+p3737-0000.dtb"
)

func main() {
	adapter.Execute(&JetsonAGXOrinInstaller{})
}

type JetsonAGXOrinInstaller struct{}

type jetsonAGXOrinExtraOptions struct{}

func (i *JetsonAGXOrinInstaller) GetOptions(extra jetsonAGXOrinExtraOptions) (overlay.Options, error) {
	return overlay.Options{
		Name: "jetson_agx_orin",
		KernelArgs: []string{
			"console=tty0",
			"console=ttyS0,115200",
			"sysctl.kernel.kexec_load_disabled=1",
			"talos.dashboard.disabled=1",
			// AGX Orin specific optimizations
			"tegra_fbmem=0x2000000@0x278000000",
			"lut_mem=0x2008@0x276000000",
			// Enhanced memory and performance settings for AGX Orin
			"nvdec_enabled",
			"nvidia_drm.modeset=1",
		},
	}, nil
}

func (i *JetsonAGXOrinInstaller) Install(options overlay.InstallOptions[jetsonAGXOrinExtraOptions]) error {
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
# sbc-jetson

This repo provides the overlay for NVIDIA Jetson Talos image.

Support for this will be based on [NVIDIA embedded lifecycle](https://developer.nvidia.com/embedded/lifecycle).

## Supported Overlay

| Overlay Name | Board       | Description         |
| ------------ | ----------- | ------------------- |
| jetson_nano  | Jetson Nano | Jetson Nano overlay |

For installation instructions (including flashing U-Boot), see the
[Jetson Nano guide](https://docs.siderolabs.com/talos/v1.14/platform-specific-installations/single-board-computers/jetson_nano).

## Artifacts

| Path                                                 | Description                              |
| ---------------------------------------------------- | ---------------------------------------- |
| `artifacts/arm64/u-boot/jetson_nano/u-boot.bin`      | Mainline U-Boot (`p3450-0000_defconfig`) |
| `artifacts/arm64/dtb/nvidia/tegra210-p3450-0000.dtb` | Device tree from the Talos kernel        |

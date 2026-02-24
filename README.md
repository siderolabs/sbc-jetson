# sbc-jetson

This repository provides **overlay support for NVIDIA Jetson single-board computers in Talos Linux**. It creates custom installers that integrate NVIDIA Jetson hardware with Talos Linux, including proper U-Boot bootloader support, device tree configurations, and hardware-specific kernel parameters.

Support for this project is based on [NVIDIA embedded lifecycle](https://developer.nvidia.com/embedded/lifecycle).

## Project Overview

**sbc-jetson** enables Talos Linux to run on NVIDIA Jetson devices by:

- **Custom Installers**: Device-specific installers that handle hardware configuration during Talos installation
- **U-Boot Integration**: Tegra-optimized U-Boot bootloader support using mainline U-Boot v2026.01
- **Device Tree Support**: Proper DTB (Device Tree Blob) installation for each Jetson variant using NVIDIA validated device trees
- **Hardware Configuration**: Optimized kernel boot parameters and system settings for Jetson hardware

## Supported Overlays

| Overlay Name     | Board            | Tegra SoC | U-Boot Config | Description              |
| ---------------- | ---------------- | --------- | ------------- | ------------------------ |
| jetson_nano      | Jetson Nano      | Tegra210  | p3450-0000    | Jetson Nano overlay      |
| jetson_orin_nano | Jetson Orin Nano | Tegra234  | N/A*          | Jetson Orin Nano overlay |
| jetson_agx_orin  | AGX Orin         | Tegra234  | N/A*          | Jetson AGX Orin overlay  |

_*Orin models currently use DTBs only (no custom U-Boot) since mainline lacks Tegra234 configs_

### Current Architecture Details

- **U-Boot Version**: v2026.01 (mainline)
- **Multi-Model Support**: All three Jetson models supported simultaneously
- **Device Trees**: 
  - Nano: `nvidia/tegra210-p3450-0000.dtb` (built by U-Boot)
  - Orin Nano: `nvidia/tegra234-p3768-0000+p3767-0000.dtb` (from NVIDIA L4T r36.4.4)
  - AGX Orin: `nvidia/tegra234-p3701-0000+p3737-0000.dtb` (from NVIDIA L4T r36.4.4)
- **Boot Configuration**: Serial console, security hardening, model-specific optimizations
- **Platform**: ARM64/AArch64
- **DTB Source**: Mainline U-Boot for Nano, NVIDIA L4T BSP for Orin models (since mainline lacks Tegra234 configs)

## Development

### Project Structure

```
├── artifacts/u-boot/          # U-Boot build configuration (multi-model support)
├── installers/                # Device-specific installers
│   ├── jetson_nano/           # Jetson Nano installer implementation
│   ├── jetson_orin_nano/      # Jetson Orin Nano installer implementation
│   └── jetson_agx_orin/       # Jetson AGX Orin installer implementation
├── profiles/                  # Talos image profiles
│   ├── jetson_nano/           # Nano-specific profile
│   ├── jetson_orin_nano/      # Orin Nano-specific profile
│   └── jetson_agx_orin/       # AGX Orin-specific profile
└── internal/                  # Shared components
```

### Current Implementation Notes

**Jetson Nano**: Uses mainline U-Boot with built-in Tegra210 support. Both U-Boot binary and device trees are built from source.

**Orin Models**: Currently use **DTBs only** from NVIDIA L4T r36.4.4 BSP. No custom U-Boot is built since mainline U-Boot v2026.01 doesn't include Tegra234 board configurations yet. These models rely on:
- Stock bootloader from NVIDIA's BSP/flashed firmware
- NVIDIA's validated device trees extracted during build
- Talos overlay system for OS-level integration

### Updating U-Boot Version

Current U-Boot version can be updated in [`artifacts/u-boot/pkg.yaml`](artifacts/u-boot/pkg.yaml):

```yaml
# Current: mainline v2026.01
- url: https://github.com/u-boot/u-boot/archive/refs/tags/v2026.01.tar.gz

# Update to newer version:
- url: https://github.com/u-boot/u-boot/archive/refs/tags/v2026.04.tar.gz
```

**Version Compatibility Reference**:

| Jetson Model | Tegra SoC | U-Boot Config | Status in Mainline  |
| ------------ | --------- | ------------- | ------------------- |
| Nano         | Tegra210  | p3450-0000    | ✅ Supported         |
| Orin Nano    | Tegra234  | N/A           | ❌ Using NVIDIA DTBs |
| AGX Orin     | Tegra234  | N/A           | ❌ Using NVIDIA DTBs |

**NVIDIA L4T Integration**: For Orin models, device trees are extracted from NVIDIA's L4T r36.4.4 BSP during build process.

### Building

Build all components:
```bash
make
```

Build specific installer:
```bash
make jetson-nano-installer
```


### Dev Dependencies

- golang
- docker
- jq
- graphviz (dot)

recommended:
- bash-completion
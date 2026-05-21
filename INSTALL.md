## Quick Install (Recommended)

A single `install.sh` script handles both flows: run it via curl/wget (below) to download from GitHub Releases, or run it from an extracted tarball to install the `dguide` binary already on disk. For remote installs, it detects your OS (macOS or Linux) and CPU architecture (x86_64, arm64, 32-bit), verifies the release checksum, and places `dguide` on your PATH.

#### Option 1 — curl, always installs the latest release

```sh
curl -fsSL https://raw.githubusercontent.com/Appdynamics/dguide/main/install.sh | sh
```

#### Option 2 — curl, pinned to a specific version

```sh
curl -fsSL https://raw.githubusercontent.com/Appdynamics/dguide/main/install.sh | sh -s -- --version 0.2.0
```

Replace `0.2.0` with any published release version (with or without the leading `v`).

#### Option 3 — wget (when curl is not available)

```sh
wget -qO- https://raw.githubusercontent.com/Appdynamics/dguide/main/install.sh | sh
```

#### Option 4 — download and inspect before running

If you prefer to review the script before executing it:

```sh
# Download the installer
curl -fsSL https://raw.githubusercontent.com/Appdynamics/dguide/main/install.sh -o install.sh

# (Optional) inspect the script
cat install.sh

# Install latest release
sh install.sh

# Or install a specific version
sh install.sh --version 0.2.0

# Show usage / help
sh install.sh --help
```

## Manual Installation

dguide takes less than few seconds to install. Download the CLI and follow below steps:
#### 1. Untar the package

```sh
tar -xzf dguide_${DGUIDE_VERSION}_linux_amd64.tar.gz -C /opt/dguide/
```

#### 2. Run the install script

```sh
/opt/dguide/install.sh
```
#### 3. Run the CLI

```sh
dguide version
```

#### Note:

_dguide_ might fail to run, due to insufficient permissions or untrusted software warnings on different platforms. Here are some workarounds you can use to execute the CLI.

On Linux ->
```sh
chmod +x dguide
./dguide
```
On [Linux(Redhat or similar with SELinux)](https://docs.redhat.com/en/documentation/red_hat_enterprise_linux/7/html/selinux_users_and_administrators_guide/sect-security-enhanced_linux-introduction-selinux_modes#sect-Security-Enhanced_Linux-Introduction-SELinux_Modes) -> 
```sh
#Check if SELinux is blocking the binary 
~$ sudo getenforce
Enforcing

# If needed, modify the SeLinux to `Permissive` 
```

On MacOs ->
```
sudo xattr -d com.apple.quarantine dguide
sudo ./install.sh
```

The below demo shows how to install and run the dguide CLI. 

![dguide-demo](img/dguide_quick_demo1742210637241.gif)


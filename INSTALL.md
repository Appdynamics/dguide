## 🛠️ Installation
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


## >_ dguide
![dguide-logo](img/dguide.png)
🪄Seamless Agent Troubleshooting: Save Time and Effort in collecting agent diagnostic information from your environment 

What it does ?
 - Automates steps to collect troubleshooting information for various language agents , no more commands to run and gather information about agents, additional configurations or even agent logs 🚀
 - Option to generate and execute custom commands as needed on the fly.
  
## Pre-requisite

Supported OS and Arch 
| OS            | 
| ------------- | 
| MacOS(Darwin) | 
| Linux         |

| Architecture         | 
| -------------------- |
| x86_64 (32 & 64bit ) | 
| arm64                |

e.g. `dguide_{version}_linux_amd64.tar.gz` can run on any Linux OS with CPU arch mathing to x86_64 (64-bit)


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

## 🤝 Steps to run the tool 

Run dguide help command to see list of supported commands.
```
~$ dguide help
```

Select he sub-command you want to run and collect information for a specific agent installed in your env.

For example, the below command prints the top 10 CPU intensive processes on your system.
``` 
~$ dguide run ptop
```

The below demo shows how to install and run the dguide CLI. 

![dguide-demo](img/dguide_quick_demo1742210637241.gif)

## Submitting changes & Join the Clan ⚔️

Generally, you should fork this repository, make changes in your own fork, and then submit a [🔗pull request](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/working-with-forks/fork-a-repo). 


`⚠️Note:`This project is open-source and hosted on our official GitHub site. We kindly request that all complaints, suggestions, and error reports be directed exclusively through the project's GitHub page. We will not be able to address concerns reported through any other channel or medium. Thank you for your understanding and cooperation.  

## Support & Report Issues

[Submit your issues, complains and suggestion](ISSUES.md)
#!/bin/sh

# ==============================================================================
# Script Name:    sysstat.sh
# Description:    This script performs system monitoring tasks including Top CPU
#                 /MEM profile, process status, system limits, and
#                 gathering vmstat information.
# Author:         Jayanta Mohanty
# Creation Date:  2024-07-08
# ==============================================================================
# NOTE: 
# 1) use /bin/sh: Ensures maximum compatibility, available on most of linux or unix-like machines.
# 2) pre-requisites packages
#       - procps (alpine, debian)
#       - procps-ng (redhat) 
# 3) Dont run ps -aef | grep <something> (# Why? Respect customer env , avoid any regex ,awk, grep commmands)
# 4) Take input to script using ENV variables i.e. before running the command , export required ENV variable and run. 
#==========================================DO NOT CHANGE!================================
#output_file="appd_dl_diagnose.log"

output_file="/dev/stdout"

add_header() {
    echo "==================================================" > "$output_file"
    echo "                System CPU & MEM Stats            " >> "$output_file"
    echo "==================================================" >> "$output_file"
    echo "Generated on: $(date)" >> "$output_file"
    echo "==================================================" >> "$output_file"
    echo "" >> "$output_file"
}

log() {
  local level=$1
  local message=$2
  local timestamp=$(date +"%Y-%m-%d %H:%M:%S")
  case $level in
    info)
      echo "${timestamp} [INFO] ${message}" 
      ;;
    debug)
      echo "${timestamp} [DEBUG] ${message}"
      ;;
    error)
      echo "${timestamp} [ERROR] ${message}"
      ;;
    *)
      echo "Invalid log level. Supported levels are info, debug, error."
      ;;
  esac
}

logf() {
    local section_title="$1"
    local inc_line="${2:-true}"  # Default to true if not provided

    if [ "$inc_line" = true ]; then
        echo "--------------------------------------------------" >> "$output_file"
    fi
    echo "$section_title" >> "$output_file"
    if [ "$inc_line" = true ]; then
        echo "--------------------------------------------------" >> "$output_file"
    fi
}

#==========================================DO NOT CHANGE!================================

#  CPU stats
gather_cpu_stats() {
    logf "Top 10 CPU consuming processes " false
    echo "USER             PID    PPID  %CPU %MEM    STIME     TIME    VSZ        RSS  COMMAND" >> "$output_file"
    for i in $(seq 1 3)
    do
        logf "\nTimestamp: $(date +%T)\n" false
        
        ps -A -eo "user pid ppid pcpu pmem stime time vsz rss args" | sort -k 4 -r | head -10 >> "$output_file"

        # Wait for 5 seconds before the next iteration, unless it's the last one
        if [ $i -lt 3 ]; then
            sleep 5
        fi
    done
}

gather_mem_stats() {
    logf "Top 10 MEM intensive processes " 
    ps -A -eo "user pid ppid pcpu pmem stime time vsz rss args" | sort -k 5 -r | head -10 >> "$output_file"
}


check_process_status() {
    logf "process status (Order by creation time)"
    #UID   PID  PPID   C STIME   TTY           TIME CMD
    ps -A -eo "user pid ppid time args" | sort -k 2 -r | head -10 >> "$output_file"
}

# Function to check system limits
check_system_limits() {
    logf "System limits"
    ulimit -a >> "$output_file"
}

gather_vmstat() {
    logf "VMStat Information"
    if [ "$(uname)" = "Darwin" ]; then
        vm_stat >> "$output_file"
    else
        vmstat -s >> "$output_file"
    fi
}

# Function to check and install vmstat if necessary
# Currently disabled 
# install_vmstat() {
#     if ! command -v vmstat > /dev/null 2>&1; then
#         echo "vmstat not found. Attempting to install..."
#         if [ "$(uname)" = "Darwin" ]; then
#             if ! command -v vm_stat > /dev/null 2>&1; then
#                 echo "vm_stat not found. Please ensure you have the necessary system utilities installed."
#                 exit 1
#             fi
#         elif [ -f /etc/alpine-release ]; then
#             apk add procps
#         elif [ -f /etc/debian_version ]; then
#             sudo apt-get update
#             sudo apt-get install -y procps
#         elif [ -f /etc/redhat-release ]; then
#             sudo yum install -y procps-ng
#         else
#             echo "Unsupported OS. Please install vmstat manually."
#             exit 1
#         fi
#     fi
# }

# Function to orchestrate the script execution
execute() {
    OS=$(uname -s)
    ARCH=$(uname -m)
    log info "Tools env: $OS/$ARCH"
    add_header # Dont remove 
    gather_cpu_stats
    gather_mem_stats
    gather_vmstat
    check_process_status
    check_system_limits

    echo "\nEnd" >> "$output_file"
    log info "Ran successfully!"
}


execute
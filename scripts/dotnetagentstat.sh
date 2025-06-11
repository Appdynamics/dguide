#!/bin/sh

OS=$(uname -s)
output_file="/dev/stdout"
source_dir="/tmp/appd"

# Initialize a count for eligible processes
eligible_count=0
eligible_processes=""
current_pid=$$  # Get the current process ID




add_header() {
	# Clear the output file if it exists
	> "$output_file"
    echo "==================================================" >> "$output_file"
    echo "                .Net Agent Stats            "       >> "$output_file"
    echo "==================================================" >> "$output_file"
    echo "Generated on: $(date +"%Y-%m-%d %H:%M:%S %Z %z")" >> "$output_file"
    echo "==================================================" >> "$output_file"
    echo "" >> "$output_file"
}

capture_system_env() {
	# Capture system-level environment variables
	echo "" >> "$output_file"
	echo "System Level Environment Variables:" >> "$output_file"
	env >> "$output_file"
	echo "-----------------------------------" >> "$output_file"
}


capture_os_details() {
	echo "" >> "$output_file"
	echo "***** Os Details *****" >> "$output_file"
	cat /etc/os-release >> "$output_file"
	echo "######################################" >> "$output_file"
	echo "" >> "$output_file"
}


capture_framework_details() {
	echo "" >> "$output_file"
	echo "***** Dotnet Runtime Details *****" >> "$output_file"
	dotnet --list-runtimes >> "$output_file"
	echo "######################################" >> "$output_file"
	echo "" >> "$output_file"
}

capture_process_details() {
	echo "***** Process Details *****" >> "$output_file"
	# Loop through all process IDs
	for pid in $(ps -e -o pid=); do
		
		
		# Exclude the current process, any bash or sh processes, and system processes (PID <= 100)
		if [ "$pid" -eq "$current_pid" ]; then
			continue
		fi


		# Check if /proc/$pid/environ exists
		if [ -e /proc/$pid/comm ]; then
		
			if cat /proc/$pid/comm | grep -E -q 'bash|sh' ; then
			continue
			fi
			
			# Check if the process contains 'dotnet' or loaded Microsoft DLLs
			if  cat /proc/$pid/comm | grep -q 'dotnet' || cat /proc/$pid/maps | grep -q 'libmscorlib.so\|libcoreclr.so\|Microsoft'; then
				# Get process name
				process_name=$(cat /proc/$pid/comm)
				
				# Append eligible process info to the variable
				eligible_processes="$eligible_processes$process_name (PID: $pid)\n"
				eligible_count=$((eligible_count + 1))
			
				# Check if /proc/$pid/environ exists
				if [ -e /proc/$pid/environ ]; then
					# Get environment variables
					env_vars=$(tr '\0' '\n' < /proc/$pid/environ)
				else
					env_vars="No environment variables available"
				fi

				# Get loaded modules
				modules=$(cat /proc/$pid/maps | awk '$6 != "" {print $6}'  | sort -u)

				# Write to output file
				{
					echo "Process Name: $process_name"
					echo "PID: $pid"
					echo ""
					echo "***** Process $process_name Environment Variables: ***** "
					echo "$env_vars"
					echo ""
					echo ""
					echo "***** Loaded Modules to process $process_name: *****"
					echo "$modules"
					echo "----------------------------------------------------------------------------"
					echo ""
					echo ""
				} >> "$output_file"

			else
				env_vars="Process not available"
			fi  
		fi
	done
	echo "######################################" >> "$output_file"
}

copy_log_files() {
	# Check if /tmp/appd/ exists and copy its contents to the dotnet folder
	if [ -d "$source_dir" ]; then
		cp -r "$source_dir/"* "${PWD}"
		echo "Copied contents from $source_dir to ${PWD}."
	else
		echo "Log does not exist."
	fi
}

print_summary() {
	# Print eligible processes and total count at the end
	echo "" >> "$output_file"
	echo "***** Eligible Processes *****" >> "$output_file"
	echo "$eligible_processes" >> "$output_file"
	echo "Total Eligible Processes: $eligible_count" >> "$output_file"

	echo "Process information collected in $output_file"

}

# Function to orchestrate the script execution
execute() {

    add_header # Dont remove

	capture_system_env
    capture_os_details
	capture_framework_details
	capture_process_details
	#@jai : use global flag -l(--logpath) to enable copy agent log feature.
	#copy_log_files
	print_summary
	echo "" >> "$output_file"
	echo "" >> "$output_file"
    echo "==================================================" >> "$output_file"
    echo "END" >> "$output_file"
    echo "==================================================" >> "$output_file"
}
# Supports Linux only
if [ "$OS" = "Linux" ]; then 
	execute
else
    echo "Trying to run on a un-supported environment [$OS]. \n👉Supported OS for dotnet agent is Linux. \n👉Please run the script on supported platform!"
fi
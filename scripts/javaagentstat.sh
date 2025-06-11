#!/bin/sh

output_file="/dev/stdout"

add_header() {
    echo "==================================================" > "$output_file"
    echo "                Java Agent Stats            "       >> "$output_file"
    echo "==================================================" >> "$output_file"
    echo "Generated on: $(date)" >> "$output_file"
    echo "==================================================" >> "$output_file"
    echo "" >> "$output_file"
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

# Check Java version
check_java_version() {
    logf "Java Stats"
    #@jai suppressing any output from the type -p java to /dev/null 
    if command -v java; then
        echo Found java executable in PATH >> "$output_file"
        _java=java
    elif [ -n "$JAVA_HOME" ] && [ -x "$JAVA_HOME/bin/java" ];  then
        echo Found java executable in JAVA_HOME >> "$output_file" 
        _java="$JAVA_HOME/bin/java"
    else
        echo "Java Not Found." >> "$output_file"
    fi

    if [[ "$_java" ]]; then
        details=$("$_java" -version 2>&1)
        version=$("$_java" -version 2>&1 | awk -F '"' '/version/ {print $2}')
        echo "$details" >> "$output_file"
        echo Version: "$version" >> "$output_file"
        if [[ "$version" > "1.7" ]]; then
            echo "Java version is more than 1.7" >> "$output_file"
        else         
            echo "Java version is less than or equal to 1.7" >> "$output_file"
        fi
    fi
}

# List Java processes
gather_java_processes() {
    logf "Java Processes"
    ps -aef | grep java| grep -v grep | awk '!/sh -c #!\/bin\/sh/' >> "$output_file"
}

# List Java agent processes
gather_java_agent_processes() {
    logf "Java Agent Processes:"
    ps -aef | grep javaagent: |grep -v grep| awk '!/sh -c #!\/bin\/sh/' >> "$output_file"
}

# Parse the Java argument and collect JAVA_AGENT_INSTALL_DIR
gather_java_agent_install_dir() {
    logf "JAVA_AGENT_INSTALL_DIR Details"
    
    ps -aef | grep javaagent | awk '{for(i=1;i<=NF;i++) if ($i ~ /^-javaagent/) print $i}' | cut -d':' -f2 | xargs dirname | while read dir; do
        if [ -n "$dir" ]; then
            echo "$dir" >> "$output_file"
        else
            logf "Could not determine any Java Agent Install Directory. Please check if the Java agent is running."
            break
        fi
    done
}

# Check permissions in each JAVA_AGENT_INSTALL_DIR
check_java_agent_install_dir_perm() {
    ps -aef | grep javaagent | awk '{for(i=1;i<=NF;i++) if ($i ~ /^-javaagent/) print $i}' | cut -d':' -f2 | xargs dirname | while read AGENT_DIR; do
        if [ -d "$AGENT_DIR" ]; then
            logf "Listing permissions in the Java agent directory: $AGENT_DIR"
            ls -laR "$AGENT_DIR" >> "$output_file"
        else
            logf "Directory does not exist: $AGENT_DIR"
        fi
    done
}

# Function to orchestrate the script execution
execute() {
    add_header # Dont remove

    check_java_version
    gather_java_processes
    gather_java_agent_processes
    gather_java_agent_install_dir
    check_java_agent_install_dir_perm

    echo "\nEnd" >> "$output_file"
}

execute

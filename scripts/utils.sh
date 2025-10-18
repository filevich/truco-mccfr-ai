#!/bin/bash

# Utility functions for tournament-based training on Slurm clusters
# This script provides common utility functions for:
# - Tournament server synchronization and connection management
# - Service health checking

# ============================================
# Common Utility Functions
# ============================================

# Function to find an available port
find_available_port() {
    local start_port=${1:-8080}
    local max_attempts=100
    local port=$start_port

    for i in $(seq 0 $max_attempts); do
        port=$((start_port + i))

        # Check if port is in valid range
        if [ $port -gt 65535 ]; then
            echo "ERROR: Exceeded maximum port number" >&2
            return 1
        fi

        # Try to connect to the port to see if it's in use
        if ! timeout 1 bash -c "echo >/dev/tcp/localhost/$port" 2>/dev/null; then
            # Port is not in use, return it
            echo $port
            return 0
        fi

        # Port is in use, try next one
        echo "Port $port is in use, trying next..." >&2
    done

    echo "ERROR: Could not find an available port after $max_attempts attempts" >&2
    return 1
}

# Function to wait for a TCP service to be ready
wait_for_service() {
    local address=$1
    local timeout=${2:-600}  # Default 10 minutes timeout
    local interval=5         # Check every 5 seconds
    local elapsed=0

    # Parse hostname and port
    local host=$(echo "$address" | cut -d: -f1)
    local port=$(echo "$address" | cut -d: -f2)

    echo "Waiting for TCP service at ${address} to be ready..."

    while [ $elapsed -lt $timeout ]; do
        # Use bash's built-in TCP connectivity test
        if timeout 3 bash -c "echo >/dev/tcp/$host/$port" >/dev/null 2>&1; then
            echo "Service at ${address} is ready! (waited ${elapsed} seconds)"
            return 0
        fi

        echo "Service not ready yet, waited ${elapsed}/${timeout} seconds..."
        sleep $interval
        elapsed=$((elapsed + interval))
    done

    echo "ERROR: Service at ${address} did not become ready within ${timeout} seconds"
    return 1
}

# Function to cleanup when training finishes (for shared tournament mode)
# Can be customized via environment variables:
#   TRAINING_JOB_PATTERN: Job name pattern to match training jobs (default: current job name)
#   SLURM_USER: Username for slurm queries (default: current user from whoami)
#   TOURNAMENT_JOB_NAME: Name of tournament server job (default: "gotournament")
cleanup_training() {
    echo "Training job finishing, checking if we should cleanup tournament..."

    # Use environment variables or defaults
    local user="${SLURM_USER:-$(whoami)}"
    local training_pattern="${TRAINING_JOB_PATTERN:-${SLURM_JOB_NAME}}"
    local tournament_job="${TOURNAMENT_JOB_NAME:-gotournament}"

    # Count active training jobs with matching pattern (excluding this one)
    ACTIVE_TRAINING=$(squeue -u "$user" -n "$training_pattern" --noheader -o "%i" | grep -v "^${SLURM_JOB_ID}$" | wc -l)

    echo "Active training jobs (excluding this one): $ACTIVE_TRAINING"
    echo "Checking for jobs matching: user=$user, name=$training_pattern"

    if [ "$ACTIVE_TRAINING" -eq 0 ]; then
        echo "This is the last training job, shutting down tournament server..."
        scancel -u "$user" -n "$tournament_job"
        # Clean up shared directory
        rm -rf $HOME/shared/tournament
        echo "Tournament server shutdown complete"
    else
        echo "Other training jobs still running, leaving tournament server active"
    fi
}

# ============================================
# Tournament Synchronization Functions
# ============================================

# Shared directories
TOURNAMENT_SHARED_DIR=$HOME/shared/tournament
LOCK_DIR=$TOURNAMENT_SHARED_DIR/.lock
CONNECTION_FILE=$TOURNAMENT_SHARED_DIR/connection.txt

# Timeouts and intervals
LOCK_TIMEOUT=300       # 5 minutes to acquire lock
LOCK_RETRY_INTERVAL=5  # Check every 5 seconds
CONNECTION_TIMEOUT=600 # 10 minutes to wait for connection info
CONNECTION_CHECK_INTERVAL=5

# Function to check if tournament server is alive
tournament_is_alive() {
    local connection_info="$1"

    if [ -z "$connection_info" ]; then
        return 1
    fi

    # Parse hostname and port
    local host=$(echo "$connection_info" | cut -d: -f1)
    local port=$(echo "$connection_info" | cut -d: -f2)

    # Try to connect to the tournament server
    if timeout 3 bash -c "echo >/dev/tcp/$host/$port" >/dev/null 2>&1; then
        return 0
    else
        return 1
    fi
}

# Function to check if tournament job exists
tournament_job_exists() {
    local user="${SLURM_USER:-$(whoami)}"
    local tournament_job="${TOURNAMENT_JOB_NAME:-gotournament}"
    local job_count=$(squeue -u "$user" -n "$tournament_job" --noheader | wc -l)
    echo "DEBUG: Found $job_count tournament job(s) in queue" >&2
    if [ "$job_count" -gt 0 ]; then
        return 0
    else
        return 1
    fi
}

# Function to acquire the lock
acquire_lock() {
    local elapsed=0

    echo "Attempting to acquire lock at: $LOCK_DIR" >&2

    while [ $elapsed -lt $LOCK_TIMEOUT ]; do
        if mkdir "$LOCK_DIR" 2>/dev/null; then
            echo "Lock acquired successfully" >&2
            return 0
        fi

        # Check who owns the lock (for debugging)
        if [ -d "$LOCK_DIR" ]; then
            echo "Lock busy (owned by another process), waiting... (${elapsed}/${LOCK_TIMEOUT}s)" >&2
            echo "DEBUG: Lock directory info:" >&2
            ls -ld "$LOCK_DIR" 2>&1 >&2 || true
        fi

        sleep $LOCK_RETRY_INTERVAL
        elapsed=$((elapsed + LOCK_RETRY_INTERVAL))
    done

    echo "ERROR: Failed to acquire lock after ${LOCK_TIMEOUT} seconds" >&2
    echo "ERROR: Lock still exists at: $LOCK_DIR" >&2
    echo "DEBUG: Stale lock info:" >&2
    ls -ld "$LOCK_DIR" 2>&1 >&2 || echo "Lock disappeared" >&2
    return 1
}

# Function to release the lock
release_lock() {
    rmdir "$LOCK_DIR" 2>/dev/null || true
    echo "Lock released" >&2
}

# Function to launch tournament job
# Uses environment variable TOURNAMENT_SCRIPT_PATH or defaults to a sensible location
launch_tournament_job() {
    echo "Launching tournament server job..." >&2

    # Use environment variable or try to find tournament.slurm in common locations
    local tournament_script="${TOURNAMENT_SCRIPT_PATH}"

    if [ -z "$tournament_script" ]; then
        # Try common locations
        if [ -f "$HOME/Workspace/truco-ai/truco-mccfr-ai/scripts/tournament.slurm" ]; then
            tournament_script="$HOME/Workspace/truco-ai/truco-mccfr-ai/scripts/tournament.slurm"
        elif [ -f "$HOME/scripts/ismcts/tournament.slurm" ]; then
            tournament_script="$HOME/scripts/ismcts/tournament.slurm"
        else
            echo "ERROR: Tournament script not found. Set TOURNAMENT_SCRIPT_PATH environment variable." >&2
            return 1
        fi
    fi

    echo "Looking for tournament script at: $tournament_script" >&2

    if [ ! -f "$tournament_script" ]; then
        echo "ERROR: Tournament script not found at: $tournament_script" >&2
        return 1
    fi

    echo "Tournament script found, submitting job..." >&2

    # Submit the tournament job and capture both output and exit code
    local job_id
    job_id=$(sbatch --parsable "$tournament_script" 2>&1)
    local submit_status=$?

    if [ $submit_status -ne 0 ]; then
        echo "ERROR: Failed to submit tournament job (exit code: $submit_status)" >&2
        echo "ERROR: sbatch output: $job_id" >&2
        return 1
    fi

    echo "Tournament job submitted successfully with ID: $job_id" >&2

    # Verify the job was actually queued
    sleep 2
    if tournament_job_exists; then
        echo "Verified: Tournament job $job_id is in queue" >&2
    else
        echo "WARNING: Tournament job was submitted but not found in queue!" >&2
        echo "DEBUG: Checking squeue for job $job_id:" >&2
        squeue -j "$job_id" 2>&1 >&2 || echo "Job not found in squeue" >&2
    fi

    return 0
}

# Function to wait for connection file to appear with valid content
wait_for_connection_info() {
    local elapsed=0

    echo "Waiting for tournament connection info..." >&2

    while [ $elapsed -lt $CONNECTION_TIMEOUT ]; do
        if [ -f "$CONNECTION_FILE" ]; then
            local connection=$(cat "$CONNECTION_FILE" 2>/dev/null)

            if [ -n "$connection" ]; then
                # Check if tournament is actually alive
                if tournament_is_alive "$connection"; then
                    echo "Tournament connection established: $connection" >&2
                    echo "$connection"
                    return 0
                else
                    echo "Connection file exists but tournament not responding yet..." >&2
                fi
            fi
        fi

        sleep $CONNECTION_CHECK_INTERVAL
        elapsed=$((elapsed + CONNECTION_CHECK_INTERVAL))

        # Show progress every 30 seconds
        if [ $((elapsed % 30)) -eq 0 ]; then
            echo "Still waiting for tournament... (${elapsed}/${CONNECTION_TIMEOUT}s)" >&2
        fi
    done

    echo "ERROR: Tournament did not become available within ${CONNECTION_TIMEOUT} seconds" >&2
    return 1
}

# Main function: acquire tournament connection
# This is the primary entry point for training scripts
# Returns the tournament connection string (host:port) on success
acquire_tournament_connection() {
    # Ensure shared directory exists
    mkdir -p "$TOURNAMENT_SHARED_DIR"

    echo "Attempting to acquire tournament connection..." >&2

    # Fast path: Check if connection already exists and is alive
    if [ -f "$CONNECTION_FILE" ]; then
        local connection=$(cat "$CONNECTION_FILE" 2>/dev/null)
        if tournament_is_alive "$connection"; then
            echo "Reusing existing tournament connection: $connection" >&2
            echo "$connection"
            return 0
        else
            echo "Existing connection file is stale, will try to restart tournament" >&2
        fi
    fi

    # Slow path: Need to potentially start tournament
    echo "No active tournament found, attempting to start one..." >&2

    # Acquire lock for exclusive access to tournament startup
    if ! acquire_lock; then
        return 1
    fi

    # Double-check after acquiring lock (another job might have started it)
    if [ -f "$CONNECTION_FILE" ]; then
        local connection=$(cat "$CONNECTION_FILE" 2>/dev/null)
        if tournament_is_alive "$connection"; then
            echo "Tournament started by another job: $connection" >&2
            release_lock
            echo "$connection"
            return 0
        fi
    fi

    # Check if tournament job exists in queue
    if ! tournament_job_exists; then
        echo "No tournament job found, launching new one..." >&2
        if ! launch_tournament_job; then
            release_lock
            return 1
        fi
    else
        echo "Tournament job already exists in queue, waiting for it to start..." >&2
    fi

    # Wait for tournament to write connection info and become available
    local connection=$(wait_for_connection_info)
    local result=$?

    # Release lock
    release_lock

    if [ $result -eq 0 ]; then
        echo "$connection"
        return 0
    else
        return 1
    fi
}

# Export the main function for use by other scripts
export -f acquire_tournament_connection

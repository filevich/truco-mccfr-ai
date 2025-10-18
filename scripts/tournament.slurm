#!/bin/bash

#SBATCH --job-name=gotournament
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=1
#SBATCH --mem=8G
#SBATCH --time=5-00:00:00
#SBATCH --partition=normal
#SBATCH --qos=normal
#SBATCH --output=/clusteruy/home/juan.filevich/batches/out/tournament/%x.%j.out
#SBATCH --error=/clusteruy/home/juan.filevich/batches/out/tournament/%x.%j.out
#SBATCH --mail-type=BEGIN,FAIL,END
#SBATCH --mail-user=juan.filevich@fing.edu.uy

# Source the utility functions
source $HOME/Workspace/truco-ai/truco-mccfr-ai/scripts/utils.sh

# Tournament-specific cleanup function
cleanup_tournament() {
    echo "Cleaning up tournament server..."

    # Remove connection info
    rm -f $HOME/shared/tournament/connection.txt
    rm -f $HOME/shared/tournament/tournament.pid

    # Kill any baseline agents we launched
    if [ -f $HOME/shared/tournament/random_agent.pid ]; then
        kill $(cat $HOME/shared/tournament/random_agent.pid) 2>/dev/null || true
        rm -f $HOME/shared/tournament/random_agent.pid
    fi

    if [ -f $HOME/shared/tournament/simple_agent.pid ]; then
        kill $(cat $HOME/shared/tournament/simple_agent.pid) 2>/dev/null || true
        rm -f $HOME/shared/tournament/simple_agent.pid
    fi

    echo "Tournament cleanup complete"
}

# Set trap to ensure cleanup on exit
trap cleanup_tournament EXIT

# Print start info
printf "Tournament server starting: $(date)\n"
echo "Job ID: $SLURM_JOB_ID"
echo "Node: $(hostname)"

# Ensure shared directory exists
mkdir -p $HOME/shared/tournament
mkdir -p $HOME/batches/out/tournament

# Get the hostname and find an available port
HOSTNAME=$(hostname)

echo "Finding available port for tournament server..."
TOURNAMENT_PORT=$(find_available_port 8080)
if [ $? -ne 0 ]; then
    echo "Failed to find available port, exiting..."
    exit 1
fi

TOURNAMENT_ADDR="${HOSTNAME}:${TOURNAMENT_PORT}"

echo "Tournament will be available at: ${TOURNAMENT_ADDR}"

# Write connection info BEFORE launching server
echo "${TOURNAMENT_ADDR}" > $HOME/shared/tournament/connection.txt
echo "Connection info written to shared file"

# Launch the Tournament Master Service
echo "Launching Tournament Master Service..."
cd ~/Workspace/truco-ai/truco-tournament/ || exit

t1k22=$HOME/data/t1k22.json \
model=$HOME/checkpoints/go-mccfr/pruned_esvmccfr_d70h0m_D70h0m_t288652014_p1_a3_2402151230.model \
go run cmd/tournament/main.go \
    -port=$TOURNAMENT_PORT \
    -level=info &> $HOME/batches/out/tournament/tournament_${SLURM_JOB_ID}.log &

TOURNAMENT_PID=$!
echo $TOURNAMENT_PID > $HOME/shared/tournament/tournament.pid
echo "Tournament server launched with PID: $TOURNAMENT_PID"

# Wait a bit for tournament to start
echo "Waiting for tournament server to initialize..."
sleep 15

# Verify the tournament is still running
if ! kill -0 $TOURNAMENT_PID 2>/dev/null; then
    echo "ERROR: Tournament server failed to start!"
    exit 1
fi

# Launch the Random Agent
echo "Launching Random Agent..."
cd ~/Workspace/truco-ai/truco-tournament/ || exit
go run cmd/cli-example/main.go \
    -active=false \
    -loop=true \
    -policy=random \
    -name=rand \
    -tournament_addr=${TOURNAMENT_ADDR} \
    -level=warn &> $HOME/batches/out/tournament/random_agent_${SLURM_JOB_ID}.log &
echo $! > $HOME/shared/tournament/random_agent.pid

# Launch the Simple Agent
echo "Launching Simple Agent..."
cd ~/Workspace/truco-ai/truco-tournament/ || exit
go run cmd/cli-example/main.go \
    -active=false \
    -loop=true \
    -policy=simple \
    -name=simpl \
    -tournament_addr=${TOURNAMENT_ADDR} \
    -level=warn &> $HOME/batches/out/tournament/simple_agent_${SLURM_JOB_ID}.log &
echo $! > $HOME/shared/tournament/simple_agent.pid

echo "Tournament services launched and ready at ${TOURNAMENT_ADDR}"
echo "Baseline agents running"
echo "Tournament server will run until time limit or manually cancelled"

# Keep the job alive by waiting for the tournament process
wait $TOURNAMENT_PID
TOURNAMENT_EXIT_CODE=$?

echo "Tournament server exited with code: $TOURNAMENT_EXIT_CODE"
echo "Tournament ended: $(date)"

# Cleanup will happen automatically via trap

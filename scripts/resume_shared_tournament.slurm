#!/bin/bash

#SBATCH --job-name=defaultname
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=1
#SBATCH --mem=60G
#SBATCH --time=5-00:00:00
#SBATCH --partition=normal
#SBATCH --qos=normal
#SBATCH --output=/clusteruy/home/juan.filevich/stdout/gotruco-tab-ai/%x.%j.out
#SBATCH --error=/clusteruy/home/juan.filevich/stdout/gotruco-tab-ai/%x.%j.out
#SBATCH --mail-type=BEGIN,END,FAIL
#SBATCH --mail-user=juan.filevich@fing.edu.uy

# Source the utility functions
source $HOME/scripts/gotruco-tab-ai/utils.sh

# Set trap to cleanup on exit
trap cleanup_training EXIT

hr (){
    printf '%*s\n' 80 | tr ' ' '-'
}

cat $0

hr
free -mh
lscpu | head -n25

hr

echo "args"
echo $@

hr

# Defaults for tournament configuration
if [ -z "$TOURNAMENT_N" ]; then
    TOURNAMENT_N=1000
fi

echo "TOURNAMENT_N: $TOURNAMENT_N"
echo "Job ID: $SLURM_JOB_ID"

hr

# Acquire tournament connection (this will start tournament if needed)
echo "=========================================="
echo "Acquiring tournament server connection..."
echo "=========================================="
TOURNAMENT_ADDR=$(acquire_tournament_connection)

if [ $? -ne 0 ] || [ -z "$TOURNAMENT_ADDR" ]; then
    echo "ERROR: Failed to acquire tournament connection"
    exit 1
fi

echo "Successfully connected to tournament at: $TOURNAMENT_ADDR"
echo ""

hr

# special arg
fmt=${10}
if [ "$fmt" = "auto" ]; then
    # Assign $SLURM_JOB_ID to $fmt
    fmt="$SLURM_JOB_ID.model"
fi

# Use the connection file path for dynamic address resolution
TOURNAMENT_CONNECTION_FILE=$HOME/shared/tournament/connection.txt

# Start training in background so we can monitor tournament health
echo "Resuming training process..."
cd $HOME/Workspace/truco-ai/truco-mccfr-ai
go run cmd/train/main.go \
    -model=${1} \
    -threads=${2} \
    -dir=${3} \
    -run=${4} \
    -prunning=${5} \
    -prunning_prob=${6} \
    -save_every=${7} \
    -eval_every=${8} \
    -prefix=${9} \
    -fmt=${fmt} \
    --tournament-addr=$TOURNAMENT_CONNECTION_FILE \
    --tournament-name=train_${SLURM_JOB_ID} \
    --tournament-n=$TOURNAMENT_N &

TRAIN_PID=$!
echo "Training started with PID: $TRAIN_PID"

# Monitor tournament health while training runs (checks every hour by default)
# This will automatically restart the tournament if it dies during training
monitor_tournament_health $TRAIN_PID 3600

# Wait for training to complete
wait $TRAIN_PID
TRAIN_EXIT_CODE=$?

hr
echo "Training finished with exit code: $TRAIN_EXIT_CODE"
echo "done $(date)"

# Cleanup will happen automatically via trap

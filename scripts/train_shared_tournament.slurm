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

# Parse mode (first argument: "train" or "resume")
MODE=$1
shift  # Remove mode from positional arguments

if [ "$MODE" != "train" ] && [ "$MODE" != "resume" ]; then
    echo "ERROR: First argument must be 'train' or 'resume'" >&2
    echo "Usage (train):  sbatch $0 train <p> <trainer> <hash> <info> <abs> <threads> <dir> <run> <prunning> <prunning_prob> <save_every> <eval_every> <prefix> <fmt>" >&2
    echo "Usage (resume): sbatch $0 resume <model> <threads> <dir> <run> <prunning> <prunning_prob> <save_every> <eval_every> <prefix> <fmt>" >&2
    exit 1
fi

echo "MODE: $MODE"

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

# Use the connection file path for dynamic address resolution
TOURNAMENT_CONNECTION_FILE=$HOME/shared/tournament/connection.txt

# Start training in background so we can monitor tournament health
if [ "$MODE" = "train" ]; then
    echo "Starting new training process..."

    # special arg for train mode
    fmt=${14}
    if [ "$fmt" = "auto" ]; then
        fmt="$SLURM_JOB_ID.model"
    fi

    cd $HOME/Workspace/truco-ai/truco-mccfr-ai
    go run cmd/train/main.go \
        -p=${1} \
        -trainer=${2} \
        -hash=${3} \
        -info=${4} \
        -abs=${5} \
        -threads=${6} \
        -dir=${7} \
        -run=${8} \
        -prunning=${9} \
        -prunning_prob=${10} \
        -save_every=${11} \
        -eval_every=${12} \
        -prefix=${13} \
        -fmt=$fmt \
        --tournament-addr=$TOURNAMENT_CONNECTION_FILE \
        --tournament-name=train_${SLURM_JOB_ID} \
        --tournament-n=$TOURNAMENT_N &
elif [ "$MODE" = "resume" ]; then
    echo "Resuming training from model checkpoint..."

    # special arg for resume mode
    fmt=${10}
    if [ "$fmt" = "auto" ]; then
        fmt="$SLURM_JOB_ID.model"
    fi

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
        -fmt=$fmt \
        --tournament-addr=$TOURNAMENT_CONNECTION_FILE \
        --tournament-name=train_${SLURM_JOB_ID} \
        --tournament-n=$TOURNAMENT_N &
fi

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

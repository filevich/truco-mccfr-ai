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
fmt=${14}
if [ "$fmt" = "auto" ]; then
    # Assign $SLURM_JOB_ID to $fmt
    fmt="$SLURM_JOB_ID.model"
fi

# Use the connection file path for dynamic address resolution
TOURNAMENT_CONNECTION_FILE=$HOME/shared/tournament/connection.txt

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
    --tournament-n=$TOURNAMENT_N

hr
echo "done $(date)"

# Cleanup will happen automatically via trap

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

# special arg
fmt=${14}
if [ "$fmt" = "auto" ]; then
    # Assign $SLURM_JOB_ID to $fmt
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
    -fmt=$fmt

hr
echo "done $(date)"
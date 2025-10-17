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
fmt=${10}
if [ "$fmt" = "auto" ]; then
    # Assign $SLURM_JOB_ID to $fmt
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
    -fmt=${fmt}

hr
echo "done $(date)"
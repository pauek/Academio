#!/bin/bash
echo "Setting environment"
source $HOME/go/src/Academio/sh/server/env.sh
cd $ACADEMIO_ROOT/webapp
echo "Starting new process"
nohup webapp &>> $ACADEMIO_ROOT/log &
echo "Done"

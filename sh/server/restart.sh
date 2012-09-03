#!/bin/bash
echo "Setting environment"
source $HOME/go/src/Academio/sh/server/env.sh
echo -n "Recompiling... "
go install Academio/webapp
echo "done"
cd $ACADEMIO_ROOT/webapp
echo "Killing old process"
pkill -9 webapp
echo "Starting new process"
nohup webapp &>> $ACADEMIO_ROOT/log &
echo "Done"

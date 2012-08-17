#!/bin/bash
echo "Setting environment"
source $HOME/go/src/Academio/server-env.sh
echo -n "Recompiling... "
go install Academio/webapp
echo "done"
cd $ACADEMIO_ROOT/webapp
sudo -s <<EOF
echo "Killing old process"
pkill -9 webapp
echo "Starting new process"
ACADEMIO_PATH=$ACADEMIO_PATH \
ACADEMIO_ROOT=$ACADEMIO_ROOT \
  nohup $(which webapp) -port=80 &>> $ACADEMIO_ROOT/log &
EOF
echo "Done"

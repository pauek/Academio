#!/bin/bash
rsync -av $ACADEMIO_ROOT/ academio:go/src/Academio
ssh academio go/src/Academio/server-restart.sh

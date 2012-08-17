#!/bin/sh
rsync -av $ACADEMIO_PATH/ academio:Academio
rsync -av $ACADEMIO_ROOT/ academio:go/src/Academio

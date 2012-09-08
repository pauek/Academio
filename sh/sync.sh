#!/bin/sh
rsync -av $ACADEMIO_PATH/ academio:Academio/Content
rsync -av $ACADEMIO_ROOT/ academio:go/src/Academio

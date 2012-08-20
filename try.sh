#!/bin/bash
ACADEMIO_PATH=$ACADEMIO_PATH \
ACADEMIO_ROOT=$ACADEMIO_ROOT \
  su -c "$(which webapp) -ssl -port=443" 

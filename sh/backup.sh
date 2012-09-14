#!/bin/bash
TARGET=/media/wd/Academio
rsync -av /pub/Academio/Videos/ $TARGET/Videos
rsync -av /pub/Academio/Problems/ $TARGET/Problems
rsync -av $HOME/Academio/Content/ $TARGET/Content

#!/bin/sh

export LD_LIBRARY_PATH=/mnt/SDCARD/RetroArch/lib:/usr/trimui/lib:$LD_LIBRARY_PATH
export PATH=/usr/trimui/bin:$PATH

./retroarch -c retroarch.cfg --menu -v

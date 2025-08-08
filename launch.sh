#!/bin/sh
CANNOLI_DIR="$(dirname "$0")"
cd "$CANNOLI_DIR" || exit 1

export LD_LIBRARY_PATH=/usr/trimui/lib:$CANNOLI_DIR/lib:$LD_LIBRARY_PATH
export PATH=/usr/trimui/bin:$PATH

echo 0 > /sys/class/led_anim/max_scale
if [ "$TRIMUI_MODEL" = "Trimui Brick" ]; then
	echo 0 > /sys/class/led_anim/max_scale_lr
	echo 0 > /sys/class/led_anim/max_scale_f1f2
fi

trimui_inputd &

export HOME=/mnt/SDCARD

while true; do
  ./cannoliOS
done

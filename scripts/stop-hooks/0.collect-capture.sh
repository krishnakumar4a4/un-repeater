#!/bin/bash

echo "Collecting capture"
SESSION_FOLDER_PATH=$1
echo "session folder path: ${SESSION_FOLDER_PATH}"


echo "killing adb shells"
ps -efl|grep "adb -s"|tr -s ' '|cut -d' ' -f3|xargs kill
echo "killed adb shells"

adb devices -l|grep -w "device" | while read line;
do 
    echo "Generating bug report -> $line"
    DEVICE_ID=`echo $line|cut -d' ' -f1`
    DEVICE_MODEL=`echo $line| cut -d' ' -f 5|cut -d':' -f2`
    BUGREPORT="${SESSION_FOLDER_PATH}/bugreport-${DEVICE_MODEL}.zip"
    adb -s ${DEVICE_ID} bugreport "${BUGREPORT}"
    BTSNOOZ=scripts/utilities/btsnooz.py scripts/utilities/extractbtsnoop.sh ${SESSION_FOLDER_PATH} ${BUGREPORT}
done

echo "bug report/s and BLE snoops generated"

echo "killing adb shells"
ps -efl|grep "adb -s"|tr -s ' '|cut -d' ' -f3|xargs kill
echo "killed adb shells"

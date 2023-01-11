#!/bin/bash

# This is a script that will allow you to create bluetooth snoop. 
# This script will handle the unzipping if it is a zip file, or will convert the btsnooz
# directly in the case of a plain text file for use with wireshark.

ROOT_DIR=$1
BTSNOOZ="${BTSNOOZ:-btsnooz.py}"

if ! hash btsnooz.py 2>/dev/null;
then
    echo "Please make sure btsnooz.py is in your path before running."
    exit 2;
fi

if [ $# -eq 0 ];
then
    echo "Usage: $0 path-to-root-dir bugreport(.txt|.zip)"
    exit 3;
fi

BUGREPORT="$2"
FILENAME="$(basename ${BUGREPORT})"
TMPDIR=$(mktemp -d "${ROOT_DIR}/extractbtsnooz_XXXXX")
LOGFILE="${ROOT_DIR}/${FILENAME%.*}.btsnooz"

trap ctrl_c INT
function ctrl_c() {
    rm -rf "${TMPDIR}"
}

if [ ! -f "${BUGREPORT}" ];
then
    echo "File ${BUGREPORT} does not exist."
    exit 4;
fi

if [ ! -d "${TMPDIR}" ];
then
    echo "Unable to create temp. dir (${TMPDIR}) :("
    exit 5;
fi

if [ "${BUGREPORT: -4}" == ".zip" ];
then
    unzip "${BUGREPORT}" -d "${TMPDIR}"
    BUGREPORT="${TMPDIR}/`ls ${TMPDIR}|grep ${FILENAME%.*}.*\.txt`"
    if [ "${BUGREPORT: -4}" != ".txt" ];then
        FILENAME="bugreport"
        BUGREPORT="${TMPDIR}/`ls ${TMPDIR}|grep ${FILENAME%.*}.*\.txt`"
        echo "bugreport not found with given name, trying with default name ${BUGREPORT}"
    fi
fi

if [ -f "${BUGREPORT}" ];
then
    python3 ${BTSNOOZ} "${BUGREPORT}" > "${LOGFILE}"
    if [ ! $? -eq 0 ];
    then
        echo "Could not extract btsnooz data from ${BUGREPORT}."
        rm -rf "${TMPDIR}"
        exit 6;
    fi
    echo "bt snoop generated in ${LOGFILE}"
else
    echo "Looks like there is no plain text bugreport (${BUGREPORT})?"
fi

rm -rf "${TMPDIR}"

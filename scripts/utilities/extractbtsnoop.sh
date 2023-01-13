#!/bin/bash

# This is a script that will allow you to create bluetooth snoop. 
# This script will handle the unzipping if it is a zip file, or will convert the btsnooz
# directly in the case of a plain text file for use with wireshark.

ROOT_DIR=$1
BTSNOOZ="${BTSNOOZ:-btsnooz.py}"

if ! hash ${BTSNOOZ} 2>/dev/null;
then
    echo "ERR: Please make sure btsnooz.py is in your path before running."
    exit 2;
fi

if [ $# -lt 2 ];
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
    echo "ERR: File ${BUGREPORT} does not exist."
    exit 4;
fi

if [ ! -d "${TMPDIR}" ];
then
    echo "ERR: Unable to create temp. dir (${TMPDIR}) :("
    exit 5;
fi

if [ "${BUGREPORT: -4}" == ".zip" ];
then
    unzip "${BUGREPORT}" -d "${TMPDIR}"
    BUGREPORT="${TMPDIR}/`ls ${TMPDIR}|grep ${FILENAME%.*}.*\.txt`"
    if [ "${BUGREPORT: -4}" != ".txt" ];then
        echo "WARN: bugreport not found with given name, trying again with default name"
        FILENAME="bugreport"
        BUGREPORT="${TMPDIR}/`ls ${TMPDIR}|grep ${FILENAME%.*}.*\.txt`"
    fi
fi

if [ -f "${BUGREPORT}" ];
then
    python3 ${BTSNOOZ} "${BUGREPORT}" > "${LOGFILE}"
    if [ ! $? -eq 0 ];
    then
        echo "WARN: Could not extract btsnooz data from ${BUGREPORT}."
        echo "INFO: Attempting manual search for BLE capture on the zip"
        CURF_FILE=`find ${TMPDIR} -iname "*.curf"|head -n 1`
        if [ ! -z "${CURF_FILE}" ];then
            cp -f ${CURF_FILE} "${ROOT_DIR}/"
            echo "INFO: Found BLE capture on manual search at ${CURF_FILE}"
            rm -rf "${TMPDIR}"
            exit 0
        else
            HCI_LOG=`find ${TMPDIR} -iname "*btsnoop_hci.log"|head -n 1`
            if [ ! -z "${HCI_LOG}" ];then
                cp -f ${HCI_LOG} "${ROOT_DIR}/"
                echo "INFO: Found BLE capture on manual search at ${HCI_LOG}"
                rm -rf "${TMPDIR}"
                exit 0
            else
                CFA_COUNT=`find ${TMPDIR} -iname "*.cfa"|wc -l`
                if [ $CFA_COUNT -ge 1 ]; then
                    while read CFA_FILE;do
                        if [ ! -z "${CFA_FILE}" ];then
                            cp -f ${CFA_FILE} "${ROOT_DIR}/"
                            echo "INFO: Found BLE capture on manual search at ${CFA_FILE}"
                        fi
                    done <<< `find ${TMPDIR} -iname "*.cfa"`
                    rm -rf "${TMPDIR}"
                    exit 0
                fi
            fi
            echo "ERR: Attempt to manually find BLE capture also failed"
        fi
        rm -rf "${TMPDIR}"
        exit 6;
    fi
    echo "INFO: BT snoop generated in ${LOGFILE}"
else
    echo "ERR: Looks like there is no plain text bugreport (${BUGREPORT})?"
fi

rm -rf "${TMPDIR}"


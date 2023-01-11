#!/bin/bash
SESSION_FOLDER_PATH=$1
echo "session folder path: ${SESSION_FOLDER_PATH}"

echo "open file to get session description"
echo "DATE: `date`" > ${SESSION_FOLDER_PATH}/session-description.md
echo "TODO: Describe the session here" >> ${SESSION_FOLDER_PATH}/session-description.md

# Open file with default editor
open -t ${SESSION_FOLDER_PATH}/session-description.md

echo "File opened"

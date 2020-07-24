#!/bin/bash

# Start the first process
/usr/local/zevenet/app/nftlb/sbin/nftlb -L #LOGSOUTPUT# -l #LOGSLEVEL# -k #KEY# -d -m #MASQUERADEMARK#
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start nftlb: $status"
  exit $status
fi

# waiting a grace time
sleep 3

# Start the second process
/goclient #KEY# #CLIENTCFG# &
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start GO client: $status"
  exit $status
fi

# Naive check runs checks once a minute to see if either of the processes exited.
# This illustrates part of the heavy lifting you need to do if you want to run
# more than one service in a container. The container exits with an error
# if it detects that either of the processes has exited.
# Otherwise it loops forever, waking up every 60 seconds

while sleep #DAEMONCHECKTIMEOUT#; do
  ps aux |grep nftlb |grep -q -v grep
  PROCESS_NFTLB_STATUS=$?
  if [ $PROCESS_NFTLB_STATUS -ne 0 ]; then
    echo "The nftlb process exited with error."
  fi

  ps aux |grep goclient |grep -q -v grep
  PROCESS_GO_STATUS=$?
  if [ $PROCESS_GO_STATUS -ne 0 ]; then
    echo "The GO client exited with error."
  fi

  if [ $PROCESS_NFTLB_STATUS -ne 0 -o $PROCESS_GO_STATUS -ne 0 ]; then
    exit 1
  fi
done


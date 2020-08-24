#!/usr/bin/env bash

# Load environment values
source .env

# Start the first process, nftlb, with devel flags
/usr/local/zevenet/app/nftlb/sbin/nftlb -L "$NFTLB_LOGS_OUTPUT" -l "$NFTLB_LOGS_LEVEL" -k "$NFTLB_KEY" -d -m "$NFTLB_MASQUERADE_MARK"
status=$?
if [ $status -ne 0 ]; then
  # If that fails, start nftlb without devel flags
  echo "Failed to start nftlb with devel flags, trying without them..."
  /usr/local/zevenet/app/nftlb/sbin/nftlb -l "$NFTLB_LOGS_LEVEL" -k "$NFTLB_KEY" -d
  status=$?
  if [ $status -ne 0 ]; then
    echo "Failed to start nftlb: $status"
    exit $status
  fi
fi

# Wait a grace time
sleep 3

# Start the second process, goclient
/goclient &
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

while sleep "$HEALTH_CHECK_TIMEOUT"; do
  ps aux | grep nftlb | grep -q -v grep
  PROCESS_NFTLB_STATUS=$?
  if [ $PROCESS_NFTLB_STATUS -ne 0 ]; then
    echo "The nftlb process exited with error."
  fi

  ps aux | grep goclient | grep -q -v grep
  PROCESS_GO_STATUS=$?
  if [ $PROCESS_GO_STATUS -ne 0 ]; then
    echo "The GO client exited with error."
  fi

  if [ $PROCESS_NFTLB_STATUS -ne 0 -o $PROCESS_GO_STATUS -ne 0 ]; then
    exit 1
  fi
done

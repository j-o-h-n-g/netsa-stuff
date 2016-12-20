#!/bin/bash

# A helper script to ease the management of rwsender regexps
# Called from /etc/sysconfig/rwsender.conf via RECEIVER_CLIENTS=`/etc/silk/rwsender.sh`
# The zzzzzzz are just a pattern which should never match as double quotes
# included in first alternation, but |'s interpreted by shell if " are missing
# Assumes all probe names are of the form P(ip).  eg P192.168.1.0

DIRECTORY=/etc/silk/rwsender

for FILE in ${DIRECTORY}/*
do
  RECEIVER=${FILE##*/}
  echo -n $RECEIVER
  echo -n " \"zzzzzzzz|P"
  cut -d# -f1 ${FILE} | tr -d '\t ' | tr '\n' '|' | sed 's/|/\.|P/g' | sed 's/\.
/\\./g'
  echo "zzzzzzzzzz\""
done

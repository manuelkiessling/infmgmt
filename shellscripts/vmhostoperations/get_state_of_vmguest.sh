#!/bin/bash

VMHOST_DNSNAME=$1
INDEX_OF_VMGUEST=$2

ssh -i ~/.ssh/infmgmt.id_rsa \
  root@$VMHOST_DNSNAME \
  virsh list --all \
  | tail --lines=+$((3 + $INDEX_OF_VMGUEST)) \
  | head --lines=1 \
  | sed "s/ \+/ /g" \
  | cut -d" " -f4-

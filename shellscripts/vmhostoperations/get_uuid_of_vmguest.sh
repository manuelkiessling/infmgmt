#!/bin/bash

VMHOST_DNSNAME=$1
VMGUEST_NAME=$2

ssh -i ~/.ssh/infmgmt.id_rsa \
  root@$VMHOST_DNSNAME \
  virsh dumpxml "$VMGUEST_NAME" \
  | grep "uuid" \
  | cut --bytes=9-44

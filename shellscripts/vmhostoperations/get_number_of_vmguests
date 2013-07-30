#!/bin/bash

VMHOST_DNSNAME=$1

ssh -i ~/.ssh/infmgmt.id_rsa \
  root@$VMHOST_DNSNAME virsh list --all \
  | tail --lines=+3 \
  | head --lines=-1 \
  | wc -l

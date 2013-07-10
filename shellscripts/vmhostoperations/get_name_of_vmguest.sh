#!/bin/bash

VMHOST_DNSNAME=$1
NUMBER_OF_MACHINE=$2

# TODO: Calculate tail line number from NUMBER_OF_MACHINE

ssh -i ~/.ssh/infmgmt.id_rsa root@$VMHOST_DNSNAME virsh list --all | tail --lines=+3+$NUMBER_OF_MACHINE | head --lines=1 | sed 's/ \\+/ /g' | cut -d' ' -f3

#!/bin/bash

if [ "$#" -lt 1 ] ; then
  echo "Usage: $0 <Security Profile>"
  echo "Where <Security Profile> is: 1, 2, or 3."
  exit 1
fi

SP=$1

echo "Patching the CSMS to disable load balancer"
patch -p1 -i config/everest/maeve-csms-no-lb.patch

if [[ $SP == 1 ]]; then
  echo "Patching the CSMS to disable WSS"
  patch -p1 -i config/everest/maeve-csms-no-wss.patch
else 
  echo "Patching the CSMS to enable EVerest organization"
  patch -p1 -i config/everest/maeve-csms-everest-org.patch
      
  echo "Patching the CSMS to enable local mo root"
  patch -p1 -i config/everest/maeve-csms-local-mo-root.patch
      
  echo "Patching the CSMS to enable local mo root"
  patch -p1 -i config/everest/maeve-csms-ignore-ocsp.patch
fi

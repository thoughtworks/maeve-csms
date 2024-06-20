#!/bin/bash

if [ "$#" -lt 1 ] ; then
  echo "Usage: $0 <Security Profile>"
  echo "Where <Security Profile> is: 1, 2, or 3."
  exit 1
fi

SP=$1

if [[ $SP == 2 || $SP == 3 ]]; then
  echo "Patching the CSMS to enable EVerest organization"
  patch -p1 -i config/everest/maeve-csms-everest-org.patch
      
  echo "Patching the CSMS to enable local mo root"
  patch -p1 -i config/everest/maeve-csms-local-mo-root.patch
      
  echo "Patching the CSMS to ignore OCSP"
  patch -p1 -i config/everest/maeve-csms-ignore-ocsp.patch
fi

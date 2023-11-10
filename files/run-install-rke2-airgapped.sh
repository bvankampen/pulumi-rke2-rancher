#!/usr/bin/bash
RKE2_ACTIVE=$(systemctl is-active rke2-server)
if [[ $RKE2_ACTIVE != "active" ]]; then
  export INSTALL_RKE2_ARTIFACT_PATH=/opt/rke2/install/
  sh /opt/rke2/install/install.sh
  systemctl enable rke2-server
  systemctl start rke2-server
fi

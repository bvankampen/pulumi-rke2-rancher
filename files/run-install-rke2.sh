#!/usr/bin/bash
RKE2_ACTIVE=$(systemctl is-active rke2-server)
if [[ $RKE2_ACTIVE != "active" ]]; then
  curl -sfL https://get.rke2.io | sh -
  systemctl enable rke2-server
  systemctl start rke2-server
fi

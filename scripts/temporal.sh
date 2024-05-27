#!/usr/bin/env bash
temporal workflow start \
  --workflow-id aar1 \
  --type WorkflowZoneUp \
  --task-queue zone \
  --input '{"name":"aar1"}'

#!/bin/sh

cd /server/SABook

export PATH="$PATH:/usr/local/instantclient_12_1"

export TNS_ADMIN="/usr/local/instantclient_12_1"
export NLS_LANG=RUSSIAN_CIS.UTF8

sleep 30 && /server/SABook/Exporter.gl -config=/server/SABook/Exporter.json -daemon=YES &


#!/bin/sh

cd /server/SABook/CardDAVMaker

sleep 30 && /server/SABook/CardDAVMaker/CardDAVMaker.gl -config=/server/SABook/CardDAVMaker/CardDAVMaker.json -daemon=YES &

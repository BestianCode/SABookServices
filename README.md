# SABookServices-03
---------------------------------------------------

## Exporter

* aptitide install libaio1
* 
* Install Oracle InstantClient and SDK
* Edit file tnsnames.ora and put to /usr/local/instantclient_12_1
* Add to file /etc/ld.so.conf.d/oracle.conf path: /usr/local/instantclient_12_1
* 
* export GOPATH=/usr/lib/go
* export NLS_LANG=RUSSIAN_CIS.UTF8
* export TNS_ADMIN="/usr/local/instantclient_12_1"
* CGO_CFLAGS=-I/usr/local/instantclient_12_1/sdk/include
* CGO_LDFLAGS=-L/usr/local/instantclient_12_1
* 
* PATH="$PATH:/usr/local/instantclient_12_1"
* 
* go get github.com/go-ldap/ldap
* go get github.com/BestianRU/gounidecode
* go get gopkg.in/goracle.v1
* go get github.com/lib/pq
* go get github.com/denisenkom/go-mssqldb

## AsteriskCIDUpdater

* go get github.com/mattn/go-sqlite3
* go get code.google.com/p/gami

## CardDAVMaker

* ////go get github.com/go-sql-driver/mysql
* go get github.com/ziutek/mymysql/godrv




# SABookServices

aptitide install libaio1

Install Oracle InstantClient and SDK

Edit file tnsnames.ora and put to /usr/local/instantclient_12_1

export GOPATH=/usr/lib/go

export NLS_LANG=RUSSIAN_CIS.UTF8

export TNS_ADMIN="/usr/local/instantclient_12_1"

CGO_CFLAGS=-I/usr/local/instantclient_12_1/sdk/include

CGO_LDFLAGS=-L/usr/local/instantclient_12_1

PATH="$PATH:/usr/local/instantclient_12_1"

go get gopkg.in/goracle.v1

go get github.com/lib/pq

go get gopkg.in/ldap.v1

go get github.com/fiam/gounidecode/unidecode

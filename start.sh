#!/bin/bash

set -e

((count = 30))
while [[ $count -ne 0 ]] ; do
    ping -c 1 db
    rc=$?
    if [[ $rc -eq 0 ]] ; then
        ((count = 1))
    fi
    ((count = count - 1))
done

if [[ $rc -eq 0 ]] ; then
    echo "DB found. Proceeding."
else
    echo "Timeout trying to find DB."
    exit 1
fi

/var/www/html/rrfc/rrfc
/etc/init.d/php7.2-fpm start
nginx -g "daemon off;"

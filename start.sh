#!/bin/bash

set -e

ping -w 30 -c 1 db
rc=$?
if [[ $rc -ne 0 ]]; then
    echo "Could not find DB."
fi

/var/www/html/rrfc/rrfc
/etc/init.d/php7.2-fpm start
nginx -g "daemon off;"

#!/bin/sh
service php8.2-fpm start
service nginx start

cat $(tail -f /var/log/nginx/*) $(tail -f /var/log/php8.2-fpm.log)
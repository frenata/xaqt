#!/usr/bin/env sh

cd /usercode/
sudo service mysql start
echo "----BEGIN-----"
mysql  mysql< /entrypoint/create_user.sql -u'root' 
mysql  ri_db < $1 -u'test' -p'test123'
mysql  mysql< /entrypoint/destroy_user.sql -u'root'
echo "----END----"


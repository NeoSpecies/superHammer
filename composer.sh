EXPECTED_SIGNATURE="$(wget -q -O - https://composer.github.io/installer.sig)"
php -r "copy('https://getcomposer.org/installer', 'composer-setup.php');"
ACTUAL_SIGNATURE="$(php -r "echo hash_file('sha384', 'composer-setup.php');")"

if [ "$EXPECTED_SIGNATURE" != "$ACTUAL_SIGNATURE" ]
then
    >&2 echo 'ERROR: Invalid installer signature'
    rm composer-setup.php
    exit 1
fi

php composer-setup.php --quiet
RESULT=$?
rm composer-setup.php

# 移动 composer.phar 到 /usr/local/bin/ 并重命名为 composer
# sudo mv composer.phar /usr/local/bin/composer
mv composer.phar /usr/local/bin/composer

# 确保 composer 可执行
# sudo chmod +x /usr/local/bin/composer
chmod +x /usr/local/bin/composer
exit $RESULT

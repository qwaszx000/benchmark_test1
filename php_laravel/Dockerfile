
#Thanks to
#https://dev.to/adnanbabakan/dockerizing-laravel-10-ubuntu-image-php-82-fpm-nginx-318p

FROM nginx

#install add-apt-repository
RUN apt update && apt install -y software-properties-common python3-launchpadlib

#php stuff
RUN add-apt-repository -y ppa:ondrej/php
RUN apt update && apt install -y php8.2 php8.2-cli php8.2-zip php8.2-bcmath \
        php8.2-fpm php8.2-curl php8.2-mbstring php8.2-mcrypt php8.2-xml php8.2-pgsql \
        php8.2-pdo php8.2-gd curl

#nginx config
COPY ./nginx/conf.d/default.conf /etc/nginx/conf.d/default.conf

# Copy the application code
COPY . /var/www/html

# Set the working directory
WORKDIR /var/www/html

# Install composer
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer
RUN chmod 755 /usr/local/bin/composer

# Install project dependencies
RUN composer install --optimize-autoloader --no-dev

# Set permissions
RUN chown -R www-data:www-data /var/www/html/storage /var/www/html/bootstrap/cache
RUN chmod +x ./start.sh

#Add nginx user to www-data group -- access to php-fpm socket
RUN usermod -a -G www-data nginx

#optimizations https://laravel.com/docs/10.x/deployment
RUN php artisan config:cache

ENTRYPOINT [ "./start.sh" ]
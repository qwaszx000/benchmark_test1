
#Thanks to
#https://medium.com/@aprilrieger/effortlessly-dockerize-your-php-laravel-app-a-step-by-step-guide-c1a6ffcc2b74

# Use PHP with Apache as the base image
FROM php:8.2-apache as web


RUN apt update && apt install -y libzip-dev zip

# Enable Apache mod_rewrite for URL rewriting
RUN a2enmod rewrite

# Install PHP extensions
RUN docker-php-ext-install pdo_mysql zip

# Configure Apache DocumentRoot to point to Laravel's public directory
# and update Apache configuration files
ENV APACHE_DOCUMENT_ROOT=/var/www/html/public
RUN sed -i -e 's!/var/www/html!${APACHE_DOCUMENT_ROOT}!g' /etc/apache2/sites-available/*.conf
RUN sed -i -e 's!/var/www/!${APACHE_DOCUMENT_ROOT}!g' /etc/apache2/apache2.conf /etc/apache2/conf-available/*.conf

#Change port to 8080(for consistency sake)
RUN sed -i -e "s/:80/:8080/" /etc/apache2/sites-available/000-default.conf
RUN sed -i -e "s/Listen 80/Listen 8080/" /etc/apache2/ports.conf

# Copy the application code
COPY . /var/www/html

# Set the working directory
WORKDIR /var/www/html

# Install composer
RUN curl -sS https://getcomposer.org/installer | php -- --install-dir=/usr/local/bin --filename=composer
RUN chmod 755 /usr/local/bin/composer

# Install project dependencies
RUN composer install

# Set permissions
RUN chown -R www-data:www-data /var/www/html/storage /var/www/html/bootstrap/cache


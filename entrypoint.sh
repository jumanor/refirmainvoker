#!/bin/bash

tail -n +5 /opt/config.properties.example > /opt/config.properties
echo "" >> /opt/config.properties
echo "clientId=$CLIENT_ID" >> /opt/config.properties
echo "clientSecret=$CLIENT_SECRET" >> /opt/config.properties

cd /opt
./main &
httpd-foreground



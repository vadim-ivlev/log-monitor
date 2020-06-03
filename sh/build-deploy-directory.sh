#!/bin/bash

rm -rf deploy/logs
rm -rf deploy/configs
rm -rf deploy/www
rm -rf deploy/notebooks
rm -rf Dockerfile-anaconda3
rm -rf deploy/docker-compose.yml

cp -R logs deploy/
cp -R configs deploy/
cp -R www deploy/
cp -R notebooks deploy/
cp -R Dockerfile-anaconda3 deploy/
cp -R docker-compose.yml deploy/
#!/bin/sh

(cd src && javac fr/itix/soapbackend/*.java)

BASE=tomcat/webapps/axis2/WEB-INF/services/
NAME=VendorBackend
mkdir -p build/$NAME/META-INF build/$NAME/fr/itix/soapbackend/
cp src/fr/itix/soapbackend/*.class build/$NAME/fr/itix/soapbackend/
cp src/services.xml build/$NAME/META-INF/

cp -rv build/$NAME $BASE/



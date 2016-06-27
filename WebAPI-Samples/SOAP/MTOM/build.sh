#!/bin/bash

NAME=SoapBackend

function die() {
  echo "ERROR: $@"
  exit 1
}

cd src || die "No source code"
javac $(find . -iname *.java) || die "Compilation error"
jar cvf ../build/$NAME.aar $(find . -iname *.class -or -iname *.xml) || die "Cannot build jar"



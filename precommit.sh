#!/bin/sh -eu

mage -v check
mage hugoRace
echo "Success"

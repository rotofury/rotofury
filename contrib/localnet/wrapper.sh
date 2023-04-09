#!/usr/bin/env sh

##
## Input parameters
##
BINARY=/gridiron/${BINARY:-gridirond}
ID=${ID:-0}
LOG=${LOG:-gridirond.log}

##
## Assert linux binary
##
if ! [ -f "${BINARY}" ]; then
	echo "The binary $(basename "${BINARY}") cannot be found. Please add the binary to the shared folder. Please use the BINARY environment variable if the name of the binary is not 'gridirond' E.g.: -e BINARY=gridirond_my_test_version"
	exit 1
fi
BINARY_CHECK="$(file "$BINARY" | grep 'ELF 64-bit LSB executable, x86-64')"
if [ -z "${BINARY_CHECK}" ]; then
	echo "Binary needs to be OS linux, ARCH amd64"
	exit 1
fi

##
## Run binary with all parameters
##
export GRIDIRON_HOME="/gridiron/node${ID}/gridirond"

if [ -d "$(dirname "${GRIDIRON_HOME}"/"${LOG}")" ]; then
  "${BINARY}" --home "${GRIDIRON_HOME}" --trace "$@" | tee "${GRIDIRON_HOME}/${LOG}"
else
  "${BINARY}" --home "${GRIDIRON_HOME}" --trace "$@"
fi

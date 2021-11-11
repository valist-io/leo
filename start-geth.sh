#!/bin/sh

geth --rinkeby \
--cache 8192 --cache.preimages \
--http --http.api eth,debug \
--ws --ws.api eth,debug
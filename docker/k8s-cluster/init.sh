#!/bin/sh

apk add --no-cache ca-certificates git bash \
&& make all \
&& make pack4docker \
&& tar -zxf open-falcon-v*.tar.gz -C build \
&& rm open-falcon-v*.tar.gz
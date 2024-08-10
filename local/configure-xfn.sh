#!/bin/bash

# Build the Docker image
docker build -t control-plane-xfn .

# Tag the image with the Gitea server URL
docker tag control-plane-xfn gitea.cnoe.localtest.me:8443/giteaadmin/control-plane-xfn:1

# Log in to the Gitea server
idpbuilder get secrets -p gitea -o json | jq -r '.[0].data.password' | docker login -u giteaAdmin --password-stdin gitea.cnoe.localtest.me:8443

# Push the image to the Gitea server
docker push gitea.cnoe.localtest.me:8443/giteaadmin/control-plane-xfn:1
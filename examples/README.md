# Quickly deploy A CF or k8s cluster.
Your MoltenCore cluster is up and running
so BUCC is up and ready and Concourse is served.

and now you want to setup a Cloudfoundry or k8s Cluster
for linux:
```
./examples/copy-pipeline cf | xclip -selection clipboard
```
for mac:
```
./examples/copy-pipeline cf | pbcopy
```

jump in to a bucc-shell
```
./dev/ssh
mc shell
##paste your clipboard##
```

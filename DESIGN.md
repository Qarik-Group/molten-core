# mc cli design

```
mc init
# after etcd, flannel
# before docker
```
- checks etcd for bucc node
-- starts election if not found
- writes drop-ins for managed unit files
- reload systemd (dbus)
- generate node config (docker certs, host, subnet) and store in etcd
- make flannel subnet persistent (remove ttl)
- write docker certs to disk


```
mc bucc-up
# after docker
```
- load bucc creds from etcds
- bucc int to generate creds
- save creds and vars to etcd
- bucc up


```
mc update-bosh-configs
# after bucc
```
- load bucc creds and vars from etcd
- configure bosh golang client with creds and vars
- load node configs (docker ip, cert and flannel subnet) from etcd
- generate and update cpi-configs
- generate and update cloud-config


```
mc bucc-shell
```
- load bucc creds and vars from etcd
- inject creds and vars into container
- start /bin/bash in bucc container

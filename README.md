# MoltenCore by Stark & Wayne
A lightweight foundation for running containerized platforms on top of bare-metal,
using: [CoreOS Container Linux](https://coreos.com/why/) and
[BUCC](https://github.com/starkandwayne/bucc) (BOSH, UAA, Credhub and Concourse).

## Project Status
This project should not be used for production systems as we still need to tackle:
- Backup & Restore
- Disaster Recovery
- Drain bosh instances on host shutdown
- Re-enable Container Linux Auto Updates

For more details about what we are planning for Phase 3 read [the blog post](https://starkandwayne.com/blog/forging-bare-metal-introducing-molte-core).

## Deployment
Use one of the following terraform projects to deploy a MoltenCore Cluster:

- [Packet MolteCore](https://github.com/starkandwayne/packet-molten-core)
- More to come

Once your cluster is deployed you can check on the status the embedded BUCC service.

## Locating BUCC
The `bucc.service` will be started by systemd on the node with the lowest
internal ip (zone: z0).

Once you have sshed into `z0` systemd can be used to check the status and the
progress of the `bucc.service` deployment.

```
systemctl status bucc.service
journalctl -f -u bucc.service
```

## Accessing BUCC
Make sure to locate your BUCC first (using the above paragraph), and make sure
it is running. Now from `z0` you can start an interactive management shell with:

```
mc shell
```

To find the Concourse login credentials run:

```
bucc info
```

The following cli's have been pre-configured:

The bosh cli [usage docs](https://bosh.io/docs/cli-v2/).
```
bosh env
```

The credhub cli
```
credhub get -n /concourse/main/moltencore
```

The concourse cli, named fly
```
fly -t mc workers
```

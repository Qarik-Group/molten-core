## explanation of dev tools

Requirements:
[ct:](https://github.com/coreos/container-linux-config-transpiler)
[terrafrom 0.12.x:](https://www.terraform.io/
[virtualbox](https://www.virtualbox.org/)
[vagrant](https://www.vagrantup.com/)
[coreos-vagrant](https://github.com/coreos/coreos-vagrant)

### vagrant-up
`./dev/vagrant-up` wil deploy a vagrant box based on the [coreos-vagrant](https://github.com/coreos/coreos-vagrant)
this will spin up 3 coreos nodes

### vagrant-ssh
```
./dev/vagrant ssh   # will connect to the first node
./dev/vagrant ssh 1 # will connect to the second node
./dev/vagrant ssh 2 # will connect to the third node
```

### vagrant-inject-mc
when you have build a new binary you can inject it with: `./dev/vagrant-inject-mc`
this will inject the binary to all deployed nodes and restart the mc.service

### vagrant-resize-box
if you have deployed all your nodes and now you have not enough disk-space
run `./dev/vagrant-resize-box`

### vagrant-destroy
!!DESTROY it ALL!!
`./dev/vagrant-destroy` REMINDER: there is no way to recover from this

### examples/copy-pipeline
BUCC is up and ready and Concourse is served.

you now want to setup a Cloudfoundry or k8s Cluster
for linux:
`./examples/copy-pipeline cf | xclip -selection clipboard`
for mac:
`./examples/copy-pipeline cf | pbcopy`

jump in to a bucc-shell
`./dev/ssh`
`mc shell`
`##paste your clipboard##`

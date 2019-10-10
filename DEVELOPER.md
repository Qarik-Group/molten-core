## explanation of dev tools

Requirements:
[ct:](https://github.com/coreos/container-linux-config-transpiler)
[terrafrom 0.12.x:](https://www.terraform.io/
[virtualbox](https://www.virtualbox.org/)
[vagrant](https://www.vagrantup.com/)
[coreos-vagrant](https://github.com/coreos/coreos-vagrant)

### vagrant-up
```
./dev/vagrant-up
```
wil deploy a vagrant box based on the [coreos-vagrant](https://github.com/coreos/coreos-vagrant)
this will spin up 3 coreos nodes

### vagrant-ssh
```
./dev/vagrant ssh   # will connect to the first node
./dev/vagrant ssh 1 # will connect to the second node
./dev/vagrant ssh 2 # will connect to the third node
```

### vagrant-inject-mc
```
./dev/vagrant-inject-mc
```
this will build and inject the binary to all deployed nodes and restart the mc.service

### vagrant-resize-box
if you have deployed all your nodes and now you have not enough disk-space
run
```
./dev/vagrant-resize-box
```

### vagrant-destroy
!!DESTROY it ALL!!
```
./dev/vagrant-destroy
```
REMINDER: there is no way to recover from this

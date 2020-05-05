# Aimpanel v2

The most easy to use game panel as a service

## Dev

### Frontend

Remember to open `master-frontend/tslint.json` with IntelliJ and apply rules to formatting in popup!

```bash
./Taskfile.sh frontend
```

### Slave VM via Vagrant

Tested on Kubuntu 18.04

```bash
# QEMU / KVM
sudo apt-get install -y qemu libvirt-bin ebtables dnsmasq-base virt-manager
#this allows access to VMs without root, reboot to apply it
sudo adduser `id -un` libvirt-qemu

# Vagrant
#https://www.vagrantup.com/downloads.html
wget https://releases.hashicorp.com/vagrant/2.2.6/vagrant_2.2.6_x86_64.deb
sudo dpkg -i vagrant_2.2.6_x86_64.deb
rm vagrant_2.2.6_x86_64.deb

# Vagrant KVM plugin
#https://github.com/vagrant-libvirt/vagrant-libvirt
sudo apt-get install -y libxslt-dev libxml2-dev libvirt-dev zlib1g-dev ruby-dev
vagrant plugin install vagrant-libvirt

# remember to tune up Vagranfile and install_dev.sh

# Run VM
vagrant up

# Access via SSH
vagrant ssh

# Apply changes to install.sh (token change etc)
vagrant provision

# Destroy VM
vagrant destroy -f
```

## Git tips

Clean old and removed remote branches in local repo
```
git remote prune origin
```

## Swagger

* https://github.com/swaggo/swag/cmd/swag

Swagger allow us to automatically build API docs and client SDK in the future, 
all this without much work from devs

```bash
go get -u github.com/swaggo/swag/cmd/swag
```
### Generate for local purpose

This will generate .json file with swagger spec

```bash
./Taskfile.sh swagger-gen-dev
```

### Serve

This will serve UI for generated swagger spec

```bash
./Taskfile.sh swagger-serve
```

## Slave

### Redis debug

```bash
redis-cli -s /opt/aimpanel/redis/redis.sock
```

### Diagnostics

```bash
ldd --version | grep ldd --color=never
```

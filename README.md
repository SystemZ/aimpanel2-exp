# Aimpanel v2

The most easy to use game panel as a service

## Dev

### Frontend

Remember to open `master-frontend/tslint.json` with IntelliJ and apply rules to automatic formatting in popup!

```bash
./Taskfile.sh frontend
```

### Master backend

Generate TLS cert for local dev env.  
Remember to visit `https://127.0.0.1:3000` and accept self signed cert in a browser at least once.  
If you don't do this, visiting frontend via webpack server will not work.

```bash
# generate private key and self signed TLS certificate
./Taskfile.sh dev-tls-crt
# show generated cert and key
ls -alh *.pem

# use local TLS cert fingerprint to whitelist slave connection to master
openssl x509 -noout -in crt.pem -fingerprint -sha256

# get remote TLS cert fingerprint
# https://askubuntu.com/questions/156620/how-to-verify-the-ssl-fingerprint-by-command-line-wget-curl
echo | openssl s_client -connect example.com:443 |& openssl x509 -fingerprint -sha256 -noout
```

### Code generation

If you change some iota ints or whatever, just run this to regenerate all syntax sugar:
```bash
./Taskfile.sh generate
```
You will need this if you encounter error like:  
> "invalid array index" compiler error signifies that the constant values have changed

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

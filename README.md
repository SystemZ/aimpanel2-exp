# Aimpanel v2

The most easy to use game panel as a service

## Dev launch

If your setup is ready...

Run MongoDB and web interface for it (PMA like)
```bash
docker-compose up
```

Run frontend in WebStorm, then visit http://127.0.0.1:8080
```bash
./Taskfile.sh frontend
# preferably in WebStorm
```

Run VM for slave
```bash
vagrant up
vagrant ssh
sudo -i
```

## Dev deployment

First time? We got you covered :)

### Frontend

Remember to open `master-frontend/tslint.json` with IntelliJ and apply rules to automatic formatting in popup!

```bash
./Taskfile.sh frontend
```

### Backend

#### ENV variables

You need to build and run master binary with something like this to make it work locally:

```
DEV_MODE=true
HTTP_FRONTEND_DIR=/home/user/Projects/aimpanel2/master-frontend/dist/
HTTP_DOCS_DIR=/home/user/Projects/aimpanel2/master/
HTTP_TLS_KEY_PATH=/home/user/Projects/aimpanel2/key.pem
HTTP_TLS_CERT_PATH=/home/user/Projects/aimpanel2/crt.pem
HTTP_TEMPLATE_DIR=/home/user/Projects/aimpanel2/master/templates/
```

#### Cert

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

#### Code generation

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

# remember to tune up Vagrantfile with resources if necessary

# Run VM
./Taskfile.sh up-slave

# Configure hosts
# if this doesn't work just put your local LAN ip in /etc/hosts like this
# 192.168.X.XXX aimpanel.local
./Taskfile.sh update-hosts-slave

# Access via SSH
vagrant ssh

# run this in VM as root to deploy slave
wget --no-check-certificate https://aimpanel.local:3000/i/<HOST TOKEN> -O- | bash -

# WIP
# Apply changes to install.sh (token change etc)
#vagrant provision

# if you are developing slave, use this to sync binary on VM
./Taskfile.sh sync-slave
# or this to auto build, sync and restart slave on VM on slave/* files change
./Taskfile.sh sync-slave-auto
```

# Destroy VM
```
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

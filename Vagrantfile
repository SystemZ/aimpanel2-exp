$script = <<-SCRIPT
Done!
SCRIPT

Vagrant.configure("2") do |config|
  config.vm.box = "generic/ubuntu1804"
  config.vm.box_check_update = false
  config.vm.provision "shell" do |s|
      s.path = "install_dev.sh"
      # uncomment this for lab connection
      #s.path = "install.sh"
      # then uncomment this and replace token
      #s.args = ["CHANGE_TO_TOKEN_FROM_MASTER"]
  end
  config.vm.provision "shell", inline: $script
  config.vm.provider :libvirt do |libvirt|
    libvirt.cpus = 4
    libvirt.memory = 2048
    # this currently doesn't work, standard size is 32 GB
    # should consume just few GBs if VM disk is not used too much
    libvirt.machine_virtual_size = 10
  end
end
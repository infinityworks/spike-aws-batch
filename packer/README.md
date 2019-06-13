# packer-builder

A project to build an AMI and store it in AWS S3

## Pre-requisits

### Install Golang go			
Download from: https://golang.org/dl/ 
Linux version: go1.11.4.linux-amd64.tar.gz
https://dl.google.com/go/go1.11.4.linux-amd64.tar.gz
Extract the binary to /usr/local/

```bash
sudo tar -C /usr/local -xzf go1.11.4.linux-amd64.tar.gz
```

Add go to the path

```bash
vi ~/.bashrc
```

Add the following

```bash
PATH="/usr/local/go/bin:$PATH"
```

Reload bashrc

```bash
. ~/.bashrc
```

Test install

```bash
go version
```

### Setup Assume-Role

Usage & install instructions here: https://github.com/remind101/assume-role 
To install

```bash
go get -u github.com/remind101/assume-role
cp ~/go/bin/assume-role ~/bin/assume-role
```

Test install (should see ENV variables output after MFA prompt)

```bash
assume-role iw-sandpit
```

Add a bash function

```bash
vi ~/.bashrc
```

Add the following

```bash
function assume-role { eval $( $(which assume-role) $@); }

function unassume-role() {
  unset ASSUMED_ROLE
  unset AWS_ACCESS_KEY_ID
  unset AWS_SECRET_ACCESS_KEY
  unset AWS_SECURITY_TOKEN
  unset AWS_SESSION_TOKEN
}
```

Reload bashrc

```bash
. ~/.bashrc
```

Test install

```bash
assume-role iw-sandpit
env | grep ASSUMED_ROLE
```

### Installing Packer
Download the binary from: https://www.packer.io/downloads.html 
Extract the binary into ~/bin/

Test the install with

```bash
packer --version
1.3.3
```

### Installing Ansible

```bash
sudo apt-get update
sudo apt-get install software-properties-common
sudo apt-add-repository --yes --update ppa:ansible/ansible
sudo apt-get install ansible
```

Test the install with

```bash
ansible --version
ansible 2.7.5 
```

## Usage

Authenticate with AWS and get a set of local creds within your env
```bash
assume-role iw-sandpit
```

Run the Packer build
```bash
packer build glacier-restore.json
```

This will output the newly created AMI details to the console and push the new AMI into the account/region specified in glacier-restore.json.
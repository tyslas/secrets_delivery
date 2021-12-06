# secrets_delivery

### setup
1. `git clone git@github.com:tyslas/secrets_delivery.git`
2. `cd secrets_delivery`
3. Ensure that you have the following tools on your machine:
   - `openssl`
   - `ssh`
   - `scp`
   - `tar`
4. Create private and public keys - these will be used for encryption, decryption, and signature verification
   - `./generate_private_public_keypair.sh`
   - send the public key to the person that you would like to exchange information with
5. To encrypt a message and send it to the person that you're exchanging information with run this script:
   ```
   ./handle_secrets.sh </path/to/sender/private_key> </path/to/recipient/public_key> </path/to/message> enc \
       <user on remote server> </path/to/ssh_private_key> <remote server IP address> <directory to copy files to>
   
   example:
   ./handle_secrets.sh ~/.ssh/private_key.pem ~/.ssh/aws_key.pub ./msg.txt enc \
       ec2-user ~/.ssh/tito_aws.pem ec2-54-193-103-167.us-west-1.compute.amazonaws.com /home/ec2-user/secrets_delivery
   ```
6. To decrypt a message that is sent to you make sure you have also cloned this repository and are in the root directory
7. Take note of the directory that the information has been sent to on your machine
   - the encrypted data will be sent in a tar format that the script will extract into a payload folder that contains:
     - `sign.sha256, encrypt.dat, signature.dat`
8. run the script with these arguments:
   ```
   ./handle_secrets.sh </path/to/recipient/private_key> </path/to/sender/public_key> dec

   example:
   ./handle_secrets.sh ~/.ssh/my_private_key.pem ~/.ssh/tito_key.pub dec
   ```
   
### demonstration
1. create and launch an AWS EC2 linux instance from the AWS console
   - add a security group for SSH that allows connections from source 'My IP'
   - use an already created SSH key or have AWS create a new one for you
   - ssh onto the EC2 with a command like this:
     ```
     ssh -i </path/to/your/private/key> ec2-user@<public_ipv4_dns>
     ```
2. install Git on your EC2
   ```
   sudo yum update -y && \
     sudo yum install git -y && \
     git clone https://github.com/tyslas/secrets_delivery.git && \
     cd secrets_delivery
   ```
3. run the script for generating a keypair on your local machine and on the EC2
   ```
   ./generate_private_public_keypair.sh
   ```
4. exchange public keys with `scp` 
   - secure copy the recipient's generated public key from the recipient machine to the sender machine
      ```
      scp -i ~/.ssh/<aws_ssh_key> ec2-user@<public_ipv4_dns>:~/.ssh/my_public_key.pem ~/.ssh/recipient_key.pub
      ```
   - secure copy the sender's generated public key from the sender machine to the recipient machine
      ```
      scp -i ~/.ssh/<aws_ssh_key> ~/.ssh/my_public_key.pem ec2-user@<public_ipv4_dns>:~/.ssh/sender_key.pub
      ```
5. encrypt and send message to recipient with the `./handle_secrets.sh` script
   ```
   ./handle_secrets.sh ~/.ssh/private_key.pem ~/.ssh/aws_key.pub enc ./msg.txt \
       ~/.ssh/<aws_ssh_key> ec2-user <public_ipv4_dns> /home/ec2-user/secrets_delivery
   ```
6. decrypt message from sender on the recipient machine with the `./handle_secrets.sh` script
   - note: be sure to run the decrypt command from the root of the Git directory of this project and that the encrypted message is also at the root of this project
   ```
   ./handle_secrets.sh ~/.ssh/my_private_key.pem ~/sender_key.pub dec
   ```
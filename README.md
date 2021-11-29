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
       ec2-user ~/.ssh/tito_aws.pem ec2-54-193-103-167.us-west-1.compute.amazonaws.com /home/ec2-user
   ```
6. To decrypt a message that is sent to you make sure you have also cloned this repository and are in the root directory
7. Take note of the directory that the information has been sent to on your machine
   - the encrypted data will be sent in a tar format that the script will extract into a payload folder that contains:
     - `sign.sha256, encrypt.dat, signature.dat`
8. run the script with these arguments:
   ```
   ./handle_secrets.sh </path/to/recipient/private_key> </path/to/sender/public_key> \
       </path/to/encrypted/data> dec </home/ec2-user/payload/signature.dat>

   example:
   ./handle_secrets.sh /home/ec2-user/.ssh/my_private_key.pem /home/ec2-user/mac_public_key.pem \
       /home/ec2-user/payload/encrypt.dat dec /home/ec2-user/payload/signature.dat
   ```
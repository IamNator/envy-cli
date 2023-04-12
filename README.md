# envy-cli

CLI Client for envy


ENVY is a cli tool for managing env variable across platforms.
[goto ENVY](https://github.com/IamNator/envy-download)

## Motivations
I needed a simple and inexpensive way of managing env variable. The best part of this, is that it's self hosted.

#### perks
Imagine, never having to copy and paste env variable again.



## INSTALLATION
```
  $ curl -sL https://raw.githubusercontent.com/IamNator/envy-cli/main/scripts/install.sh | bash
```

## HOW

1. The env varibles are encryted then encoded before uploading to a remote server
2. When retrieved, they are decoded then decryted before display or writing to a file
3. A secret key is used for both encrypting and decrypting the values. Please don't loose or forget the key
4. They Secret keys are not stored by the remote server, they are only known to the client


## USE


1. set a single variable
```
  $ envy -set AWS_SECRET=aeewq45243gfe
  
  output: 
  $ done
```

2. get a signle variable
```
  $ envy -get AWS_SECRET
  
  output:
  $ AWS_SECRET=aeewq45243gfe
```

3. upload env variables from a file
```
  $ envy -source .env 
  
  output:
  $ uploading
  $ done
```

4. download env variable to a *.env file
```
  $ envy -dir /home/staging.env
  
  output:
  $ downloading
  $ done.
```

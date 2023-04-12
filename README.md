# envy-cli
CLI Client for envy


ENVY is a cli tool for managing env variable across platforms.


You can set, get or retrieve to a .env file.

[CLOUD SERVER](https://github.com/IamNator/envy-download)

HOW

The env variable are encryted before they are sent to the cloud. When they are retrieved, a key is also used to decryt them


USE


1. set a single varible
```
  $ envy -set AWS_SECRET=aeewq45243gfe
```

2. get a signle variable
```
  $ envy -get AWS_SECRET
```

3. upload env variables from a file
```
  $ envy -source .env 
```

4. download env variable to a *.env file
```
  $ envy -dir /home/staging.env
```

EXLEA
==================================

Deploy Notes
----------------------------------

- to build
`docker image build -t victron/exleacar:latest .`
or with ver. tag
`docker image build -t victron/exleacar:2020-06-21 -t victron/exleacar:latest .`
- push 

- to run 
`docker run --name exlea -d --restart=unless-stopped --net=host -v /home/ubuntu/exle:/home --user $(id -u):$(id -g) victron/exleacar:latest -u <user> -p <password> -w 30 -vvv`


- rebuild builder image
`docker image build --target builder -t victron/exleacar_builder:latest .`
# ad-server
I have used beego MVC framework and visual studio code editor for the development of this project. You can setup the project by following way:

### You will find a mysql db backup file named ad_server_db.sql. Actually, It contains all the necessary statements to create  db tables and insert some dummy data.
### Install golang from here: https://golang.org/dl/. 

After installing go please set the gopath by following way:

* Open bash_profile using: vi ~/.bash_profile
* Then add following lines:
```
export GOROOT=/usr/local/go
export PATH=$GOROOT/bin:$PATH
export GOPATH=/Users/{your_user_name}/work (In my case it’s ‘rakib’. This will be the work directory)
export PATH=$GOPATH/bin:$PATH
```
* You will need to install Beego and the Bee dev tool. Run the following commands from the terminal to install the beego:
```
$ go get github.com/astaxie/beego
$ go get github.com/beego/bee
```
* Now copy ad-server folder. Then go to cd $GOPATH/src directory and paste my ‘ad-server’ folder.
* Go to cd $GOPATH/src/ad-server directory

* Before you run the project, please change the mysql connection string which will be found in main.go file at line #13.
orm.RegisterDataBase("default", "mysql", "username:password@/ad_server_db?charset=utf8")
* Execute following command in the terminal to run the solution:  
```
bee run
```
* It will run the project at HTTP port 8080
* Open any browser and then hit this url http://localhost:8080 to see the running ad-server.



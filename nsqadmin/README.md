## nsqadmin

`nsqadmin` is a Web UI to view aggregated cluster stats in realtime and perform various
administrative tasks.

Read the [docs](http://nsq.io/components/nsqadmin.html)


由于https://www.npmjs.com在国内访问不稳定，因此建议使用国内镜向站点https://npm.taobao.org



## Working Locally

 1. `$ npm install`                 ## downlaod gulp-sass issue if you occur
 2. `$ ./gulp clean watch` or `$ ./gulp clean build`
 3. `$ go-bindata --debug --pkg=nsqadmin --prefix=static/build static/build/...`
 4. `$ go build && ./nsqadmin`
 5. make changes (repeat step 5 if you make changes to any Go code)
 6. `$ go-bindata --pkg=nsqadmin --prefix=static/build static/build/...`
 7. commit other changes and `bindata.go`





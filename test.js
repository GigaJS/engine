const http = require('http')
const mongodb = require('mongodb')
// console.log(mongodb)
mongodb.createClient({url: "mongodb://root:jABZKj0Wrk@127.0.0.1:27017/?authSource=admin&readPreference=primary&appname=MongoDB%20Compass&directConnection=true&ssl=false"}).then(client => {
    let users = client.db("qwave").collection('users');
    // console.log(users)
    console.log('RES: ', users.find())
}).catch(err => {
    console.log(err)
})

http.request({
    url: "https://postman-echo.com/get?foo1=bar1&foo2=bar2",
    headers: {
        'Test': 'Headers'
    }
}).then(ok => {
    console.log(ok);
    console.log(ok.status);
    console.log(ok.body.toString());
}).catch(err => {
    console.log(err);
})


const server = http.createServer()

const group = server.createGroup('/')
group.get('/world', ctx => {
    ctx.reply({
        hello: 'user'
    })
})

server.listen(':3000')
setTimeout(() => {

}, 60000)

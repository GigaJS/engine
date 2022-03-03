const dbURL = "mongodb://root:jABZKj0Wrk@127.0.0.1:27017/?authSource=admin&readPreference=primary&directConnection=true&ssl=false"

const http = require('http')
const mongodb = require('mongodb')

mongodb.createClient({url: dbURL}).then(client => {
    let users = client.db("qwave").collection('users');
    users.findOne({role: {$gt: 6, $lt: 9}}, { skip: 1 }).then(res => {
        console.log('RES: ', res)
    }).catch(err => {
        console.log("Error ", err)
    })
}).catch(err => {
    console.log("Error")
    console.log(err)
})


/*

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
*/

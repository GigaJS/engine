const dbURL = "mongodb://root:jABZKj0Wrk@127.0.0.1:27017/?authSource=admin&readPreference=primary&directConnection=true&ssl=false"

const http = require('http')
const mongodb = require('mongodb')
const colors = require('colors')

function info(...msgs) {
    console.log(colors.brightBlue + "info   " + colors.reset, ...msgs)
}
function err(...msgs) {
    console.log(colors.brightRed + colors.bold + "error  " + colors.reset, ...msgs)
}

info("Listening users")
mongodb.createClient({ url: dbURL }).then(client => {
    info("Mongodb connection success")
    let db = client.db("qwave");
    let users = db.collection('users');
    let qwave = db.collection('qwave');
    // users.findOne({ role: { $gt: 6, $lt: 9 } }, { skip: 0 }).then(res => {
    // console.log('RES: ', res)
    // }).catch(err => {
    // console.log("Error ", err)
    // })

err("pizda")

    qwave.insert({ test: true, date: Date.now() }).then(id => {
        info("Insert id: " + id)
    })

    qwave.find({ test: true }).then(results => {
        info("Qwave resuts: ", results)
    })

    users.find({ role: 10 }, { limit: 1 }).then(result => {
        console.log(result)
    }).catch(err => {
        console.log('Error occured ', err)
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

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

function await(promise, cb) {
    promise
        .then((ok) => {
            try {
                cb(ok, null)
            } catch (e) {
                err("executing cb failed ", e)
            }
        })
        .catch(err => {
            console.log("Promise error ", err)
            cb(null, err)
        })
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

    try {
        console.log("a ", 'f')
        await(qwave.find({}), (res, err) => {
            console.log("a ", 'g')

            info(res)
            await(qwave.deleteMany({}), (res, err) => {
                info("Deleted ", res)

                await(qwave.insert({ test: true, date: Date.now }), (res, err) => {
                    info("Inserted entry with id ", res)
                    await(qwave.findOne({ _id: mongodb.objectId(res) }), (res) => {
                        info("Found ", res)
                    })
                });
            });
        });
    } catch (e) {
        console.log(e)
    }

    console.log("ff")
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


import express from 'express';
import bodyParser from 'body-parser';
import {join} from 'path'

const app = express();

const ROOT = 'js.bottle.remotehack.space';

var urlencodedParser = bodyParser.urlencoded({ extended: false })

app.all('/', urlencodedParser, (req, res, next) => {

    const {host} = req.headers;

    if(host === ROOT) {
        // 
        if(req.method === 'GET') {
            res.sendFile(join(__dirname, '/../www/landing.html'))
            return;
        }

        if(req.method === 'POST') {
            res.send(`TODO: CREATE SUBDOMAIN: ${req.body.subdomain}`)
        
            return;
        }

    }

    next();

})

app.all('*', (req, res) => {
    const {host} = req.headers;
    res.send(`(NODE) FALLBACK2 - ${host}`)
})

app.listen(8000, () => {
    console.log("started")
});

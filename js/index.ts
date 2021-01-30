
import express from 'express';
import bodyParser from 'body-parser';
import {join} from 'path'
import { RequestsStore } from './request-store';
import { valid } from './util';

const app = express();

const ROOT = 'js.bottle.remotehack.space';

const SUBREG = /^(?<subdomain>\w[\w-]*)\.js\.bottle\.remotehack\.space$/;

const store = new RequestsStore({storeDir: join(__dirname, '/../data')});

var urlencodedParser = bodyParser.urlencoded({ extended: false })

app.all('/', urlencodedParser, async (req, res, next) => {

    const {host} = req.headers;

    if(host === ROOT) {
        // 
        if(req.method === 'GET') {
            res.sendFile(join(__dirname, '/../www/landing.html'))
            return;
        }

        if(req.method === 'POST') {

            const {subdomain} = req.body

            if(valid(subdomain)) {
                if(await store.exists({subdomain})) {

                    res.redirect(subdomain + ROOT)
    
                    return;
                } else {

                    await store.createStore({subdomain})

                    res.send(`CREATED SUBDOMAIN: ${req.body.subdomain}`)
                    return;
                }
            }


            return;
        }

    }

    next();

})

app.all('*', async (req, res) => {
    const {host} = req.headers;

    const match =host?.match(SUBREG);
    if(match) {

        const {subdomain} = (match.groups as any)


        const timestamp = Date.now().toString();
        const content = JSON.stringify([req.headers, req.url])

        await store.write({subdomain, timestamp, content})

        const result = await store.read({subdomain});

        res.contentType('text/plain');
        res.send(result);

        return;
    }

    res.send(`(NODE) Cool `)
})

app.listen(8000, () => {
    console.log("started")
});

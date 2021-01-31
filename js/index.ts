
import express from 'express';
import bodyParser from 'body-parser';
import {join} from 'path'
import { RequestsStore } from './request-store';
import { valid } from './util';

const app = express();

const PORT = 8000;

const ROOT = 'js.bottle.remotehack.space';

const SUBREG = /^(?<subdomain>\w[\w-]*)\.js\.bottle\.remotehack\.space$/;

const store = new RequestsStore({storeDir: join(__dirname, '/../data')});

var urlencodedParser = bodyParser.urlencoded({ extended: false })

app.all('/', urlencodedParser, async (req, res, next) => {

    const {host} = req.headers;

    if (host !== ROOT) {
        next();
        return;
    }

    if(req.method === 'GET') {
        res.sendFile(join(__dirname, '/../www/landing.html'))
        return;
    }

    if (req.method === 'POST') {

        const {subdomain} = req.body

        if (valid(subdomain)) {
            if (!await store.exists({subdomain})) {
                await store.createStore({subdomain})
            }

            res.redirect(`http://${subdomain}.${ROOT}`)
        } else {
            res.status(404).send(`${subdomain} is not a valid subdomain name!`);
        }

        return;
    }
})

app.all('*', async (req, res) => {
    const {host} = req.headers;

    const match = host?.match(SUBREG);
    const {subdomain} = match?.groups || {};

    if (!match || !subdomain) {
        const msg = `Failed to match for subdomain on host ${host}`;
        console.error(msg);
        res.status(404).send(msg);
        return;
    }

    const timestamp = Date.now().toString();
    const content = JSON.stringify([req.headers, req.url])
    try {
        await store.write({subdomain, timestamp, content})

        const result = await store.read({subdomain});

        res.contentType('text/plain');
        res.send(result);
    } catch (e) {
        console.error(e);
        res.status(404).send(`Nobody has created this domain yet.`);
    }
})

app.listen(PORT, () => {
    console.log(`started on port ${PORT}`);
});


// A web server that listens to 8080
import express from 'express';

const app = express();

app.get('/', (req, res) => {res.send("hello from node!")})


app.listen(8000, () => {
    console.log("started")
});

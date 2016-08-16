const http = require('http');
const os = require('os');
const express = require('express')

const PORT = 3333;

const app = express();
let healthCheck = true;

const home = (req, res) => {
  healthCheck ? res.status(200) : res.status(500);
  res.write(`Current healthcheck response : ${healthCheck}`);
  res.end();
}

app.get('/', home);

app.get('/switch', (req, res) => {
  healthCheck = (healthCheck == false)
  home(req, res);
})

app.listen(PORT);


console.log('Running on port ' + PORT);

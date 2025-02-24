const HyperExpress = require('hyper-express');
const webserver = new HyperExpress.Server();

// Create GET route to serve 'Hello World'
webserver.get('/test_plain', (request, response) => {
    response.type("text/plain")
    response.send('Hello world!');
})

webserver.listen(8080)
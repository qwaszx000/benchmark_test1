const express = require('express')
const app = express()
const port = 8080

app.get('/test_plain', (req, res) => {
    res.type("text/plain")
    res.send('Hello world!')
})

app.listen(port, () => {})


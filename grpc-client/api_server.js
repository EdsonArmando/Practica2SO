const express = require('express');
const app = express();
var morgan = require('morgan');
var cors = require('cors');

//Settings
const port = 5000;

//Middlewares
app.use(express.json());
app.use(morgan('dev'));
app.use(cors());

//Routes
//app.use('/',require('./routes/index'))
app.use('/',require('./routes/juegos.js'));

app.listen(port,()=>{
    console.log('Servidor en el puerto, maquina 1 : ', port);
});
var express = require('express');
var router = express.Router();
const client = require('../gRPC_client')

router.post('/ejecutarJuego',  function(req, res) {
    const data_juego = {
        game_id : req.body.game_id,
        players : req.body.players,
    }
    
    client.EjecutarJuego(data_juego, function(err, response) {
        res.status(200).json({mensaje: response.message})
    });
});

router.get('/listarJuegos',  function(req, res) {
    const rows = [];
    const call = client.LogsJuego();
    call.on('data', function(data) {
        rows.push(data);
    });
    call.on('end', function() {
        console.log('Data obtenida con exito');
        res.status(200).json({data:rows});
    });
    call.on('error', function(e) {
        console.log('Error al obtener la data',e);
    });
    /*
    call.on('status', function(status) {
        // process status
    });
    */
});

router.get('/',  function(req, res) {
    res.status(200).json(1);
});

module.exports = router;
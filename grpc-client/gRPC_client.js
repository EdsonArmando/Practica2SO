var PROTO_PATH = './proto/fase2.proto';

var parseArgs = require('minimist');
var grpc = require('@grpc/grpc-js');
var protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
     longs: String,
     enums: String,
     defaults: true,
     oneofs: true
    });
var juego_proto = grpc.loadPackageDefinition(packageDefinition).proto;

var argv = parseArgs(process.argv.slice(2), {
  string: 'target'
});
var target;
if (argv.target) {
  target = argv.target;
} else {
  target = '127.0.0.1:50051';
}
var client = new juego_proto.IniciarJuego(target, grpc.credentials.createInsecure());
module.exports = client;

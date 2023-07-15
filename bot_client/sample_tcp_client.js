config = require('../config.json');
net = require('net');

// Variables
const HOSTNAME = config.ServerIP;
const PORT = config.TcpPort;
var socket;

// Connect to server
var socket = net.connect(PORT, HOSTNAME, () => {
  console.log("connected")
}).on("error", (err)=>{console.log("couldn't connect"); process.exit()});

// Console input
process.stdin.on("data", (input) => {
  socket.write(input.toString().trim());
})

// Communication
socket.on("data", (data) => {
  console.log(data.toString().trim());
})

// Connection closed
socket.on("end", () => {
  console.log("disconnected")
  process.exit()
})
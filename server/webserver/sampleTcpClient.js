// These require statements need Node v18+ --> for older versions of Node,
// replace "A = require('B');" with "import A from 'B';"
config = require("../../config.json");
net = require("net");

// Variables
const HOSTNAME = config.ServerIP;
const PORT = config.TcpPort;
var socket;

// Connect to server
var socket = net.connect(PORT, HOSTNAME, () => {
  console.log("connected")
}).on("error", (err)=>{console.log("couldn't connect\n", err); process.exit()});

// Console input
process.stdin.on("data", (input) => {
  console.log(input)
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
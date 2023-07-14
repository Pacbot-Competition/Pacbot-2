<script>
  const socket = new WebSocket('ws://localhost:3002');
  let good = false;

  socket.addEventListener('open', (event) => {
    good = true;
    console.log('WebSocket connection established');
    socket.send('Hello, server!');
  });

  socket.addEventListener('message', (event) => {
    console.log(`Received: ${event.data}`);
  });

  socket.addEventListener('close', (event) => {
    console.log('WebSocket connection closed');
    good = false;
  });
</script>

<h1> {good ? "Good!" : "Offline"} </h1>
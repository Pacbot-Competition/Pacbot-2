<style>

  :root {
    font-family: Inter, system-ui, Avenir, Helvetica, Arial, sans-serif;

    color-scheme: light dark;
    color: rgba(255, 255, 255, 0.87);
    background-color: #242424;

    font-synthesis: none;
    text-rendering: optimizeLegibility;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    -webkit-text-size-adjust: 100%;
  }

  .maze-space {

    /* Positioning */
    position: absolute;
    top: 5vh;
    left: 5vw;
  }

</style>

<script>

  // Imports
  import config from '../../config.json';
  import Maze from './lib/Maze.svelte';
  import Pellets from './lib/Pellets.svelte';
  import Pacman from './lib/Pacman.svelte';
  import Ghost from './lib/Ghost.svelte';
  import MpsCounter from './lib/MpsCounter.svelte';
  import Ticker from './lib/Ticker.svelte';

  // Creating a websocket client
  var socket = new WebSocket(`ws://${config.ServerIP}:${config.WebSocketPort}`);
  socket.binaryType = 'arraybuffer';
  let socketOpen = false;

  /* TODO - Add error handling for when the websocket client is not open, to avoid console errors */

  /* 
    This generates an empty array of pellet states 
    (0 = none, 1 = pellet, 2 = super)
  */
  let pelletGrid = [];
  for (let row = 0; row < 31; row++) {
    pelletGrid[row] = [];
    for (let col = 0; col < 28; col++) {
      pelletGrid[row][col] = 0;
    }
  }

  /*
    We use a circular queue (with a fixed max capacity to keep track of the times
    of the most recent messages. For every message we receive, we should add this
    time to the queue and remove all times longer than a millisecond ago. The 
    length of this array will be the MPS (messages per second), which should be
    synced with the frame rate of the game engine if there is no lag.
  */
  const MPS_BUFFER_SIZE = 2 * config.GameFPS; // Allow double the specified FPS
  let mpsBuffer = new Array(MPS_BUFFER_SIZE);
  let mpsIdxLeft = 0;
  let mpsIdxRight = 1;
  let mpsAvg = 0;
  mpsBuffer[0] = Date.now();

  // Keep track of the number of ticks elapsed (from the server)
  let currTicks = 0;

  // Keep track of the ticks per update (from the server)
  let updatePeriod = 12;

  // Keep track of the game mode (from the server)
  /* TODO: Make use of game mode for changing the ticker color */
  let gameMode = 0;

  // Initial states for all the agents
  let pacmanRowState = 23;
  let pacmanColState = 14;

  let redRowState = 11;
  let redColState = 13 | 0xc0; // left
  let redFrightState = 0 | 128;

  let pinkRowState = 13 | 0x40; // down
  let pinkColState = 13;
  let pinkFrightState = 0 | 128;

  let cyanRowState = 14 | 0xc0; // up
  let cyanColState = 11;
  let cyanFrightState = 0 | 128;

  let orangeRowState = 14 | 0xc0; // up
  let orangeColState = 15;
  let orangeFrightState = 0 | 128;

  // Handling a new connection
  socket.addEventListener('open', (_) => {
    console.log('WebSocket connection established');
    socket.send('Hello, server!');
    socketOpen = true;
  });

  // When the ticker is clicked, send a message
  let paused = false;
  $: {if (socketOpen) {
    socket.send(paused ? 'p' : 'P');
  }}

  // Message events
  socket.addEventListener('message', (event) => {
    if (event.data instanceof ArrayBuffer) {

      // Log the time
      const ts = Date.now();
      mpsBuffer[mpsIdxRight] = ts;
      mpsAvg++;
      mpsIdxRight = (mpsIdxRight + 1) % MPS_BUFFER_SIZE;
      // Leeway of about 2%
      while (ts - mpsBuffer[mpsIdxLeft] > 1020 && mpsIdxLeft != mpsIdxRight) {
        mpsIdxLeft = (mpsIdxLeft + 1) % MPS_BUFFER_SIZE;
        mpsAvg--;
      }

      // Binary frame for manually parsing the input
      let view = new DataView(event.data);

      if (view) {

        // Keep track of the byte index we are reading
        let byteIdx = 0;

        // Get the current ticks from the server
        currTicks = view.getUint16(byteIdx, false);
        byteIdx += 2;

        // Get the update ticks from the server
        updatePeriod = view.getUint8(byteIdx, false);
        byteIdx += 1;

        // Get the game mode from the server
        gameMode = view.getUint8(byteIdx, false);
        byteIdx += 1;

        // Parse ghost data
        redRowState       = view.getUint8(byteIdx, false); byteIdx += 1;
        redColState       = view.getUint8(byteIdx, false); byteIdx += 1;
        redFrightState    = view.getUint8(byteIdx, false); byteIdx += 1;
        pinkRowState      = view.getUint8(byteIdx, false); byteIdx += 1;
        pinkColState      = view.getUint8(byteIdx, false); byteIdx += 1;
        pinkFrightState   = view.getUint8(byteIdx, false); byteIdx += 1;
        cyanRowState      = view.getUint8(byteIdx, false); byteIdx += 1;
        cyanColState      = view.getUint8(byteIdx, false); byteIdx += 1;
        cyanFrightState   = view.getUint8(byteIdx, false); byteIdx += 1;
        orangeRowState    = view.getUint8(byteIdx, false); byteIdx += 1;
        orangeColState    = view.getUint8(byteIdx, false); byteIdx += 1;
        orangeFrightState = view.getUint8(byteIdx, false); byteIdx += 1;

        // Parse pellet data
        for (let row = 0; row < 31; row++) {
          const binRow = view.getUint32(byteIdx, false);
          for (let col = 0; col < 28; col++) {
            let superPellet = ((row === 3) || (row === 23)) && ((col === 1) || (col === 26));
            pelletGrid[row][col] = ((binRow >> col) & 1) ? (superPellet ? 2 : 1) : 0;
          }
          byteIdx += 4;
        }

        // Trigger an update for the pellets
        pelletGrid = pelletGrid;
      }
    }
  });

  // Event on close
  socket.addEventListener('close', (_) => {
    socketOpen = false;
    console.log('WebSocket connection closed');
  });

  // Track the size of the window, to determine the grid size
  let innerWidth = 0;
  let innerHeight = 0;
  $: gridSize = 0.9 * ((innerHeight * 28 < innerWidth * 31) ? (innerHeight / 31) : (innerWidth / 28));

  // Calculate the remainder when currTicks is divided by updatePeriod
  $: modTicks = currTicks % updatePeriod

</script>

<svelte:window bind:innerWidth bind:innerHeight />

<div class='maze-space'>
  <Maze {gridSize} />
  <Pellets {pelletGrid} {gridSize} />
  <Pacman {gridSize} {pacmanRowState} {pacmanColState} />

  <!-- SVG Ghost Sprites -->
  <Ghost {gridSize}
         {modTicks}
         {updatePeriod} 
         rowState={redRowState}
         colState={redColState}
         frightState={redFrightState}
         color='red'/>
  
  <Ghost {gridSize}
         {modTicks}
         {updatePeriod}
         rowState={pinkRowState}
         colState={pinkColState}
         frightState={pinkFrightState}
         color='pink'/>

  <Ghost {gridSize}
         {modTicks}
         {updatePeriod}
         rowState={cyanRowState}
         colState={cyanColState}
         frightState={cyanFrightState}
         color='cyan'/>

  <Ghost {gridSize}
         {modTicks}
         {updatePeriod} 
         rowState={orangeRowState} 
         colState={orangeColState} 
         frightState={orangeFrightState}
         color='orange'/>

  <MpsCounter {gridSize} {mpsAvg} />
  <Ticker {gridSize} {modTicks} {updatePeriod} bind:paused/>
</div>
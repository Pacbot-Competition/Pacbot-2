<style>

  .pacman {
    position: absolute;
    top: 0;
    left: 0;
    background-color: yellow;
    clip-path: polygon(-100% -100%, -100% 200%, 200% 200%, 100% 80%, 45% 50%, 100% 20%, 200% -100%);
    border-radius: 50%;
  }

  .eating {
    animation: eat 0.5s ease 1;
  }

  @keyframes eat {
    50% {
      clip-path: polygon(-100% -100%, -100% 200%, 200% 200%, 100% 50%, 45% 50%, 100% 50%, 200% -100%);
    }
  }

</style>

<script>
  export let gridSize;
  export let pacmanRow;
  export let pacmanCol;

  $: posX = gridSize * pacmanCol;
  $: posY = gridSize * pacmanRow;

</script>

<div class="pacman eating" 
     style:width="{gridSize}px" 
     style:height="{gridSize}px" 
     style:left="{posX}px"
     style:top="{posY}px"/>
<style>

  .pacman {

    /* Positioning */
    position: absolute;

    /* Formatting */
    clip-path: polygon(-100% -100%, -100% 200%, 200% 200%, 100% 80%, 45% 50%, 100% 20%, 200% -100%);
    background-color: yellow;
    border-radius: 50%;

    /* Grid-size related CSS */
    width: var(--grid-size);
    height: var(--grid-size);
    transform: scale(1.5);
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
  export let pacmanRowState;
  export let pacmanColState;

  $: posX = gridSize * pacmanColState;
  $: posY = gridSize * pacmanRowState;

</script>

<div class="pacman eating" 
     style:--grid-size="{gridSize}px"
     style:left="{posX}px"
     style:top="{posY}px"/>
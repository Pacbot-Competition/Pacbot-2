<style>

  /* Pacman (yellow circle) */
  .pacman {

    /* Positioning */
    position: absolute;

    /* Formatting */
    background-color: yellow;
    border-radius: 50%;

    /* Grid-size related CSS */
    width: var(--grid-size);
    height: var(--grid-size);
    transform: scale(1.5) rotate(var(--dir-angle));
  }

  /* Clip path (mouth of Pacman) */
  .clip {
    clip-path: polygon(-100% -100%, -100% 200%,
                        200% 200%, 100% 80%, 45% 50%,
                        100% 20%, 200% -100%);
  }

  /* Pacman eating */
  .eating {
    animation: eat 0.5s ease 1;
  }

  /* Pacman eating animation (clip path changes) */
  @keyframes eat {
    50% {
      clip-path: polygon(-100% -100%, -100% 200%,
                          200% 200%, 100% 50%, 45% 50%,
                          100% 50%, 200% -100%);
    }
  }

</style>

<script>

  // Grid size, same as for other components
  export let gridSize;

  // Pacman state
  export let pacmanRowState;
  export let pacmanColState;

  // Using the & operator to pick out the 5 lowest bits
  $: posX = pacmanColState & 0b11111
  $: posY = pacmanRowState & 0b11111

  // Hide the Pacman if bit 5 (32) of either coordinate is set
  $: showPacman = ((pacmanRowState | pacmanColState) & 0b100000) ? false : true;

  // The below code is a sign-extension trick, taking advantage of 32-bit
  // integer representations in JavaScript
  $: dirX = ((pacmanColState >> 6) << 30) >> 30
  $: dirY = ((pacmanRowState >> 6) << 30) >> 30

  // Based on the direction, decide the rotation amount
  let rotation = 0;
  let clip = true;
  $: {
    clip = true;
    if (dirX === 1 && dirY === 0) {
      rotation = 0;
    } else if (dirY === 1 && dirX === 0) {
      rotation = 90;
    } else if (dirX === -1 && dirY === 0) {
      rotation = 180;
    } else if (dirY === -1 && dirX === 0) {
      rotation = 270;
    } else {
      clip = false;
    }
  }

</script>

{#if showPacman}
  <div
    class='pacman {clip ? 'clip' : ''} eating'
    style:--grid-size='{gridSize}px'
    style:--dir-angle='{rotation}deg'
    style:left='{gridSize * posX}px'
    style:top='{gridSize * posY}px'
  />
{/if}
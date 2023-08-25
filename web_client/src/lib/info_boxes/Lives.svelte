<style>

  /* Information box for remaining lives */
  .lives-box {

    /* Positioning */
    position: absolute;
    text-align: center;
    z-index: 1;

    /* Formatting */
    background-color: rgba(0,0,0,0.3);
    border: none;
    display: flex;
    align-items: center;
    justify-content: center;

    /* Grid-size related CSS */
    width:       calc(5   * var(--grid-size));
    height:      calc(3   * var(--grid-size));
    line-height: calc(1.1 * var(--grid-size));
    left:        calc(23   * var(--grid-size));
    top:         calc(16  * var(--grid-size));
    font-size:   calc(0.8 * var(--grid-size));
  }

</style>

<script>
  import Pacman from "../agents/Pacman.svelte";

  export let gridSize;
  export let currLives;
  export let Directions;
</script>

<div class='lives-box' style:--grid-size="{gridSize}px">
  {#if currLives > 1}
    <Pacman
      {gridSize}
      pacmanRowState={1}
      pacmanColState={1 | Directions.Right}
    />
  {/if}
  {#if currLives > 2}
    <Pacman
      {gridSize}
      pacmanRowState={1}
      pacmanColState={3 | Directions.Right}
    />
  {/if}
</div>
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
    left:        calc(23  * var(--grid-size));
    top:         calc(16  * var(--grid-size));
    font-size:   calc(0.8 * var(--grid-size));
  }

  /* Text to show extra lives, above the total 3 */
  .extra-lives-text {

    /* Positioning */
    position: absolute;
    z-index: 2;

    /* Grid-size related CSS */
    right:       calc(0.2  * var(--grid-size));
    top:         calc(0.2 * var(--grid-size));
    font-size:   calc(0.5  * var(--grid-size));
    line-height: calc(0.5  * var(--grid-size));
  }

  /* Text to show GAME OVER once no lives are left */
  .game-over-text {

    /* Positioning */
    position: absolute;
    text-align: center;
    z-index: 1;

    /* Formatting */
    color: lightcoral;
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;

    width:       calc(5   * var(--grid-size));
    height:      calc(3   * var(--grid-size));
    line-height: calc(1.1 * var(--grid-size));
    left:        calc(0   * var(--grid-size));
    top:         calc(0   * var(--grid-size));
    font-size:   calc(0.8 * var(--grid-size));
  }

  /* Control of header margins */
    h2.game-over-text {
    margin: 0;
  }

</style>

<script>
  import Pacman from "../agents/Pacman.svelte";

  export let gridSize;
  export let currLives;
  export let Directions;
</script>

<div
  class='lives-box'
  style:--grid-size='{gridSize}px'
>
  {#if currLives == 0}
    <!-- No lives left -->
    <h2 class='game-over-text'>
      <div>GAME</div>
      <div>OVER</div>
    </h2>
  {/if}
  {#if currLives > 1}
    <!-- Second Pacman life -->
    <Pacman
      {gridSize}
      pacmanRowState={1}
      pacmanColState={1 | Directions.Right}
    />
  {/if}
  {#if currLives > 2}
    <!-- Third Pacman life -->
    <Pacman
      {gridSize}
      pacmanRowState={1}
      pacmanColState={3 | Directions.Right}
    />
  {/if}
  {#if currLives > 3}
    <!-- More than 3 lives -->
    <div class='extra-lives-text'>
      +{currLives - 3}
    </div>
  {/if}
</div>
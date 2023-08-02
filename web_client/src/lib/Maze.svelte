<style>
  .maze {

    /* Positioning */
    position: absolute;
    top: 0;
    left: 0;
    z-index: -1;

    /* Interaction */
    user-select: none;

    /* Grid-size related CSS */
    width:  calc(28 * var(--grid-size));
    height: calc(31 * var(--grid-size));
  }
</style>

<script>
  export let gridSize;
</script>

<img class="maze" 
     src="maze.svg" 
     alt='Pacman Maze' 
     style:--grid-size="{gridSize}px" 
     draggable="false"/>
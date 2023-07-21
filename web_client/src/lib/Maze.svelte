<style>
  .maze {
    position: absolute;
    top: 0;
    left: 0;
    user-select: none;
  }
</style>

<script>
  export let gridSize;
</script>

<img class="maze" src="maze.png" alt='Pacman Maze' style:width="{28*gridSize}px" style:height="{31*gridSize}px" draggable="false"/>
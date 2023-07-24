<!-- MPS = Messages Per Second -->

<style>

  .mps-counter {
    position: absolute;
    background-color: rgba(0,0,0,0.3);
    text-align: center;
  }

</style>

<script>
  export let gridSize;
  export let mpsAvg;
</script>

<h2 class='mps-counter' 
     style:width="{5*gridSize}px"
     style:height="{3*gridSize}px"
     style:line-height="{3*gridSize}px"
     style:left="{23*gridSize}px"
     style:top="{10*gridSize}px"
     style:font-size="{gridSize*0.9}px">
  MPS: {mpsAvg}
</h2>
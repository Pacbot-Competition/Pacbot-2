<style>

  .top-left {
    position: absolute;
    top: 0;
    left: 0;
  }

  .row {
    display: flex;
  }

  .grid-element {
    border: none;
    display: flex;
    background-color: rgba(0,0,0,0);
    padding: 0;
    margin: 0;
    box-shadow: inset 0px 0px 1px 1px rgba(255,255,255,0.05);
    justify-content: center;
    align-items: center;
  }

  .grid-element:hover {
    cursor: pointer;
    border: 2px solid rgba(255, 255, 0, 0.4);
    border-radius: 40%;
    box-shadow: none;
  }

  .grid-element:focus-visible {
    outline: none;
  }

  .hidden:hover {
    cursor: auto;
  }

  .pellet {
    background-color: #fff;
    display: block;
    color: black;
  }

  .grid-element .hidden {
    opacity: 0;
    cursor: auto;
  }

  .grid-element .hidden:hover {
    border: none;
    cursor: auto;
  }

  .grid-element .super {
    border-radius: 40%;
    transform: scale(3);
    animation: blinker 0.3s linear infinite;
  }

  @keyframes blinker {
    50% {
      opacity: 0.2;
    }
  }

</style>

<script>

  export let gridSize;
  export let grid;

  let innerWidth = 0
  let innerHeight = 0

  const hello = (i, j) => {
    console.log("hello from " + i + " , " + j)
  }

  const pelletMods = [' hidden', '', ' super']

  $: pellet_size = ~~(gridSize/6 + 0.5)  

</script>

<svelte:window bind:innerWidth bind:innerHeight />

<div class="top-left"> 
  {#each {length:31} as _, i}
    <div class="row">
      {#each {length:28} as _, j}
        <button on:click={() => hello(i, j)} class={"grid-element" + pelletMods[grid[i][j]]} style:width="{gridSize}px" style:height="{gridSize}px">
          <span class={"pellet" + pelletMods[grid[i][j]]} style:width="{pellet_size}px" style:height="{pellet_size}px"/>
        </button>
      {/each}
    </div>
  {/each}
</div>

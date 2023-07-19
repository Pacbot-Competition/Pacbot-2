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
    display: flex;
    background-color: rgba(0,0,0,0);
    /*box-shadow: inset 0px 0px 1px 1px #f00;*/
    padding: 0;
    margin: 0;
    justify-content: center;
    align-items: center;
  }

  .pellet {
    background-color: #fff;
    display: block;
    color: black;
  }

  .hidden {
    opacity: 0;
  }

  .hidden:hover {
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

  import { afterUpdate } from "svelte";
  export let grid;

  let innerWidth = 0
  let innerHeight = 0

  $: grid_size = (innerHeight * 28 < innerWidth * 31) ? (innerHeight / 31) : (innerWidth / 28)

  afterUpdate(() => {
    console.log('updated')
  })

  var hello = (i, j) => {
    console.log("hello from " + i + " , " + j)
  }

  var pelletMods = [' hidden', '', ' super']

  var pelletState = (i, j) => {
    if (i === 0 || i === 30 || j === 0 || j === 27) return 0;
    let rowCondition = (i === 3) || (i === 23);
    let colCondition = (j === 1) || (j === 26);
    if (grid[i-1][j-1]) {
      if (rowCondition && colCondition) {
        return 2
      } else return 1
    } else return 0
  }

</script>

<svelte:window bind:innerWidth bind:innerHeight />

<div class="top-left"> 
  {#each {length:31} as _, i}
    <div class="row">
      {#each {length:28} as _, j}
        <button on:click={() => hello(i, j)} class={"grid-element" + pelletMods[pelletState(i, j)]} style:width="{grid_size}px" style:height="{grid_size}px">
          <span class={"pellet" + pelletMods[pelletState(i, j)]} style:padding="{grid_size/12}px"/>
        </button>
      {/each}
    </div>
  {/each}
</div>

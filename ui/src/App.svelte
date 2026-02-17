<script>
  import logo from './assets/logo.jpeg'
  import Fetch from './lib/Fetch.svelte'
</script>

<main>
  <div>
    <a href="https://github.com/ryan-dayrit" target="_blank" rel="noreferrer">
      <img src={logo} class="logo" alt="Logo" />
    </a>
  </div>

  <div class="card">
    <Fetch />
  </div>

</main>

<style>
  .logo {
    height: 6em;
    padding: 1.5em;
    will-change: filter;
    transition: filter 300ms;
  }
  .logo:hover {
    filter: drop-shadow(0 0 2em #646cffaa);
  }
  .logo.svelte:hover {
    filter: drop-shadow(0 0 2em #ff3e00aa);
  }
</style>

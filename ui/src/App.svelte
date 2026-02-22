<script>
  import Router from 'svelte-spa-router';
  import Home from './pages/Home.svelte';
  import AlbumsList from './pages/AlbumsList.svelte';
  import AddAlbum from './pages/AddAlbum.svelte';

  const routes = {
    '/': Home,
    '/albums': AlbumsList,
    '/add': AddAlbum,
  };
</script>

<main>
  <nav class="nav">
    <a href="#/">Music Service</a>
    <a href="#/albums">Albums</a>
    <a href="#/add">Add Album</a>
  </nav>

  <div class="content">
    <Router {routes} />
  </div>
</main>

<style>
  .nav {
    display: flex;
    gap: 1.5rem;
    padding: 1rem 0 2rem;
    border-bottom: 1px solid rgba(255, 255, 255, 0.1);
    margin-bottom: 2rem;
  }

  .nav a {
    font-weight: 500;
  }

  .content {
    min-height: 200px;
  }
</style>

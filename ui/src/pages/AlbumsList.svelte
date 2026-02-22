<script>
  let albums;
  let isLoading = false;
  let error = null;

  async function fetchAlbums() {
    isLoading = true;
    error = null;
    try {
      const response = await fetch('http://localhost:3000/api/v1/albums');

      if (!response.ok) {
        throw new Error('Network response was not ok');
      }

      albums = await response.json();
    } catch (err) {
      error = err.message;
    } finally {
      isLoading = false;
    }
  }
</script>

<div class="albums-list">
  <h2>Albums</h2>

  <button onclick={fetchAlbums} disabled={isLoading}>
    {#if isLoading}
      Loading...
    {:else}
      Fetch Albums from REST API
    {/if}
  </button>

  {#if error}
    <p class="error">Error: {error}</p>
  {:else if albums}
    <table>
      <thead>
        <tr>
          <th>Id</th>
          <th>Title</th>
          <th>Artist</th>
          <th>Price</th>
        </tr>
      </thead>
      <tbody>
        {#each albums as album}
          <tr>
            <td>{album.id ?? album.Id}</td>
            <td>{album.title ?? album.Title}</td>
            <td>{album.artist ?? album.Artist}</td>
            <td>{album.price ?? album.Price}</td>
          </tr>
        {/each}
      </tbody>
    </table>
  {:else}
    <p>Click the button to fetch albums from REST API.</p>
  {/if}
</div>

<style>
  .albums-list {
    text-align: left;
  }

  h2 {
    margin-bottom: 1rem;
  }

  table {
    width: 100%;
    border-collapse: collapse;
    margin-top: 1rem;
  }

  th,
  td {
    border: 1px solid #ddd;
    padding: 8px;
    text-align: left;
  }

  th {
    background-color: #f2f2f2;
  }

  .error {
    color: #f87171;
  }
</style>

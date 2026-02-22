<script>
  let title = '';
  let artist = '';
  let price = '';
  let isSubmitting = false;
  let error = null;
  let success = false;

  async function handleSubmit(e) {
    e.preventDefault();
    isSubmitting = true;
    error = null;
    success = false;

    try {
      const response = await fetch('http://localhost:3000/api/v1/album', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title, artist, price: parseFloat(price) || 0 }),
      });

      if (!response.ok) {
        const err = await response.json();
        throw new Error(err.error || `HTTP ${response.status}`);
      }

      success = true;
      title = '';
      artist = '';
      price = '';
    } catch (err) {
      error = err.message;
    } finally {
      isSubmitting = false;
    }
  }
</script>

<div class="add-album">
  <h2>Add Album</h2>

  <form onsubmit={handleSubmit}>
    <label for="title">Title</label>
    <input
      id="title"
      type="text"
      bind:value={title}
      required
      placeholder="e.g. Blue Train"
    />

    <label for="artist">Artist</label>
    <input
      id="artist"
      type="text"
      bind:value={artist}
      required
      placeholder="e.g. John Coltrane"
    />

    <label for="price">Price</label>
    <input
      id="price"
      type="number"
      bind:value={price}
      step="0.01"
      min="0"
      required
      placeholder="e.g. 56.99"
    />

    <button type="submit" disabled={isSubmitting}>
      {#if isSubmitting}Adding...{:else}Add Album{/if}
    </button>
  </form>

  {#if error}
    <p class="error">Error: {error}</p>
  {/if}

  {#if success}
    <p class="success">Album added successfully.</p>
  {/if}
</div>

<style>
  .add-album {
    text-align: left;
    max-width: 400px;
    margin: 0 auto;
  }

  h2 {
    margin-bottom: 1.5rem;
  }

  form {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  label {
    font-weight: 500;
    margin-bottom: -0.5rem;
  }

  input {
    padding: 0.6em;
    font-size: 1em;
    border-radius: 8px;
    border: 1px solid #646cff;
    background-color: #1a1a1a;
    color: inherit;
    font-family: inherit;
  }

  input:focus {
    outline: 2px solid #646cff;
    outline-offset: 2px;
  }

  input::placeholder {
    color: rgba(255, 255, 255, 0.4);
  }

  .error {
    color: #f87171;
    margin-top: 1rem;
  }

  .success {
    color: #4ade80;
    margin-top: 1rem;
  }
</style>

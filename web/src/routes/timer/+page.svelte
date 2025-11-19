<script>
  import { onMount } from 'svelte';
  
  let currentTimer = null;
  let loading = false;
  let formData = {
    description: '',
    project_id: 'demo-project',
    user_id: 'demo-user'
  };

  onMount(async () => {
    await getCurrentTimer();
  });

  async function getCurrentTimer() {
    try {
      const response = await fetch('/api/time/current?user_id=demo-user');
      const result = await response.json();
      currentTimer = result.success ? result.data : null;
    } catch (error) {
      console.error('Failed to get current timer:', error);
    }
  }

  async function startTimer() {
    if (!formData.description.trim()) {
      alert('Please enter a description');
      return;
    }

    loading = true;
    try {
      const response = await fetch('/api/time/start', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(formData)
      });
      const result = await response.json();
      if (result.success) {
        currentTimer = result.data;
        formData.description = '';
      } else {
        alert('Failed to start timer: ' + result.error);
      }
    } catch (error) {
      console.error('Failed to start timer:', error);
      alert('Failed to start timer');
    } finally {
      loading = false;
    }
  }

  async function stopTimer() {
    loading = true;
    try {
      const response = await fetch('/api/time/stop?user_id=demo-user', {
        method: 'POST'
      });
      const result = await response.json();
      if (result.success) {
        currentTimer = null;
      } else {
        alert('Failed to stop timer: ' + result.error);
      }
    } catch (error) {
      console.error('Failed to stop timer:', error);
      alert('Failed to stop timer');
    } finally {
      loading = false;
    }
  }

  async function pauseTimer() {
    loading = true;
    try {
      const response = await fetch('/api/time/pause?user_id=demo-user', {
        method: 'POST'
      });
      const result = await response.json();
      if (result.success) {
        currentTimer = result.data;
      }
    } catch (error) {
      console.error('Failed to pause timer:', error);
    } finally {
      loading = false;
    }
  }

  async function resumeTimer() {
    loading = true;
    try {
      const response = await fetch('/api/time/resume?user_id=demo-user', {
        method: 'POST'
      });
      const result = await response.json();
      if (result.success) {
        currentTimer = result.data;
      }
    } catch (error) {
      console.error('Failed to resume timer:', error);
    } finally {
      loading = false;
    }
  }

  function formatDuration(seconds) {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    const secs = seconds % 60;
    return `${hours.toString().padStart(2, '0')}:${minutes.toString().padStart(2, '0')}:${secs.toString().padStart(2, '0')}`;
  }
</script>

<div class="max-w-2xl mx-auto space-y-8">
  <!-- Header -->
  <div class="text-center">
    <h2 class="text-3xl font-bold text-gray-900 mb-2">Time Tracker</h2>
    <p class="text-gray-600">Track your work sessions</p>
  </div>

  <!-- Current Timer Display -->
  {#if currentTimer}
    <div class="bg-white rounded-lg shadow-lg p-8 text-center">
      <div class="text-6xl font-mono font-bold text-blue-600 mb-4">
        {formatDuration(currentTimer.duration || 0)}
      </div>
      <div class="text-lg font-medium text-gray-900 mb-2">
        {currentTimer.description}
      </div>
      <div class="text-sm text-gray-500 mb-6">
        Project: {currentTimer.project_id}
      </div>
      <div class="flex justify-center space-x-4">
        {#if currentTimer.is_running}
          <button 
            on:click={pauseTimer}
            disabled={loading}
            class="bg-yellow-600 hover:bg-yellow-700 disabled:bg-gray-400 text-white font-medium py-2 px-6 rounded-md"
          >
            {loading ? 'Pausing...' : '⏸️ Pause'}
          </button>
        {:else}
          <button 
            on:click={resumeTimer}
            disabled={loading}
            class="bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white font-medium py-2 px-6 rounded-md"
          >
            {loading ? 'Resuming...' : '▶️ Resume'}
          </button>
        {/if}
        <button 
          on:click={stopTimer}
          disabled={loading}
          class="bg-red-600 hover:bg-red-700 disabled:bg-gray-400 text-white font-medium py-2 px-6 rounded-md"
        >
          {loading ? 'Stopping...' : '⏹️ Stop'}
        </button>
      </div>
    </div>
  {:else}
    <!-- Start Timer Form -->
    <div class="bg-white rounded-lg shadow p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">Start New Timer</h3>
      <div class="space-y-4">
        <div>
          <label for="description" class="block text-sm font-medium text-gray-700 mb-1">
            What are you working on?
          </label>
          <input 
            type="text" 
            id="description"
            bind:value={formData.description}
            placeholder="e.g., Designing homepage, Bug fixes, Client meeting"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          >
        </div>
        <div>
          <label for="project" class="block text-sm font-medium text-gray-700 mb-1">
            Project
          </label>
          <select 
            id="project"
            bind:value={formData.project_id}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="demo-project">Demo Project</option>
            <option value="client-website">Client Website</option>
            <option value="mobile-app">Mobile App</option>
          </select>
        </div>
        <button 
          on:click={startTimer}
          disabled={loading || !formData.description.trim()}
          class="w-full bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-3 px-4 rounded-md"
        >
          {loading ? 'Starting...' : '▶️ Start Timer'}
        </button>
      </div>
    </div>
  {/if}
</div>
<script>
  import { onMount } from 'svelte';
  
  let apiHealth = null;
  let currentTimer = null;
  let loading = false;

  onMount(async () => {
    await checkApiHealth();
    await getCurrentTimer();
  });

  async function checkApiHealth() {
    try {
      const response = await fetch('/api/health');
      apiHealth = await response.json();
    } catch (error) {
      console.error('API health check failed:', error);
      apiHealth = { status: 'error', error: error.message };
    }
  }

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
    loading = true;
    try {
      const response = await fetch('/api/time/start', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          user_id: 'demo-user',
          project_id: 'demo-project',
          description: 'Working on freelancer app'
        })
      });
      const result = await response.json();
      if (result.success) {
        currentTimer = result.data;
      }
    } catch (error) {
      console.error('Failed to start timer:', error);
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
      }
    } catch (error) {
      console.error('Failed to stop timer:', error);
    } finally {
      loading = false;
    }
  }
</script>

<div class="space-y-8">
  <!-- Header -->
  <div class="text-center">
    <h2 class="text-3xl font-bold text-gray-900 mb-2">Freelancer Dashboard</h2>
    <p class="text-gray-600">Track your time and manage your freelance projects</p>
  </div>

  <!-- API Status -->
  <div class="bg-white rounded-lg shadow p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">API Status</h3>
    {#if apiHealth}
      <div class="flex items-center space-x-2">
        {#if apiHealth.status === 'healthy'}
          <div class="w-3 h-3 bg-green-400 rounded-full"></div>
          <span class="text-green-700 font-medium">API Connected</span>
        {:else}
          <div class="w-3 h-3 bg-red-400 rounded-full"></div>
          <span class="text-red-700 font-medium">API Error</span>
        {/if}
      </div>
      <div class="mt-2 text-sm text-gray-500">
        Architecture: {apiHealth.architecture || 'Unknown'} | Database: {apiHealth.database || 'Unknown'}
      </div>
    {:else}
      <div class="text-gray-500">Checking API status...</div>
    {/if}
  </div>

  <!-- Timer Widget -->
  <div class="bg-white rounded-lg shadow p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">Quick Timer</h3>
    
    {#if currentTimer}
      <div class="text-center">
        <div class="text-4xl font-mono font-bold text-blue-600 mb-2">
          ⏱️ Running
        </div>
        <div class="text-gray-600 mb-2">
          Project: {currentTimer.project_id}
        </div>
        <div class="text-gray-500 mb-4">
          {currentTimer.description}
        </div>
        <div class="text-sm text-gray-400 mb-4">
          Started: {new Date(currentTimer.start_time).toLocaleTimeString()}
        </div>
        <button 
          on:click={stopTimer}
          disabled={loading}
          class="bg-red-600 hover:bg-red-700 disabled:bg-gray-400 text-white font-medium py-2 px-6 rounded-md"
        >
          {loading ? 'Stopping...' : 'Stop Timer'}
        </button>
      </div>
    {:else}
      <div class="text-center">
        <div class="text-4xl font-mono font-bold text-gray-400 mb-2">
          ⏸️ Idle
        </div>
        <div class="text-gray-500 mb-4">
          No active timer
        </div>
        <button 
          on:click={startTimer}
          disabled={loading}
          class="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-2 px-6 rounded-md"
        >
          {loading ? 'Starting...' : 'Start Timer'}
        </button>
      </div>
    {/if}
  </div>

  <!-- Quick Stats -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
    <div class="bg-white rounded-lg shadow p-6">
      <h4 class="text-lg font-semibold text-gray-900 mb-2">Today</h4>
      <div class="text-2xl font-bold text-blue-600">2h 30m</div>
      <div class="text-sm text-gray-500">Time tracked</div>
    </div>
    
    <div class="bg-white rounded-lg shadow p-6">
      <h4 class="text-lg font-semibold text-gray-900 mb-2">This Week</h4>
      <div class="text-2xl font-bold text-green-600">18h 45m</div>
      <div class="text-sm text-gray-500">Total hours</div>
    </div>
    
    <div class="bg-white rounded-lg shadow p-6">
      <h4 class="text-lg font-semibold text-gray-900 mb-2">Earnings</h4>
      <div class="text-2xl font-bold text-purple-600">$1,406</div>
      <div class="text-sm text-gray-500">This month</div>
    </div>
  </div>
</div>
<script>
  import { onMount } from 'svelte';

  let clients = [];
  let projects = [];
  let loading = false;
  let showClientModal = false;
  let showProjectModal = false;
  let selectedClient = null;

  // Form data
  let clientForm = {
    name: '',
    email: '',
    company: '',
    phone: '',
    address: ''
  };

  let projectForm = {
    client_id: '',
    name: '',
    description: '',
    hourly_rate: ''
  };

  onMount(async () => {
    await loadClients();
    await loadProjects();
  });

  async function loadClients() {
    try {
      const response = await fetch('/api/client/list?user_id=demo-user');
      const result = await response.json();
      clients = result.success ? result.data || [] : [];
    } catch (error) {
      console.error('Failed to load clients:', error);
      clients = [];
    }
  }

  async function loadProjects() {
    try {
      const response = await fetch('/api/project/list?user_id=demo-user');
      const result = await response.json();
      projects = result.success ? result.data || [] : [];
    } catch (error) {
      console.error('Failed to load projects:', error);
      projects = [];
    }
  }

  async function createClient() {
    loading = true;
    try {
      const response = await fetch('/api/client/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: 'demo-user',
          ...clientForm
        })
      });
      
      const result = await response.json();
      if (result.success) {
        await loadClients();
        resetClientForm();
        showClientModal = false;
      } else {
        alert('Failed to create client: ' + (result.error || 'Unknown error'));
      }
    } catch (error) {
      console.error('Failed to create client:', error);
      alert('Failed to create client: ' + error.message);
    } finally {
      loading = false;
    }
  }

  async function createProject() {
    loading = true;
    try {
      const response = await fetch('/api/project/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: 'demo-user',
          client_id: projectForm.client_id,
          name: projectForm.name,
          description: projectForm.description,
          hourly_rate: parseFloat(projectForm.hourly_rate)
        })
      });
      
      const result = await response.json();
      if (result.success) {
        await loadProjects();
        resetProjectForm();
        showProjectModal = false;
      } else {
        alert('Failed to create project: ' + (result.error || 'Unknown error'));
      }
    } catch (error) {
      console.error('Failed to create project:', error);
      alert('Failed to create project: ' + error.message);
    } finally {
      loading = false;
    }
  }

  function resetClientForm() {
    clientForm = {
      name: '',
      email: '',
      company: '',
      phone: '',
      address: ''
    };
  }

  function resetProjectForm() {
    projectForm = {
      client_id: '',
      name: '',
      description: '',
      hourly_rate: ''
    };
  }

  function openProjectModal(clientId = '') {
    projectForm.client_id = clientId;
    showProjectModal = true;
  }

  function getClientName(clientId) {
    const client = clients.find(c => c.id === clientId);
    return client ? client.name : 'Unknown Client';
  }

  function getClientProjects(clientId) {
    return projects.filter(p => p.client_id === clientId);
  }
</script>

<svelte:head>
  <title>Client Management</title>
</svelte:head>

<div class="space-y-8">
  <!-- Header -->
  <div class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-bold text-gray-900">Client Management</h2>
      <p class="text-gray-600">Manage your clients and projects</p>
    </div>
    <div class="space-x-4">
      <button 
        on:click={() => showClientModal = true}
        class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md"
      >
        Add Client
      </button>
      <button 
        on:click={() => openProjectModal()}
        class="bg-green-600 hover:bg-green-700 text-white font-medium py-2 px-4 rounded-md"
      >
        Add Project
      </button>
    </div>
  </div>

  <!-- Clients Grid -->
  <div class="grid grid-cols-1 lg:grid-cols-2 xl:grid-cols-3 gap-6">
    {#each clients as client}
      <div class="bg-white rounded-lg shadow p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-900">{client.name}</h3>
          <span class="text-sm text-gray-500">{client.company}</span>
        </div>
        
        <div class="space-y-2 text-sm text-gray-600">
          <div class="flex items-center">
            <span class="w-4 h-4 mr-2">📧</span>
            <span>{client.email}</span>
          </div>
          {#if client.phone}
            <div class="flex items-center">
              <span class="w-4 h-4 mr-2">📞</span>
              <span>{client.phone}</span>
            </div>
          {/if}
        </div>

        <!-- Client Projects -->
        <div class="mt-4">
          <div class="flex items-center justify-between mb-2">
            <h4 class="text-sm font-medium text-gray-700">Projects</h4>
            <button 
              on:click={() => openProjectModal(client.id)}
              class="text-blue-600 hover:text-blue-800 text-sm"
            >
              Add Project
            </button>
          </div>
          
          {#each getClientProjects(client.id) as project}
            <div class="bg-gray-50 rounded p-2 mb-2">
              <div class="flex justify-between items-start">
                <div>
                  <div class="font-medium text-sm">{project.name}</div>
                  <div class="text-xs text-gray-600">{project.description}</div>
                </div>
                <div class="text-right">
                  <div class="text-sm font-semibold text-green-600">
                    ${project.hourly_rate}/hr
                  </div>
                  <div class="text-xs text-gray-500">{project.status}</div>
                </div>
              </div>
            </div>
          {:else}
            <div class="text-sm text-gray-500 italic">No projects yet</div>
          {/each}
        </div>
      </div>
    {:else}
      <div class="col-span-full text-center py-12">
        <div class="text-gray-500">
          <div class="text-4xl mb-4">👥</div>
          <h3 class="text-lg font-medium mb-2">No clients yet</h3>
          <p class="text-sm">Start by adding your first client</p>
          <button 
            on:click={() => showClientModal = true}
            class="mt-4 bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md"
          >
            Add Your First Client
          </button>
        </div>
      </div>
    {/each}
  </div>
</div>

<!-- Client Modal -->
{#if showClientModal}
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4">
      <div class="px-6 py-4 border-b">
        <h3 class="text-lg font-semibold text-gray-900">Add New Client</h3>
      </div>
      
      <form on:submit|preventDefault={createClient} class="px-6 py-4 space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Name</label>
          <input 
            type="text" 
            bind:value={clientForm.name}
            required
            class="w-full border border-gray-300 rounded-md px-3 py-2"
            placeholder="Client name"
          >
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Email</label>
          <input 
            type="email" 
            bind:value={clientForm.email}
            required
            class="w-full border border-gray-300 rounded-md px-3 py-2"
            placeholder="client@example.com"
          >
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Company</label>
          <input 
            type="text" 
            bind:value={clientForm.company}
            required
            class="w-full border border-gray-300 rounded-md px-3 py-2"
            placeholder="Company name"
          >
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Phone</label>
          <input 
            type="tel" 
            bind:value={clientForm.phone}
            class="w-full border border-gray-300 rounded-md px-3 py-2"
            placeholder="Phone number"
          >
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Address</label>
          <textarea 
            bind:value={clientForm.address}
            rows="3"
            class="w-full border border-gray-300 rounded-md px-3 py-2"
            placeholder="Client address"
          ></textarea>
        </div>
        
        <div class="flex justify-end space-x-3 pt-4">
          <button 
            type="button"
            on:click={() => { showClientModal = false; resetClientForm(); }}
            class="bg-gray-300 hover:bg-gray-400 text-gray-700 font-medium py-2 px-4 rounded-md"
          >
            Cancel
          </button>
          <button 
            type="submit"
            disabled={loading}
            class="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-2 px-4 rounded-md"
          >
            {loading ? 'Creating...' : 'Create Client'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}

<!-- Project Modal -->
{#if showProjectModal}
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4">
      <div class="px-6 py-4 border-b">
        <h3 class="text-lg font-semibold text-gray-900">Add New Project</h3>
      </div>
      
      <form on:submit|preventDefault={createProject} class="px-6 py-4 space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Client</label>
          <select 
            bind:value={projectForm.client_id}
            required
            class="w-full border border-gray-300 rounded-md px-3 py-2"
          >
            <option value="">Select a client</option>
            {#each clients as client}
              <option value={client.id}>{client.name} - {client.company}</option>
            {/each}
          </select>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Project Name</label>
          <input 
            type="text" 
            bind:value={projectForm.name}
            required
            class="w-full border border-gray-300 rounded-md px-3 py-2"
            placeholder="Project name"
          >
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
          <textarea 
            bind:value={projectForm.description}
            rows="3"
            class="w-full border border-gray-300 rounded-md px-3 py-2"
            placeholder="Project description"
          ></textarea>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Hourly Rate ($)</label>
          <input 
            type="number" 
            bind:value={projectForm.hourly_rate}
            required
            min="0"
            step="0.01"
            class="w-full border border-gray-300 rounded-md px-3 py-2"
            placeholder="75.00"
          >
        </div>
        
        <div class="flex justify-end space-x-3 pt-4">
          <button 
            type="button"
            on:click={() => { showProjectModal = false; resetProjectForm(); }}
            class="bg-gray-300 hover:bg-gray-400 text-gray-700 font-medium py-2 px-4 rounded-md"
          >
            Cancel
          </button>
          <button 
            type="submit"
            disabled={loading}
            class="bg-green-600 hover:bg-green-700 disabled:bg-gray-400 text-white font-medium py-2 px-4 rounded-md"
          >
            {loading ? 'Creating...' : 'Create Project'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
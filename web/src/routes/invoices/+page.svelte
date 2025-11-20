<script>
  import { onMount } from 'svelte';

  let invoices = [];
  let clients = [];
  let projects = [];
  let loading = false;
  let showInvoiceModal = false;

  // Form data
  let invoiceForm = {
    client_id: '',
    project_id: '',
    hourly_rate: '',
    time_entries: []
  };

  // Mock time entries for demo
  let mockTimeEntries = [];

  onMount(async () => {
    await loadInvoices();
    await loadClients();
    await loadProjects();
    generateMockTimeEntries();
  });

  async function loadInvoices() {
    try {
      const response = await fetch('/api/invoice/list?user_id=demo-user');
      const result = await response.json();
      invoices = result.success ? result.data || [] : [];
    } catch (error) {
      console.error('Failed to load invoices:', error);
      invoices = [];
    }
  }

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

  function generateMockTimeEntries() {
    // Generate mock time entries for demo
    const now = new Date();
    mockTimeEntries = [
      {
        id: '1',
        description: 'Frontend development',
        start_time: new Date(now.getTime() - 4 * 60 * 60 * 1000).toISOString(),
        end_time: new Date(now.getTime() - 2 * 60 * 60 * 1000).toISOString(),
        duration: 7200
      },
      {
        id: '2',
        description: 'API integration',
        start_time: new Date(now.getTime() - 24 * 60 * 60 * 1000).toISOString(),
        end_time: new Date(now.getTime() - 22 * 60 * 60 * 1000).toISOString(),
        duration: 7200
      },
      {
        id: '3',
        description: 'Testing and debugging',
        start_time: new Date(now.getTime() - 2 * 24 * 60 * 60 * 1000).toISOString(),
        end_time: new Date(now.getTime() - 2 * 24 * 60 * 60 * 1000 + 90 * 60 * 1000).toISOString(),
        duration: 5400
      }
    ];
    invoiceForm.time_entries = [...mockTimeEntries];
  }

  async function generateInvoice() {
    loading = true;
    try {
      const response = await fetch('/api/invoice/generate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: 'demo-user',
          client_id: invoiceForm.client_id,
          project_id: invoiceForm.project_id,
          hourly_rate: parseFloat(invoiceForm.hourly_rate),
          time_entries: invoiceForm.time_entries
        })
      });
      
      const result = await response.json();
      if (result.success) {
        await loadInvoices();
        resetInvoiceForm();
        showInvoiceModal = false;
      } else {
        alert('Failed to generate invoice: ' + (result.error || 'Unknown error'));
      }
    } catch (error) {
      console.error('Failed to generate invoice:', error);
      alert('Failed to generate invoice: ' + error.message);
    } finally {
      loading = false;
    }
  }

  async function updateInvoiceStatus(invoiceId, newStatus) {
    try {
      const response = await fetch('/api/invoice/status', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: 'demo-user',
          invoice_id: invoiceId,
          status: newStatus
        })
      });
      
      const result = await response.json();
      if (result.success) {
        await loadInvoices();
      } else {
        alert('Failed to update status: ' + (result.error || 'Unknown error'));
      }
    } catch (error) {
      console.error('Failed to update status:', error);
      alert('Failed to update status: ' + error.message);
    }
  }

  function resetInvoiceForm() {
    invoiceForm = {
      client_id: '',
      project_id: '',
      hourly_rate: '',
      time_entries: [...mockTimeEntries]
    };
  }

  function getClientName(clientId) {
    const client = clients.find(c => c.id === clientId);
    return client ? client.name : 'Unknown Client';
  }

  function getProjectName(projectId) {
    const project = projects.find(p => p.id === projectId);
    return project ? project.name : 'Unknown Project';
  }

  function getStatusBadge(status) {
    const badges = {
      'draft': 'bg-gray-100 text-gray-800',
      'sent': 'bg-blue-100 text-blue-800', 
      'paid': 'bg-green-100 text-green-800',
      'overdue': 'bg-red-100 text-red-800'
    };
    return badges[status] || 'bg-gray-100 text-gray-800';
  }

  function formatDuration(seconds) {
    const hours = Math.floor(seconds / 3600);
    const minutes = Math.floor((seconds % 3600) / 60);
    return `${hours}h ${minutes}m`;
  }

  function calculateTotalHours() {
    return invoiceForm.time_entries.reduce((total, entry) => total + entry.duration, 0);
  }

  function calculateTotalAmount() {
    const totalSeconds = calculateTotalHours();
    const totalHours = totalSeconds / 3600;
    return totalHours * parseFloat(invoiceForm.hourly_rate || 0);
  }

  $: filteredProjects = invoiceForm.client_id 
    ? projects.filter(p => p.client_id === invoiceForm.client_id)
    : projects;

  $: {
    // Auto-fill hourly rate when project changes
    if (invoiceForm.project_id) {
      const project = projects.find(p => p.id === invoiceForm.project_id);
      if (project) {
        invoiceForm.hourly_rate = project.hourly_rate.toString();
      }
    }
  }
</script>

<svelte:head>
  <title>Invoice Management</title>
</svelte:head>

<div class="space-y-8">
  <!-- Header -->
  <div class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-bold text-gray-900">Invoice Management</h2>
      <p class="text-gray-600">Generate and manage your invoices</p>
    </div>
    <button 
      on:click={() => showInvoiceModal = true}
      class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md"
    >
      Generate Invoice
    </button>
  </div>

  <!-- Summary Cards -->
  <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex items-center">
        <div class="text-2xl mr-3">📄</div>
        <div>
          <div class="text-2xl font-bold text-blue-600">{invoices.length}</div>
          <div class="text-sm text-gray-500">Total Invoices</div>
        </div>
      </div>
    </div>
    
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex items-center">
        <div class="text-2xl mr-3">✅</div>
        <div>
          <div class="text-2xl font-bold text-green-600">
            {invoices.filter(inv => inv.status === 'paid').length}
          </div>
          <div class="text-sm text-gray-500">Paid</div>
        </div>
      </div>
    </div>
    
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex items-center">
        <div class="text-2xl mr-3">⏳</div>
        <div>
          <div class="text-2xl font-bold text-orange-600">
            {invoices.filter(inv => inv.status === 'sent').length}
          </div>
          <div class="text-sm text-gray-500">Pending</div>
        </div>
      </div>
    </div>
    
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex items-center">
        <div class="text-2xl mr-3">💰</div>
        <div>
          <div class="text-2xl font-bold text-purple-600">
            ${invoices.reduce((sum, inv) => sum + inv.total_amount, 0).toFixed(2)}
          </div>
          <div class="text-sm text-gray-500">Total Value</div>
        </div>
      </div>
    </div>
  </div>

  <!-- Invoices List -->
  <div class="bg-white rounded-lg shadow">
    <div class="px-6 py-4 border-b">
      <h3 class="text-lg font-semibold text-gray-900">Recent Invoices</h3>
    </div>
    
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Invoice #
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Client
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Project
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Amount
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Date
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          {#each invoices as invoice}
            <tr class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                {invoice.number}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {getClientName(invoice.client_id)}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {getProjectName(invoice.project_id)}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                ${invoice.total_amount.toFixed(2)}
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="inline-flex px-2 py-1 text-xs font-semibold rounded-full {getStatusBadge(invoice.status)}">
                  {invoice.status.charAt(0).toUpperCase() + invoice.status.slice(1)}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {new Date(invoice.created_at).toLocaleDateString()}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <div class="flex space-x-2">
                  {#if invoice.status === 'draft'}
                    <button 
                      on:click={() => updateInvoiceStatus(invoice.id, 'sent')}
                      class="text-blue-600 hover:text-blue-800"
                    >
                      Send
                    </button>
                  {/if}
                  {#if invoice.status === 'sent'}
                    <button 
                      on:click={() => updateInvoiceStatus(invoice.id, 'paid')}
                      class="text-green-600 hover:text-green-800"
                    >
                      Mark Paid
                    </button>
                  {/if}
                  <button class="text-gray-600 hover:text-gray-800">
                    View
                  </button>
                </div>
              </td>
            </tr>
          {:else}
            <tr>
              <td colspan="7" class="px-6 py-12 text-center text-gray-500">
                <div class="text-4xl mb-4">📄</div>
                <div class="text-lg font-medium mb-2">No invoices yet</div>
                <div class="text-sm">Generate your first invoice from time entries</div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>

<!-- Invoice Generation Modal -->
{#if showInvoiceModal}
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg shadow-xl w-full max-w-2xl mx-4 max-h-screen overflow-y-auto">
      <div class="px-6 py-4 border-b">
        <h3 class="text-lg font-semibold text-gray-900">Generate Invoice</h3>
      </div>
      
      <form on:submit|preventDefault={generateInvoice} class="px-6 py-4 space-y-6">
        <!-- Client and Project Selection -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label for="client" class="block text-sm font-medium text-gray-700 mb-1">Client</label>
            <select 
              id="client"
              bind:value={invoiceForm.client_id}
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
            <label for="project" class="block text-sm font-medium text-gray-700 mb-1">Project</label>
            <select 
              id="project"
              bind:value={invoiceForm.project_id}
              required
              class="w-full border border-gray-300 rounded-md px-3 py-2"
            >
              <option value="">Select a project</option>
              {#each filteredProjects as project}
                <option value={project.id}>{project.name}</option>
              {/each}
            </select>
          </div>
        </div>

        <!-- Hourly Rate -->
        <div>
          <label for="hourly-rate" class="block text-sm font-medium text-gray-700 mb-1">Hourly Rate ($)</label>
          <input 
            id="hourly-rate"
            type="number" 
            bind:value={invoiceForm.hourly_rate}
            required
            min="0"
            step="0.01"
            class="w-full border border-gray-300 rounded-md px-3 py-2"
            placeholder="75.00"
          >
        </div>

        <!-- Time Entries Preview -->
        <div>
          <h4 class="text-md font-medium text-gray-700 mb-3">Time Entries (Demo Data)</h4>
          <div class="border border-gray-200 rounded-md max-h-60 overflow-y-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500">Description</th>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500">Duration</th>
                  <th class="px-4 py-2 text-left text-xs font-medium text-gray-500">Amount</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                {#each invoiceForm.time_entries as entry}
                  <tr>
                    <td class="px-4 py-2 text-sm text-gray-900">{entry.description}</td>
                    <td class="px-4 py-2 text-sm text-gray-500">{formatDuration(entry.duration)}</td>
                    <td class="px-4 py-2 text-sm text-gray-900">
                      ${((entry.duration / 3600) * parseFloat(invoiceForm.hourly_rate || 0)).toFixed(2)}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
          
          <!-- Totals -->
          <div class="mt-4 bg-gray-50 rounded-md p-4">
            <div class="flex justify-between items-center mb-2">
              <span class="text-sm font-medium text-gray-700">Total Hours:</span>
              <span class="text-sm text-gray-900">{formatDuration(calculateTotalHours())}</span>
            </div>
            <div class="flex justify-between items-center">
              <span class="text-lg font-semibold text-gray-900">Total Amount:</span>
              <span class="text-lg font-bold text-blue-600">${calculateTotalAmount().toFixed(2)}</span>
            </div>
          </div>
        </div>
        
        <div class="flex justify-end space-x-3 pt-4">
          <button 
            type="button"
            on:click={() => { showInvoiceModal = false; resetInvoiceForm(); }}
            class="bg-gray-300 hover:bg-gray-400 text-gray-700 font-medium py-2 px-4 rounded-md"
          >
            Cancel
          </button>
          <button 
            type="submit"
            disabled={loading || !invoiceForm.client_id || !invoiceForm.project_id || !invoiceForm.hourly_rate}
            class="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-2 px-4 rounded-md"
          >
            {loading ? 'Generating...' : 'Generate Invoice'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
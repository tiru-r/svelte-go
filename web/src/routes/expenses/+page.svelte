<script>
  import { onMount } from 'svelte';

  let expenses = [];
  let projects = [];
  let loading = false;
  let showExpenseModal = false;
  let selectedProject = '';

  // Form data
  let expenseForm = {
    project_id: '',
    category: '',
    description: '',
    amount: ''
  };

  // Filter and summary
  let selectedCategory = '';
  let categories = ['Travel', 'Software', 'Hardware', 'Office Supplies', 'Meals', 'Other'];
  let monthlyTotal = 0;

  onMount(async () => {
    await loadExpenses();
    await loadProjects();
  });

  async function loadExpenses() {
    try {
      const response = await fetch('/api/expense/list?user_id=demo-user');
      const result = await response.json();
      expenses = result.success ? result.data || [] : [];
      calculateTotals();
    } catch (error) {
      console.error('Failed to load expenses:', error);
      expenses = [];
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

  async function createExpense() {
    loading = true;
    try {
      const response = await fetch('/api/expense/create', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: 'demo-user',
          project_id: expenseForm.project_id,
          category: expenseForm.category,
          description: expenseForm.description,
          amount: parseFloat(expenseForm.amount)
        })
      });
      
      const result = await response.json();
      if (result.success) {
        await loadExpenses();
        resetExpenseForm();
        showExpenseModal = false;
      } else {
        alert('Failed to create expense: ' + (result.error || 'Unknown error'));
      }
    } catch (error) {
      console.error('Failed to create expense:', error);
      alert('Failed to create expense: ' + error.message);
    } finally {
      loading = false;
    }
  }

  async function deleteExpense(expenseId) {
    if (!confirm('Are you sure you want to delete this expense?')) return;
    
    try {
      const response = await fetch('/api/expense/delete', {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: 'demo-user',
          expense_id: expenseId
        })
      });
      
      const result = await response.json();
      if (result.success) {
        await loadExpenses();
      } else {
        alert('Failed to delete expense: ' + (result.error || 'Unknown error'));
      }
    } catch (error) {
      console.error('Failed to delete expense:', error);
      alert('Failed to delete expense: ' + error.message);
    }
  }

  function resetExpenseForm() {
    expenseForm = {
      project_id: '',
      category: '',
      description: '',
      amount: ''
    };
  }

  function getProjectName(projectId) {
    const project = projects.find(p => p.id === projectId);
    return project ? project.name : 'Unknown Project';
  }

  function calculateTotals() {
    const currentMonth = new Date().getMonth();
    const currentYear = new Date().getFullYear();
    
    monthlyTotal = expenses
      .filter(expense => {
        const expenseDate = new Date(expense.created_at);
        return expenseDate.getMonth() === currentMonth && expenseDate.getFullYear() === currentYear;
      })
      .reduce((total, expense) => total + expense.amount, 0);
  }

  function getCategoryTotal(category) {
    return expenses
      .filter(expense => expense.category === category)
      .reduce((total, expense) => total + expense.amount, 0);
  }

  function getCategoryIcon(category) {
    const icons = {
      'Travel': '✈️',
      'Software': '💻',
      'Hardware': '🖥️',
      'Office Supplies': '📎',
      'Meals': '🍽️',
      'Other': '📦'
    };
    return icons[category] || '📦';
  }

  $: filteredExpenses = selectedCategory 
    ? expenses.filter(expense => expense.category === selectedCategory)
    : expenses;

  $: filteredByProject = selectedProject
    ? filteredExpenses.filter(expense => expense.project_id === selectedProject)
    : filteredExpenses;
</script>

<svelte:head>
  <title>Expense Tracking</title>
</svelte:head>

<div class="space-y-8">
  <!-- Header -->
  <div class="flex justify-between items-center">
    <div>
      <h2 class="text-3xl font-bold text-gray-900">Expense Tracking</h2>
      <p class="text-gray-600">Track and categorize your business expenses</p>
    </div>
    <button 
      on:click={() => showExpenseModal = true}
      class="bg-blue-600 hover:bg-blue-700 text-white font-medium py-2 px-4 rounded-md"
    >
      Add Expense
    </button>
  </div>

  <!-- Summary Cards -->
  <div class="grid grid-cols-1 md:grid-cols-4 gap-6">
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex items-center">
        <div class="text-2xl mr-3">💰</div>
        <div>
          <div class="text-2xl font-bold text-blue-600">${monthlyTotal.toFixed(2)}</div>
          <div class="text-sm text-gray-500">This Month</div>
        </div>
      </div>
    </div>
    
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex items-center">
        <div class="text-2xl mr-3">📊</div>
        <div>
          <div class="text-2xl font-bold text-green-600">{expenses.length}</div>
          <div class="text-sm text-gray-500">Total Expenses</div>
        </div>
      </div>
    </div>
    
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex items-center">
        <div class="text-2xl mr-3">📈</div>
        <div>
          <div class="text-2xl font-bold text-purple-600">
            ${(expenses.reduce((sum, exp) => sum + exp.amount, 0) / Math.max(expenses.length, 1)).toFixed(2)}
          </div>
          <div class="text-sm text-gray-500">Average</div>
        </div>
      </div>
    </div>
    
    <div class="bg-white rounded-lg shadow p-6">
      <div class="flex items-center">
        <div class="text-2xl mr-3">🏷️</div>
        <div>
          <div class="text-2xl font-bold text-orange-600">{categories.length}</div>
          <div class="text-sm text-gray-500">Categories</div>
        </div>
      </div>
    </div>
  </div>

  <!-- Filters -->
  <div class="bg-white rounded-lg shadow p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">Filters</h3>
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Category</label>
        <select 
          bind:value={selectedCategory}
          class="w-full border border-gray-300 rounded-md px-3 py-2"
        >
          <option value="">All Categories</option>
          {#each categories as category}
            <option value={category}>{getCategoryIcon(category)} {category}</option>
          {/each}
        </select>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Project</label>
        <select 
          bind:value={selectedProject}
          class="w-full border border-gray-300 rounded-md px-3 py-2"
        >
          <option value="">All Projects</option>
          {#each projects as project}
            <option value={project.id}>{project.name}</option>
          {/each}
        </select>
      </div>
    </div>
  </div>

  <!-- Category Breakdown -->
  <div class="bg-white rounded-lg shadow p-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">Category Breakdown</h3>
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-4">
      {#each categories as category}
        {#if getCategoryTotal(category) > 0}
          <div class="text-center p-3 bg-gray-50 rounded-lg">
            <div class="text-2xl mb-1">{getCategoryIcon(category)}</div>
            <div class="text-sm font-medium text-gray-700">{category}</div>
            <div class="text-lg font-bold text-blue-600">${getCategoryTotal(category).toFixed(2)}</div>
          </div>
        {/if}
      {/each}
    </div>
  </div>

  <!-- Expenses List -->
  <div class="bg-white rounded-lg shadow">
    <div class="px-6 py-4 border-b">
      <h3 class="text-lg font-semibold text-gray-900">
        Recent Expenses ({filteredByProject.length})
      </h3>
    </div>
    
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Date
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Description
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Category
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Project
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Amount
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          {#each filteredByProject as expense}
            <tr class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {new Date(expense.created_at).toLocaleDateString()}
              </td>
              <td class="px-6 py-4 text-sm text-gray-900">
                {expense.description}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <span class="inline-flex items-center">
                  {getCategoryIcon(expense.category)}
                  <span class="ml-2">{expense.category}</span>
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                {getProjectName(expense.project_id)}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                ${expense.amount.toFixed(2)}
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <button 
                  on:click={() => deleteExpense(expense.id)}
                  class="text-red-600 hover:text-red-800"
                >
                  Delete
                </button>
              </td>
            </tr>
          {:else}
            <tr>
              <td colspan="6" class="px-6 py-12 text-center text-gray-500">
                <div class="text-4xl mb-4">💸</div>
                <div class="text-lg font-medium mb-2">No expenses found</div>
                <div class="text-sm">Start tracking your business expenses</div>
              </td>
            </tr>
          {/each}
        </tbody>
      </table>
    </div>
  </div>
</div>

<!-- Expense Modal -->
{#if showExpenseModal}
  <div class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center z-50">
    <div class="bg-white rounded-lg shadow-xl w-full max-w-md mx-4">
      <div class="px-6 py-4 border-b">
        <h3 class="text-lg font-semibold text-gray-900">Add New Expense</h3>
      </div>
      
      <form on:submit|preventDefault={createExpense} class="px-6 py-4 space-y-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Project</label>
          <select 
            bind:value={expenseForm.project_id}
            required
            class="w-full border border-gray-300 rounded-md px-3 py-2"
          >
            <option value="">Select a project</option>
            {#each projects as project}
              <option value={project.id}>{project.name}</option>
            {/each}
          </select>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Category</label>
          <select 
            bind:value={expenseForm.category}
            required
            class="w-full border border-gray-300 rounded-md px-3 py-2"
          >
            <option value="">Select a category</option>
            {#each categories as category}
              <option value={category}>{getCategoryIcon(category)} {category}</option>
            {/each}
          </select>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
          <input 
            type="text" 
            bind:value={expenseForm.description}
            required
            class="w-full border border-gray-300 rounded-md px-3 py-2"
            placeholder="Expense description"
          >
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Amount ($)</label>
          <input 
            type="number" 
            bind:value={expenseForm.amount}
            required
            min="0"
            step="0.01"
            class="w-full border border-gray-300 rounded-md px-3 py-2"
            placeholder="0.00"
          >
        </div>
        
        <div class="flex justify-end space-x-3 pt-4">
          <button 
            type="button"
            on:click={() => { showExpenseModal = false; resetExpenseForm(); }}
            class="bg-gray-300 hover:bg-gray-400 text-gray-700 font-medium py-2 px-4 rounded-md"
          >
            Cancel
          </button>
          <button 
            type="submit"
            disabled={loading}
            class="bg-blue-600 hover:bg-blue-700 disabled:bg-gray-400 text-white font-medium py-2 px-4 rounded-md"
          >
            {loading ? 'Adding...' : 'Add Expense'}
          </button>
        </div>
      </form>
    </div>
  </div>
{/if}
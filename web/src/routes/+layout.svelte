<script>
  import '../app.css';
</script>

<main class="min-h-screen bg-gray-50">
  <nav class="bg-white shadow-sm border-b border-gray-200">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <div class="flex items-center">
          <h1 class="text-xl font-semibold text-gray-900">⏱️ Freelancer Tracker</h1>
        </div>
        <div class="flex space-x-4">
          <a href="/" class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">Dashboard</a>
          <a href="/timer" class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">Timer</a>
          <a href="/projects" class="text-gray-600 hover:text-gray-900 px-3 py-2 rounded-md text-sm font-medium">Projects</a>
        </div>
      </div>
    </div>
  </nav>
  
  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <slot />
  </div>
</main>
<script>
	import { formatDate, apiRequest, publishEvent, getNATSStats } from '$lib/utils.js';
	import { onMount } from 'svelte';

	/** @type {string} */
	let healthStatus = 'Checking...';
	
	/** @type {string} */
	let currentDate = '';

	/** 
	 * @type {{
	 *   connections: number,
	 *   subscriptions: number,
	 *   messages: {in: number, out: number},
	 *   bytes: {in: number, out: number},
	 *   uptime: string
	 * }|null} 
	 */
	let natsStats = null;

	/** @type {number} */
	let eventCount = 0;

	/**
	 * Check the API health status
	 * @returns {Promise<void>}
	 */
	async function checkHealth() {
		// Publish event for health check
		await publishEvent('health_check_request', { 
			timestamp: Date.now(),
			source: 'frontend_button'
		});

		const response = await apiRequest('/api/health');
		if (response.success) {
			healthStatus = `✅ ${response.data.status}`;
		} else {
			healthStatus = `❌ ${response.error}`;
		}
		
		eventCount++;
		await refreshNATSStats();
	}

	/**
	 * Refresh NATS statistics
	 * @returns {Promise<void>}
	 */
	async function refreshNATSStats() {
		const response = await getNATSStats();
		if (response.success) {
			natsStats = response.data;
		}
	}

	/**
	 * Send a test event
	 * @returns {Promise<void>}
	 */
	async function sendTestEvent() {
		await publishEvent('user_interaction', {
			action: 'test_event_button_click',
			timestamp: Date.now(),
			userAgent: navigator.userAgent
		});
		
		eventCount++;
		await refreshNATSStats();
	}

	onMount(() => {
		currentDate = formatDate(new Date(), {
			formatOptions: { 
				weekday: 'long', 
				year: 'numeric', 
				month: 'long', 
				day: 'numeric' 
			}
		});
		checkHealth();
		refreshNATSStats();
	});
</script>

<div class="min-h-screen bg-gray-100 flex items-center justify-center p-4">
	<div class="max-w-2xl w-full space-y-6">
		<!-- Main App Info -->
		<div class="bg-white rounded-lg shadow-md p-6">
			<h1 class="text-2xl font-bold text-gray-900 mb-4">Svelte + Go App</h1>
			<p class="text-gray-600 mb-4">Event-driven architecture with embedded NATS</p>
			
			<div class="space-y-2 text-sm">
				<div class="flex justify-between">
					<span class="text-gray-500">Date:</span>
					<span class="text-gray-700">{currentDate}</span>
				</div>
				<div class="flex justify-between">
					<span class="text-gray-500">API Status:</span>
					<span class="text-gray-700">{healthStatus}</span>
				</div>
				<div class="flex justify-between">
					<span class="text-gray-500">Events Sent:</span>
					<span class="text-gray-700">{eventCount}</span>
				</div>
			</div>
			
			<div class="flex space-x-3 mt-4">
				<button 
					class="flex-1 bg-blue-500 hover:bg-blue-600 text-white font-medium py-2 px-4 rounded transition-colors"
					on:click={checkHealth}
				>
					Health Check
				</button>
				<button 
					class="flex-1 bg-green-500 hover:bg-green-600 text-white font-medium py-2 px-4 rounded transition-colors"
					on:click={sendTestEvent}
				>
					Send Event
				</button>
			</div>
		</div>

		<!-- NATS Statistics -->
		{#if natsStats}
		<div class="bg-white rounded-lg shadow-md p-6">
			<h2 class="text-xl font-bold text-gray-900 mb-4">NATS Server Stats</h2>
			
			<div class="grid grid-cols-2 gap-4 text-sm">
				<div class="space-y-2">
					<div class="flex justify-between">
						<span class="text-gray-500">Connections:</span>
						<span class="text-gray-700">{natsStats?.connections ?? 0}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-gray-500">Subscriptions:</span>
						<span class="text-gray-700">{natsStats?.subscriptions ?? 0}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-gray-500">Uptime:</span>
						<span class="text-gray-700">{natsStats?.uptime ?? 'N/A'}</span>
					</div>
				</div>
				
				<div class="space-y-2">
					<div class="flex justify-between">
						<span class="text-gray-500">Messages In:</span>
						<span class="text-gray-700">{natsStats?.messages?.in ?? 0}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-gray-500">Messages Out:</span>
						<span class="text-gray-700">{natsStats?.messages?.out ?? 0}</span>
					</div>
					<div class="flex justify-between">
						<span class="text-gray-500">Bytes In:</span>
						<span class="text-gray-700">{natsStats?.bytes?.in ? (natsStats.bytes.in / 1024).toFixed(1) : '0.0'} KB</span>
					</div>
				</div>
			</div>
			
			<button 
				class="mt-4 w-full bg-gray-500 hover:bg-gray-600 text-white font-medium py-2 px-4 rounded transition-colors"
				on:click={refreshNATSStats}
			>
				Refresh Stats
			</button>
		</div>
		{/if}
	</div>
</div>
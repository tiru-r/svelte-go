/**
 * @fileoverview Utility functions for the application
 */

/**
 * Formats a date to a human-readable string
 * @param {Date} date - The date to format
 * @param {Object} [options] - Formatting options
 * @param {string} [options.locale='en-US'] - The locale to use
 * @param {Intl.DateTimeFormatOptions} [options.formatOptions] - Date formatting options
 * @returns {string} The formatted date string
 * @example
 * formatDate(new Date())
 * // Returns: "12/25/2023"
 */
export function formatDate(date, options = {}) {
  const { locale = 'en-US', formatOptions = {} } = options;
  return date.toLocaleDateString(locale, formatOptions);
}

/**
 * Debounces a function call
 * @param {Function} func - The function to debounce
 * @param {number} delay - The delay in milliseconds
 * @returns {Function} The debounced function
 * @example
 * const debouncedSearch = debounce(searchFunction, 300);
 */
export function debounce(func, delay) {
  /** @type {ReturnType<typeof setTimeout>|undefined} */
  let timeoutId;
  return function (/** @type {any[]} */ ...args) {
    clearTimeout(timeoutId);
    timeoutId = setTimeout(() => func(...args), delay);
  };
}

/**
 * @typedef {Object} ApiResponse
 * @property {boolean} success - Whether the request was successful
 * @property {*} data - The response data
 * @property {string} [error] - Error message if request failed
 */

/**
 * @typedef {Object} NATSMessageStats
 * @property {number} in - Number of messages received
 * @property {number} out - Number of messages sent
 */

/**
 * @typedef {Object} NATSByteStats
 * @property {number} in - Number of bytes received
 * @property {number} out - Number of bytes sent
 */

/**
 * @typedef {Object} NATSStats
 * @property {number} connections - Number of active connections
 * @property {number} subscriptions - Number of active subscriptions
 * @property {NATSMessageStats} messages - Message statistics
 * @property {NATSByteStats} bytes - Byte transfer statistics
 * @property {string} uptime - Server uptime duration
 */

/**
 * Makes an API request to the backend
 * @param {string} endpoint - The API endpoint
 * @param {Object} [options] - Request options
 * @param {'GET'|'POST'|'PUT'|'DELETE'} [options.method='GET'] - HTTP method
 * @param {Object} [options.body] - Request body for POST/PUT requests
 * @returns {Promise<ApiResponse>} The API response
 * @example
 * const response = await apiRequest('/api/health');
 * if (response.success) {
 *   console.log(response.data);
 * }
 */
export async function apiRequest(endpoint, options = {}) {
  const { method = 'GET', body } = options;
  
  try {
    const response = await fetch(endpoint, {
      method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: body ? JSON.stringify(body) : undefined,
    });
    
    const data = await response.json();
    
    return {
      success: response.ok,
      data: response.ok ? data : null,
      error: response.ok ? null : data.message || 'Request failed',
    };
  } catch (error) {
    return {
      success: false,
      data: null,
      error: error instanceof Error ? error.message : 'Request failed',
    };
  }
}

/**
 * Publishes an event to the backend event system
 * @param {string} type - The event type
 * @param {Object} data - The event data
 * @returns {Promise<ApiResponse>} The API response
 * @example
 * await publishEvent('button_click', { button: 'refresh', timestamp: Date.now() });
 */
export async function publishEvent(type, data) {
  return apiRequest('/api/events', {
    method: 'POST',
    body: { type, data }
  });
}

/**
 * Gets NATS server statistics
 * @returns {Promise<ApiResponse & {data: NATSStats}>} The NATS stats response
 * @example
 * const stats = await getNATSStats();
 * if (stats.success) {
 *   console.log('Messages:', stats.data.messages);
 *   console.log('Uptime:', stats.data.uptime);
 * }
 */
export async function getNATSStats() {
  return apiRequest('/api/nats/stats');
}
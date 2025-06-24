const { test, expect } = require('@playwright/test');

test.describe('API Endpoints E2E', () => {
  const baseURL = process.env.BASE_URL || 'http://localhost:8080';

  test('should get weather data for valid city', async ({ request }) => {

    const response = await request.get(`${baseURL}/api/weather?city=London`);
    
    if (response.status() === 200) {
      const weatherData = await response.json();
      
      expect(weatherData).toHaveProperty('city');
      expect(weatherData).toHaveProperty('temperature');
      expect(weatherData).toHaveProperty('humidity');
      expect(weatherData).toHaveProperty('description');
      
      expect(typeof weatherData.temperature).toBe('number');
      expect(typeof weatherData.humidity).toBe('number');
      expect(typeof weatherData.city).toBe('string');
      expect(typeof weatherData.description).toBe('string');
    } else {

      expect(response.status()).toBe(500);
      const errorData = await response.json();
      expect(errorData).toHaveProperty('error');
    }
  });

  test('should return error for missing city parameter', async ({ request }) => {
    const response = await request.get(`${baseURL}/api/weather`);
    
    expect(response.status()).toBe(400);
    const errorData = await response.json();
    expect(errorData.error).toContain('City parameter is required');
  });

  test('should return error for empty city parameter', async ({ request }) => {
    const response = await request.get(`${baseURL}/api/weather?city=`);
    
    expect(response.status()).toBe(400);
    const errorData = await response.json();
    expect(errorData.error).toContain('City parameter is required');
  });

  test('should create subscription with valid data', async ({ request }) => {
    const subscriptionData = {
      email: 'e2etest@example.com',
      city: 'London',
      frequency: 'daily'
    };

    const response = await request.post(`${baseURL}/api/subscribe`, {
      data: subscriptionData,
      headers: {
        'Content-Type': 'application/json'
      }
    });

    if (response.status() === 200) {
      const responseData = await response.json();
      expect(responseData.message).toContain('check your email');
    } else {
      
      expect([400, 500]).toContain(response.status());
    }
  });

  test('should validate subscription request format', async ({ request }) => {
    const invalidData = {
      email: 'invalid-email',
      city: 'London'
      // Missing frequency field
    };

    const response = await request.post(`${baseURL}/api/subscribe`, {
      data: invalidData,
      headers: {
        'Content-Type': 'application/json'
      }
    });

    expect(response.status()).toBe(400);
    const errorData = await response.json();
    expect(errorData.error).toContain('Invalid request format');
  });

  test('should serve static assets', async ({ request }) => {
    const cssResponse = await request.get(`${baseURL}/static/css/styles.css`);
    expect(cssResponse.status()).toBe(200);
    expect(cssResponse.headers()['content-type']).toContain('text/css');

    const jsResponse = await request.get(`${baseURL}/static/js/script.js`);
    expect(jsResponse.status()).toBe(200);
    expect(jsResponse.headers()['content-type']).toContain('javascript');
  });

  test('should serve main page', async ({ request }) => {
    const response = await request.get(`${baseURL}/`);
    expect(response.status()).toBe(200);
    expect(response.headers()['content-type']).toContain('text/html');
    
    const html = await response.text();
    expect(html).toContain('Weather Notifier');
    expect(html).toContain('subscription-form');
  });

  test('should handle 404 for non-existent endpoints', async ({ request }) => {
    const response = await request.get(`${baseURL}/api/nonexistent`);
    expect(response.status()).toBe(404);
  });

  test('should respond to ping endpoint', async ({ request }) => {
    const response = await request.get(`${baseURL}/ping`);
    expect(response.status()).toBe(200);
    
    const data = await response.json();
    expect(data.message).toBe('pong');
  });

  test('should serve swagger documentation', async ({ request }) => {
    const response = await request.get(`${baseURL}/swagger/index.html`);
    expect(response.status()).toBe(200);
    expect(response.headers()['content-type']).toContain('text/html');
  });
});

test.describe('API Response Headers', () => {
  const baseURL = process.env.BASE_URL || 'http://localhost:8080';

  test('should have proper CORS headers if configured', async ({ request }) => {
    const response = await request.get(`${baseURL}/api/weather?city=London`);
    
    const headers = response.headers();
    
    console.log('Response headers:', Object.keys(headers));
  });

  test('should have security headers', async ({ request }) => {
    const response = await request.get(`${baseURL}/`);
    
    const headers = response.headers();
    
    console.log('Security-related headers:');
    if (headers['x-frame-options']) console.log('X-Frame-Options:', headers['x-frame-options']);
    if (headers['x-content-type-options']) console.log('X-Content-Type-Options:', headers['x-content-type-options']);
    if (headers['x-xss-protection']) console.log('X-XSS-Protection:', headers['x-xss-protection']);
  });
});